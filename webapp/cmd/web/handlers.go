package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

var pathTemplate = "./templates/"

type TemplateData struct {
	IP   string
	Data map[string]any
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]any)
	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "Hit this page at "+time.Now().UTC().String())
	}
	_ = app.render(w, r, "home.page.html", &TemplateData{
		IP:   app.ipFromContext(r.Context()),
		Data: td,
	})
}

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "profile.page.html", &TemplateData{
		IP: app.ipFromContext(r.Context()),
	})
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.DB.GetUserByEmail(email)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Println(password, user.FirstName)

	_ = app.Session.RenewToken(r.Context())

	app.Session.Put(r.Context(), "flash", "Successfully logged in")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (app *application) render(w http.ResponseWriter, _ *http.Request, t string, data *TemplateData) error {
	templParsed, err := template.ParseFiles(path.Join(pathTemplate, t), path.Join(pathTemplate, "layout.page.html"))
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}
	err = templParsed.Execute(w, data)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return err
	}
	return nil
}
