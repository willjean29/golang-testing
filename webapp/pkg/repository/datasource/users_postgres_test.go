package datasource

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"
	"webapp/pkg/data"
	"webapp/pkg/repository"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "users_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var pool *dockertest.Pool
var testDB *sql.DB
var testRepo repository.Repository

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	pool = p
	opt := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opt)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = createTables()
	if err != nil {
		log.Fatalf("Could not create tables: %s", err)
	}

	testRepo = &PostgresDB{DB: testDB}

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Exit(code)
}

func createTables() error {
	tableSql, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		log.Println("Error reading sql file", err)
		return err
	}

	_, err = testDB.Exec(string(tableSql))
	if err != nil {
		log.Println("Error creating tables", err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Errorf("pingDB() failed: %s", err)
	}
}

func TestPostgresDBInsertUser(t *testing.T) {
	cleanDatabase()
	tests := []struct {
		name    string
		user    data.User
		wantErr bool
	}{
		{
			name: "insert user",
			user: data.User{
				Email:     "admin@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password",
				IsAdmin:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "error inserting user when fail sql query",
			user: data.User{
				Email:     "admin@example.com",
				FirstName: strings.Repeat("Jhon", 256),
				LastName:  "Doe",
				Password:  "password",
				IsAdmin:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "error inserting user when fail generate password hash",
			user: data.User{
				Email:     "admin@example.com",
				FirstName: "Jhon",
				LastName:  "Doe",
				Password:  strings.Repeat("secret", 256),
				IsAdmin:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := testRepo.InsertUser(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostgresDB.InsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && id != 1 {
				t.Errorf("PostgresDB.InsertUser() = %v, want %v", id, 1)
			}
		})
	}
	cleanDatabase()
}

func TestPostgresDBAllUsers(t *testing.T) {
	cleanDatabase()
	user := data.User{
		Email:     "admin@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	testRepo.InsertUser(user)

	users, err := testRepo.AllUsers()
	if err != nil {
		t.Errorf("PostgresDB.AllUsers() error = %v", err)
	}

	if len(users) <= 0 {
		t.Errorf("PostgresDB.AllUsers() = %v", len(users))
	}
	cleanDatabase()
}

func TestPostgresDBGetUser(t *testing.T) {
	cleanDatabase()
	newUser := data.User{
		Email:     "admin@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, err := testRepo.InsertUser(newUser)

	if err != nil {
		t.Errorf("PostgresDB.InsertUser() error = %v", err)
	}
	user, err := testRepo.GetUser(id)
	log.Println("user", user)
	if err != nil {
		t.Errorf("PostgresDB.GetUser() error = %v", err)
	}

	cleanDatabase()
}

func TestPostgresDBGetUserByEmail(t *testing.T) {
	cleanDatabase()
	newUser := data.User{
		Email:     "admin@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := testRepo.InsertUser(newUser)

	if err != nil {
		t.Errorf("PostgresDB.InsertUser() error = %v", err)
	}

	user, err := testRepo.GetUserByEmail(newUser.Email)
	if err != nil {
		t.Errorf("PostgresDB.GetUserByEmail() error = %v", err)
	}

	if user.Email != newUser.Email {
		t.Errorf("PostgresDB.GetUserByEmail() = %v, want %v", user.Email, newUser.Email)
	}
	cleanDatabase()
}

func TestPostgresDBUpdateUser(t *testing.T) {
	cleanDatabase()
	newUser := data.User{
		Email:     "admin@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		IsAdmin:   1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	id, _ := testRepo.InsertUser(newUser)
	user, _ := testRepo.GetUser(id)
	user.FirstName = "Jean"
	user.LastName = "Osco"

	err := testRepo.UpdateUser(*user)
	if err != nil {
		t.Errorf("PostgresDB.UpdateUser() error = %v", err)
	}

	user, _ = testRepo.GetUser(id)
	if user.FirstName != "Jean" && user.LastName != "Osco" {
		t.Errorf("expected Jean Osco, got %s %s", user.FirstName, user.LastName)
	}
}

func cleanDatabase() {
	_, err := testDB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	if err != nil {
		log.Println("Error cleaning database", err)
	}
}
