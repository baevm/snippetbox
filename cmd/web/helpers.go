package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"snippetbox/internal/templates"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *App) serverError(w http.ResponseWriter, error error) {
	// get stack trace for current goroutine
	trace := fmt.Sprintf("%s\n%s", error.Error(), debug.Stack())
	app.errLogger.Print(2, trace)

	if app.debug {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *App) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *App) render(w http.ResponseWriter, status int, page string, data *templates.TemplateData) {
	ts, ok := app.templateCache[page]

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

func (app *App) newTemplateData(r *http.Request) *templates.TemplateData {
	return &templates.TemplateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *App) DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()

	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)

	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
	}

	return nil
}

func (app *App) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)

	if !ok {
		return false
	}

	return isAuthenticated
}
