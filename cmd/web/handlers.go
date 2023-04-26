package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *app) Home(w http.ResponseWriter, r *http.Request) {
	// servemux "/" working like catch all
	// we need to check if path is really "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	p, err := template.ParseFiles(files...)

	if err != nil {
		app.errLogger.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = p.ExecuteTemplate(w, "base", nil)

	if err != nil {
		app.errLogger.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *app) SnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Hello username with id: %d", id)
}

func (app *app) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	// we need to check for method
	// because servemux doesnt support rest
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http error shortcut
		// for WriteHeader, Write
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`"snippet":"textTestText"`))
}
