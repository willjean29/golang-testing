package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
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
			if res.Code != http.StatusOK {
				t.Error("Expected status 200, got", res.Code)
			}
		})
	}
}

func Test_app_authRequired(t *testing.T) {
	testUser := &data.User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		IsAdmin:   1,
	}
	tokens, _ := app.generateTokenPair(testUser)
	tests := []struct {
		name          string
		token         string
		expectedError bool
		setHeader     bool
	}{
		{
			name:          "Valid token",
			token:         fmt.Sprintf("Bearer %s", tokens.Token),
			expectedError: false,
			setHeader:     true,
		},
		{
			name:          "Expired token",
			token:         fmt.Sprintf("Bearer %s", expiredToken),
			expectedError: true,
			setHeader:     true,
		},
	}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerTest := app.authRequired(nextHandler)
			req := httptest.NewRequest("GET", "/testing", nil)
			if tt.setHeader {
				req.Header.Add("Authorization", tt.token)
			}
			res := httptest.NewRecorder()

			handlerTest.ServeHTTP(res, req)
			log.Println("res.Code", res.Code)
			if res.Code != http.StatusUnauthorized && tt.expectedError {
				t.Error("Expected status 401, got", res.Code)
			}

			if res.Code != http.StatusOK && !tt.expectedError {
				t.Error("Expected status 200, got", res.Code)
			}
		})
	}
}
