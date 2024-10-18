package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_application_handlers(t *testing.T) {
	var testRoutes = []struct {
		name               string
		url                string
		expectedStatusCode int
	}{
		{"home", "/", http.StatusOK},
		{"404", "/finish", http.StatusNotFound},
	}

	var app application
	app = application{}

	router := app.routes()

	ts := httptest.NewTLSServer(router)
	defer ts.Close()
	pathTemplate = "./../../templates/"
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
	var app application
	app = application{}

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
