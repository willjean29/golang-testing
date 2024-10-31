package main

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
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
		{"profile", "/user/profile", http.StatusOK},
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
		req = addContextAndSessionToRequest(req)
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
		name             string
		body             io.Reader
		expectedStatus   int
		expectedLocation string
	}{
		{"valid credentials", strings.NewReader(url.Values{"email": {"admin@example.com"}, "password": {"secret"}}.Encode()), http.StatusSeeOther, "/user/profile"},
		{"valid credentials (password)", strings.NewReader(url.Values{"email": {"admin@example.com"}, "password": {"secret123"}}.Encode()), http.StatusSeeOther, "/"},
		{"invalid credential (email)", strings.NewReader(url.Values{"email": {"admin2@example.com"}, "password": {"secret"}}.Encode()), http.StatusSeeOther, "/"},
		{"invalid form", strings.NewReader(url.Values{"email": {"jean@gmail"}}.Encode()), http.StatusSeeOther, "/"},
		{"invalid body", strings.NewReader("%"), http.StatusBadRequest, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", tt.body)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req = addContextAndSessionToRequest(req)
			res := httptest.NewRecorder()

			app.Login(res, req)

			if res.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.Code)
			}

			location, err := res.Result().Location()

			if err == nil {
				if location.String() != tt.expectedLocation {
					t.Errorf("expected location %q, got %q", tt.expectedLocation, location.String())
				}
			} else {
				if tt.expectedLocation != "" {
					t.Errorf("expected location %q, got %v", tt.expectedLocation, err)
				}
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

func Test_app_uploadFiles(t *testing.T) {
	// setup pipe
	r, w := io.Pipe()

	// create a new writer
	writer := multipart.NewWriter(w)

	// create a wait group, and add 1 to it
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// simulate a file upload by writing a file to the writer
	go simulatePngUpload("./../../static/img/test/img.png", writer, t, wg)

	// read from the reader end of the pipe
	request := httptest.NewRequest("POST", "/upload", r)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	// call app.uploadFiles
	uploadFiles, err := app.UploadFiles(request, "./../../static/img/")
	if err != nil {
		t.Error(err)
	}
	// perform test
	if _, err := os.Stat(fmt.Sprintf("./../../static/img/%s", uploadFiles[0].OriginalFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exist: %s", err.Error())
	}
	// clean up
	_ = os.Remove(fmt.Sprintf("./../../static/img/%s", uploadFiles[0].OriginalFileName))
}

func simulatePngUpload(fileToUpload string, writer *multipart.Writer, t *testing.T, wg *sync.WaitGroup) {
	defer writer.Close()
	defer wg.Done()

	// create a new form data field "file" with value being filename
	part, err := writer.CreateFormFile("file", path.Base(fileToUpload))
	if err != nil {
		t.Error(err)
	}

	// open the file
	file, err := os.Open(fileToUpload)
	if err != nil {
		t.Error(err)
	}

	// decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		t.Error("Error decoded image :", err)
	}

	// write the image to our io.Writer
	err = png.Encode(part, img)
}
