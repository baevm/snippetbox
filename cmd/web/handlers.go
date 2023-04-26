package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
	"text/template"
)

func (app *app) Home(w http.ResponseWriter, r *http.Request) {
	// servemux "/" working like catch all
	// we need to check if path is really "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	_, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	p, err := template.ParseFiles(files...)

	if err != nil {
		app.serverError(w, err)
		return
	}

	err = p.ExecuteTemplate(w, "base", nil)

	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *app) SnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

func (app *app) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	// we need to check for method
	// because servemux doesnt support rest
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http error shortcut
		// for WriteHeader, Write
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Ultra title"
	content := "Hey content snippet LOL"
	expires := 7

	id, err := app.snippets.Create(title, content, expires)

	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
