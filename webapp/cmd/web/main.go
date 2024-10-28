package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"webapp/pkg/data"
	"webapp/pkg/db"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
	DSN     string
	DB      db.PostgresConn
}

func main() {
	gob.Register(data.User{})
	app := &application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users ", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()

	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}

	defer conn.Close()

	app.DB = db.PostgresConn{DB: conn}

	app.Session = getSession()

	log.Println("Starting server on :8080")

	err = http.ListenAndServe("localhost:8080", app.routes())

	if err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
