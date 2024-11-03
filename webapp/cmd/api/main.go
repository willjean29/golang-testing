package main

import (
	"flag"
	"log"
	"net/http"
	"webapp/pkg/repository"
	"webapp/pkg/repository/datasource"
)

type application struct {
	DSN       string
	DB        repository.Repository
	Domain    string
	JWTSecret string
}

func main() {
	var app application
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for the application e.g company.com")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users", "Postgres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "secret123", "Secret for JWT")
	flag.Parse()

	conn, err := app.connectToDB()

	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}
	defer conn.Close()

	app.DB = &datasource.PostgresDB{DB: conn}

	log.Println("Starting server on :8080")

	err = http.ListenAndServe("localhost:8080", app.routes())

	if err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
