package datasource

import (
	"database/sql"
	"webapp/pkg/data"
)

type TestDB struct {
	users []*data.User
}

func (m *TestDB) Connection() *sql.DB {
	return nil
}

// AllUsers returns all users as a slice of *data.User
func (m *TestDB) AllUsers() ([]*data.User, error) {
	var users []*data.User
	return users, nil
}

// GetUser returns one user by id
func (m *TestDB) GetUser(id int) (*data.User, error) {
	var user data.User

	return &user, nil
}

// GetUserByEmail returns one user by email address
func (m *TestDB) GetUserByEmail(email string) (*data.User, error) {
	var user data.User

	return &user, nil
}

// UpdateUser updates one user in the database
func (m *TestDB) UpdateUser(u data.User) error {

	return nil
}

// DeleteUser deletes one user from the database, by id
func (m *TestDB) DeleteUser(id int) error {

	return nil
}

// InsertUser inserts a new user into the database, and returns the ID of the newly inserted row
func (m *TestDB) InsertUser(user data.User) (int, error) {

	return 1, nil
}

// ResetPassword is the method we will use to change a user's password.
func (m *TestDB) ResetPassword(id int, password string) error {

	return nil
}

// InsertUserImage inserts a user profile image into the database.
func (m *TestDB) InsertUserImage(i data.UserImage) (int, error) {

	return 1, nil
}
