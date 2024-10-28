package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(app.addIPToContext)
	router.Use(app.Session.LoadAndSave)

	router.Get("/", app.Home)
	router.Post("/login", app.Login)

	router.Route("/user", func(r chi.Router) {
		r.Use(app.auth)
		r.Get("/profile", app.Profile)
	})

	fs := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static", fs))

	return router
}
