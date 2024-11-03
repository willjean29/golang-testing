package main

import "net/http"

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// Authenticate user
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
