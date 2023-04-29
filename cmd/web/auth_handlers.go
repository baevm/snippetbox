package main

import (
	"errors"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
)

type UserSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *App) UserSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserSignupForm{}
	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *App) UserSignupPost(w http.ResponseWriter, r *http.Request) {
	var form UserSignupForm
	err := app.DecodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validating
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cant be empty")
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cant be empty")
	form.CheckField(validator.NotBlank(form.Password), "passwowrd", "This field cant be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be valid email")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field cant be less than 8 characters length")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Create(form.Name, form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "User succesfully created!")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *App) UserLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UserLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *App) UserLoginPost(w http.ResponseWriter, r *http.Request) {
	var form UserLoginForm
	err := app.DecodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validating
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cant be empty")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cant be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", "This field must be valid email")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field cant be less than 8 characters length")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	// auth
	id, err := app.users.Authenticate(form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// renew session ID
	err = app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, err)
		return
	}

	// add current user to session
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	url := app.sessionManager.PopString(r.Context(), "redirectURL")

	if url == "" {
		http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (app *App) UserLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
