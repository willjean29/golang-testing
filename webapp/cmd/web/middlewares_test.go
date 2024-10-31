package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"webapp/pkg/data"
)

func Test_application_addIPToContext(t *testing.T) {
	tests := []struct {
		headerName   string
		headerValue  string
		address      string
		emptyAddress bool
	}{
		{"", "", ":8080", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.3.2.1", "", false},
		{"", "", "hello:word", false},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := r.Context().Value(contentUserKey)
		if val == nil {
			t.Error(contentUserKey, "not found in context")
		}
		ip, ok := val.(string)
		if !ok {
			t.Error("value in context is not a string")
		}
		t.Log("IP in context is", ip)
	})

	for _, tt := range tests {
		handlerTest := app.addIPToContext(nextHandler)
		req := httptest.NewRequest("GET", "/testing", nil)

		if tt.emptyAddress {
			req.RemoteAddr = ""
		}

		if len(tt.headerName) > 0 {
			req.Header.Set(tt.headerName, tt.headerValue)
		}

		if len(tt.address) > 0 {
			req.RemoteAddr = tt.address
		}

		handlerTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func Test_application_ipFromContext(t *testing.T) {
	const mockIP = "192.3.2.1"
	ctx := context.Background()
	ctx = context.WithValue(ctx, contentUserKey, mockIP)
	ip := app.ipFromContext(ctx)

	if ip != mockIP {
		t.Errorf("expected %s, got %s", mockIP, ip)
	}
}

func Test_application_auth(t *testing.T) {
	tests := []struct {
		name   string
		isAuth bool
	}{
		{"logged in", true},
		{"not logged in", false},
	}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlerTest := app.auth(nextHandler)
			req := httptest.NewRequest("GET", "/testing", nil)
			req = addContextAndSessionToRequest(req)
			if tt.isAuth {
				app.Session.Put(req.Context(), "user", data.User{})
			}
			res := httptest.NewRecorder()
			handlerTest.ServeHTTP(res, req)
			if tt.isAuth && res.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, res.Code)
			}
			if !tt.isAuth && res.Code != http.StatusTemporaryRedirect {
				t.Errorf("expected status code %d, got %d", http.StatusTemporaryRedirect, res.Code)
			}
		})
	}
}
