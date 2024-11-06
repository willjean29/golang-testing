package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_app_authenticate(t *testing.T) {
	tests := []struct {
		name         string
		expectedBody string
		expectedCode int
	}{
		{
			name:         "Valid credentials",
			expectedBody: `{"email": "admin@example.com", "password": "secret"}`,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid credentials - password",
			expectedBody: `{"email": "admin@example.com", "password": "secret123"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid credentials - email",
			expectedBody: `{"email": "admin2@example.com", "password": "secret123"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Invalid body",
			expectedBody: `""`,
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := strings.NewReader(tt.expectedBody)
			res := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth", body)
			app.authenticate(res, req)

			if res.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedCode, res.Code)
			}
		})
	}

}
