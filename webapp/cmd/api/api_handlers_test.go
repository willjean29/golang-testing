package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"webapp/pkg/data"
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

func Test_app_refresh(t *testing.T) {
	testUser := &data.User{
		ID:        1,
		Email:     "admin@example.com",
		FirstName: "Admin",
		LastName:  "User",
		IsAdmin:   1,
	}
	tests := []struct {
		name               string
		token              string
		expectedStatusCode int
		resetRefreshTime   bool
	}{
		{
			name:               "Valid token",
			token:              "",
			expectedStatusCode: http.StatusOK,
			resetRefreshTime:   true,
		},
		{
			name:               "Expired token",
			token:              expiredToken,
			expectedStatusCode: http.StatusBadRequest,
			resetRefreshTime:   false,
		},
		{
			name:               "Valid token but not yet ready to expired",
			token:              "",
			expectedStatusCode: http.StatusTooEarly,
			resetRefreshTime:   false,
		},
	}
	oldRefreshTokenExpiry := jwtRefreshTokenExpiry
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tkn string
			if tt.token == "" {
				if tt.resetRefreshTime {
					jwtRefreshTokenExpiry = time.Second * 1
				}
				tokens, _ := app.generateTokenPair(testUser)
				tkn = tokens.RefreshToken
			} else {
				tkn = tt.token
			}
			postData := strings.NewReader(url.Values{"refresh_token": {tkn}}.Encode())
			res := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/refresh-token", postData)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			app.refresh(res, req)

			if res.Code != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, res.Code)
			}
			jwtRefreshTokenExpiry = oldRefreshTokenExpiry
		})
	}

}
