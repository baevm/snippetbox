package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	router := chi.NewRouter()

	// custom not found
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	}))

	// file server that serves static content
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})

	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	router.Get("/", app.Home)
	router.Get("/snippet/view/{id}", app.SnippetView)
	router.Get("/snippet/create", app.SnippetCreate)
	router.Post("/snippet/create", app.SnippetCreatePost)

	middlewares := alice.New(app.recoverPanic, app.logRequests, headerMiddleware)

	return middlewares.Then(router)
}
