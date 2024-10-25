package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
	DSN     string
	DB      *sql.DB
}

func main() {
	app := &application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users ", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()

	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}

	app.DB = conn

	app.Session = getSession()

	log.Println("Starting server on :8080")

	err = http.ListenAndServe("localhost:8080", app.routes())

	if err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
