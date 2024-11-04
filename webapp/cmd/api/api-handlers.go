package main

import (
	"net/http"
)

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {

}

func (app *application) refresh(w http.ResponseWriter, r *http.Request) {
	// Refresh user
}

func (app *application) allUsers(w http.ResponseWriter, r *http.Request) {
	// All users
}

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	// Get user
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {
	// Update user
}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Delete user
}

func (app *application) insertUser(w http.ResponseWriter, r *http.Request) {
	// Insert user
}
