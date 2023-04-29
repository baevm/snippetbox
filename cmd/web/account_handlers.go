package main

import (
	"errors"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
)

type UpdatePasswordForm struct {
	CurrentPassword         string `form:"currentPassword"`
	NewPassword             string `form:"newPassword"`
	NewPasswordConfirmation string `form:"newPasswordConfirmation"`
	validator.Validator     `form:"-"`
}

func (app *App) AccountView(w http.ResponseWriter, r *http.Request) {
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	user, err := app.users.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)

	data.Account = user

	app.render(w, http.StatusOK, "account.tmpl.html", data)
}

func (app *App) AccountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	var form UpdatePasswordForm
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	err := app.DecodePostForm(r, &form)

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// validating
	form.CheckField(validator.NotBlank(form.CurrentPassword), "currentPassword", "This field cant be empty")
	form.CheckField(validator.NotBlank(form.NewPassword), "newPassword", "This field cant be empty")
	form.CheckField(validator.NotBlank(form.NewPasswordConfirmation), "newPasswordConfirmation", "This field cant be empty")
	form.CheckField(validator.MinChars(form.NewPassword, 8), "newPassword", "New password cant be less then 8 characters")
	form.CheckField(validator.Equals(form.NewPassword, form.NewPasswordConfirmation), "newPasswordConfirmation", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "password.tmpl.html", data)
		return
	}

	err = app.users.PasswordUpdate(id, form.CurrentPassword, form.NewPassword)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddFieldError("currentPassword", "Incorrect current password")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "password.tmpl.html", data)
			return
		} else {
			app.serverError(w, err)
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Password has been changed successfully.")
	http.Redirect(w, r, "/account/view", http.StatusSeeOther)
}

func (app *App) AccountPasswordUpdateView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = UpdatePasswordForm{}
	app.render(w, http.StatusOK, "password.tmpl.html", data)
}
