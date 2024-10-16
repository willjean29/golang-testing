package main

import (
	"html/template"
	"net/http"
	"path"
)

var pathTemplate = "./templates/"

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.html", &TemplateData{})
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	templParsed, err := template.ParseFiles(path.Join(pathTemplate, t))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}
	data.IP = app.ipFromContext(r.Context())
	err = templParsed.Execute(w, data)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}
	return nil
}
