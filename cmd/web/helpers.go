package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *app) serverError(w http.ResponseWriter, error error) {
	// get stack trace for current goroutine
	trace := fmt.Sprintf("%s\n%s", error.Error(), debug.Stack())
	app.errLogger.Print(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *app) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *app) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *app) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templaceCache[page]

	if !ok {
		app.serverError(w, fmt.Errorf("template %s doesnt exist", page))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *app) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
