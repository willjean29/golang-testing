package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_app_enableCORS(t *testing.T) {

	tests := []struct {
		name   string
		method string
	}{
		{
			name:   "OPTIONS method",
			method: "OPTIONS",
		},
		{
			name:   "GET method",
			method: "GET",
		},
	}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerTest := app.enableCORS(nextHandler)

			req := httptest.NewRequest(tt.method, "/testing", nil)
			res := httptest.NewRecorder()

			handlerTest.ServeHTTP(res, req)
			log.Println("HEADER", req.Method)
			if res.Code != http.StatusOK {
				t.Error("Expected status 200, got", res.Code)
			}
		})
	}

}
