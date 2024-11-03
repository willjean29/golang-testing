package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Post("/auth", app.authenticate)
	router.Post("/refresh-token", app.refresh)

	// Protected routes
	router.Route("/users", func(r chi.Router) {
		r.Get("/", app.allUsers)
		r.Get("/{userId}", app.getUser)
		r.Put("/{userId}", app.updateUser)
		r.Delete("/{userId}", app.deleteUser)
		r.Post("/", app.insertUser)
	})

	return router
}