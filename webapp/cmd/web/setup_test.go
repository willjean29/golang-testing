package main

import (
	"fmt"
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	fmt.Println("Setting up")
	app.Session = getSession()
	os.Exit(m.Run())
}
