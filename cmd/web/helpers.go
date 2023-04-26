package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
