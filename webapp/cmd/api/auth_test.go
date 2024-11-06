package main

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_app_getTokenFromHeaderAndVerify(t *testing.T) {
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
		issuer        string
	}{
		{
			name:          "Valid token",
			token:         fmt.Sprintf("Bearer %s", tokens.Token),
			expectedError: false,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "Expired token",
			token:         fmt.Sprintf("Bearer %s", expiredToken),
			expectedError: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "Empty header",
			token:         "",
			expectedError: true,
			setHeader:     false,
		},
		{
			name:          "Invalid token",
			token:         fmt.Sprintf("Bearer %s1", tokens.Token),
			expectedError: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "No bearer header",
			token:         fmt.Sprintf("Bea %s1", tokens.Token),
			expectedError: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "Three header parts",
			token:         fmt.Sprintf("Bea %s 1", tokens.Token),
			expectedError: true,
			setHeader:     true,
			issuer:        app.Domain,
		},
		{
			name:          "Wrong issuer",
			token:         fmt.Sprintf("Bearer %s", tokens.Token),
			expectedError: true,
			setHeader:     true,
			issuer:        "company.com",
		},
	}

	for _, tt := range tests {
		t.Run((tt.name), func(t *testing.T) {
			if tt.issuer != app.Domain {
				app.Domain = tt.issuer
				tokens, _ = app.generateTokenPair(testUser)
			}
			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.setHeader {
				req.Header.Add("Authorization", tt.token)
			}

			_, _, err := app.getTokenFromHeaderAndVerify(res, req)

			if err != nil && !tt.expectedError {
				t.Errorf("Expected no error, but got %v", err)
			}

			if err == nil && tt.expectedError {
				t.Error("Expected error, but got nil")
			}

			app.Domain = "example.com"

		})
	}

}
