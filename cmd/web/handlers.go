package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"strconv"
)

func (app *app) Home(w http.ResponseWriter, r *http.Request) {
	// servemux "/" working like catch all
	// we need to check if path is really "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
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

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
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
