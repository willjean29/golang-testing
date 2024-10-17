package main

import (
	"fmt"
	"log"
	"net/http"
)

type application struct {
}

func main() {
	app := &application{}

	router := app.routes()

	fmt.Println("Starting server on :8080")

	err := http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatal("Error starting server: " + err.Error())
	}
}
