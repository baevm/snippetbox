package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *app) routes() http.Handler {
	router := chi.NewRouter()

	// middlewares
	router.Use(app.recoverPanic, app.logRequests, headerMiddleware)
	// custom not found
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	}))

	// file server that serves static content
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// routes with session middleware
	router.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave)
		r.Get("/", app.Home)
		r.Get("/snippet/view/{id}", app.SnippetView)
		r.Get("/snippet/create", app.SnippetCreate)
		r.Post("/snippet/create", app.SnippetCreatePost)
	})

	return router
}
