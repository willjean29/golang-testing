package main

import (
	"html/template"
	"net/http"
)

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}

func (app *application) render(w http.ResponseWriter, _ *http.Request, t string, data *TemplateData) error {
	templParsed, err := template.ParseFiles("./templates/" + t)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}
	err = templParsed.Execute(w, data)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	return nil
}
