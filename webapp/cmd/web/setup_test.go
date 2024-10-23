package main

import (
	"fmt"
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	fmt.Println("Setting up")
	os.Exit(m.Run())
}
