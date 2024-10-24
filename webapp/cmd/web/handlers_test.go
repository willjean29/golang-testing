package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setup() {
	pathTemplate = "./../../templates/"
}

func Test_application_handlers(t *testing.T) {
	setup()
	var testRoutes = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/finish", http.StatusNotFound},
	}

	router := app.routes()

	ts := httptest.NewTLSServer(router)
	defer ts.Close()

	for _, tr := range testRoutes {
		resp, err := ts.Client().Get(ts.URL + tr.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != tr.expectedStatusCode {
			t.Errorf("expected status code %d, got %d", tr.expectedStatusCode, resp.StatusCode)
		}
	}
}

func Test_application_render(t *testing.T) {
	setup()

	tests := []struct {
		name          string
		templateName  string
		templateData  *TemplateData
		expectedError bool
	}{
		{"valid template", "home.page.html", &TemplateData{}, false},
		{"invalid template", "nonexistent.page.html", &TemplateData{}, true},
		{"invalid data", "test.page.html", &TemplateData{}, true},
	}

	for _, tt := range tests {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)

		err := app.render(res, req, tt.templateName, tt.templateData)
		if tt.expectedError {
			if err == nil {
				t.Errorf("%s: expected error but got none", tt.name)
			}
		} else {
			if err != nil {
				t.Errorf("%s: did not expect error but got %v", tt.name, err)
			}
			if res.Code != http.StatusOK {
				t.Errorf("%s: expected status code %d, got %d", tt.name, http.StatusOK, res.Code)
			}
		}

	}
}

func Test_application_home(t *testing.T) {
	setup()
	tests := []struct {
		name         string
		contentKey   string
		contentValue string
		statusCode   int
		expected     string
	}{
		{"exist session", "test", "jean@gmail", http.StatusOK, "From session: jean@gmail"},
		{"no exist session", "", "unknown", http.StatusOK, "From session:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req = addContextAndSessionToRequest(req)
			app.Session.Put(req.Context(), tt.contentKey, tt.contentValue)
			res := httptest.NewRecorder()

			app.Home(res, req)

			if res.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, res.Code)
			}

			if !strings.Contains(res.Body.String(), tt.expected) {
				t.Errorf("expected %q, got %q", tt.expected, res.Body.String())
			}
		})
	}

}

func Test_application_login(t *testing.T) {
	tests := []struct {
		name           string
		body           io.Reader
		expected       string
		expectedStatus int
	}{
		{"valid form", strings.NewReader(url.Values{"email": {"jean@gmail"}, "password": {"123456"}}.Encode()), "Email: jean@gmail", http.StatusOK},
		{"invalid form", strings.NewReader(url.Values{"email": {"jean@gmail"}}.Encode()), "Form is not valid", http.StatusOK},
		{"invalid body", strings.NewReader("%"), "Bad Request", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", tt.body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			res := httptest.NewRecorder()

			app.Login(res, req)

			if res.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.Code)
			}

			if strings.TrimSpace(res.Body.String()) != strings.TrimSpace(tt.expected) {
				t.Errorf("expected %q, got %q", tt.expected, res.Body.String())
			}
		})
	}
}

func getCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, contentUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request) *http.Request {
	req = req.WithContext(getCtx())
	session, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))
	return req.WithContext(session)
}
