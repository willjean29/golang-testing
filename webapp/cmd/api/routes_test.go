package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/test", "GET"},
		{"/auth", "POST"},
		{"/refresh-token", "POST"},
		{"/users/", "GET"},
		{"/users/{userId}", "GET"},
		{"/users/{userId}", "DELETE"},
		{"/users/", "POST"},
	}

	app = application{}

	router := app.routes()

	chiRoutes := router.(chi.Routes)

	for _, route := range registered {
		if !routeExists(chiRoutes, route.route, route.method) {
			t.Errorf("route %s %s not found", route.method, route.route)
		}
	}
}

func routeExists(routes chi.Routes, testRoute string, testMethod string) bool {
	found := false
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}
		return nil
	}
	_ = chi.Walk(routes, walkFunc)
	return found
}
