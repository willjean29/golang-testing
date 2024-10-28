package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"webapp/pkg/repository/datasource"
)

var app application

func TestMain(m *testing.M) {
	fmt.Println("Setting up")
	app.Session = getSession()
	app.DSN = "host=localhost port=5432 user=postgres password=postgres dbname=users"
	conn, err := app.connectToDB()

	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}
	app.DB = &datasource.PostgresDB{DB: conn}
	os.Exit(m.Run())
}
