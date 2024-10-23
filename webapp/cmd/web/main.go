package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
}

func main() {
	app := &application{}

	app.Session = getSession()

	fmt.Println("Starting server on :8080")

	err := http.ListenAndServe("localhost:8080", app.routes())

	if err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
