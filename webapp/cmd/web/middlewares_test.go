package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
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

	var app application

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
	var app application
	const mockIP = "192.3.2.1"
	ctx := context.Background()
	ctx = context.WithValue(ctx, contentUserKey, mockIP)
	ip := app.ipFromContext(ctx)

	if ip != mockIP {
		t.Errorf("expected %s, got %s", mockIP, ip)
	}
}
