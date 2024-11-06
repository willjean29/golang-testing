package main

import (
	"log"
	"os"
	"testing"
	"webapp/pkg/data"
	"webapp/pkg/repository/datasource"
)

var app application
var expiredToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXVkIjoiZXhhbXBsZS5jb20iLCJleHAiOjE3MzA1MDkwNTcsImlzcyI6ImV4YW1wbGUuY29tIiwibmFtZSI6IkFkbWluIFVzZXIiLCJzdWIiOiIxIn0.cVrukpv1lOlgkEXJqeF_5_L3ZcxbON8h3LIG9j5-6Ts"
var users = []*data.User{
	{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		Password:  "$2a$14$ajq8Q7fbtFRQvXpdCq7Jcuy.Rx1h/L4J60Otx.gyNLbAYctGMJ9tK",
	},
}

func TestMain(m *testing.M) {
	log.Println("Setting up unit test")
	app.Domain = "example.com"
	app.JWTSecret = "test123"
	app.DB = &datasource.TestDB{
		Users: users,
	}
	os.Exit(m.Run())
}
