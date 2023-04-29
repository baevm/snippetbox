package main

import (
	"net/http"
	"snippetbox/ui"

	"github.com/go-chi/chi/v5"
)

func (app *app) routes() http.Handler {
	router := chi.NewRouter()

	// global middlewares
	router.Use(app.recoverPanic, app.logRequests, headerMiddleware)
	// custom not found
	router.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	}))

	router.HandleFunc("/ping", ping)

	// file server with embed filesystem that serves static content
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handle("/static/*", fileServer)

	router.Route("/user", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
		r.Get("/signup", app.UserSignup)
		r.Get("/login", app.UserLogin)
		r.Post("/signup", app.UserSignupPost)
		r.Post("/login", app.UserLoginPost)
		r.Post("/logout", app.UserLogoutPost)
	})

	router.Route("/account", func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate, app.requireAuth)

		r.Get("/view", app.AccountView)
		r.Get("/password/update", app.AccountPasswordUpdateView)
		r.Post("/password/update", app.AccountPasswordUpdate)
	})

	// routes with session middleware
	router.Group(func(r chi.Router) {
		r.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
		r.Get("/", app.Home)
		r.Get("/snippet/view/{id}", app.SnippetView)
		r.Get("/about", app.AboutView)

		// with auth middleware
		r.With(app.requireAuth).Get("/snippet/create", app.SnippetCreate)
		r.With(app.requireAuth).Post("/snippet/create", app.SnippetCreatePost)
	})

	return router
}
