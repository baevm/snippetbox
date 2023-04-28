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

	router.Route("/user", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
		r.Get("/signup", app.UserSignup)
		r.Get("/login", app.UserLogin)
		r.Post("/signup", app.UserSignupPost)
		r.Post("/login", app.UserLoginPost)
		r.Post("/logout", app.UserLogoutPost)
	})

	// routes with session middleware
	router.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
		r.Get("/", app.Home)
		r.Get("/snippet/view/{id}", app.SnippetView)

		// require auth for creating snippets
		r.With(app.requireAuth).Get("/snippet/create", app.SnippetCreate)
		r.With(app.requireAuth).Post("/snippet/create", app.SnippetCreatePost)
	})

	return router
}
