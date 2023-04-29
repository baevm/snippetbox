package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"snippetbox/internal/tests"
	"testing"
)

// unit
func Test_Ping(t *testing.T) {
	// response recorder
	rr := httptest.NewRecorder()

	// dummy http request
	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	// call handler with responseRecorder and request
	ping(rr, r)

	// get result
	rs := rr.Result()

	// compare result with expected
	tests.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	tests.Equal(t, string(body), "OK")
}

func Test_PingE2E(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")

	tests.Equal(t, status, http.StatusOK)
	tests.Equal(t, string(body), "OK")
}

func Test_SnippetView(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testCases := []struct {
		name    string
		url     string
		expCode int
		expBody string
	}{
		{
			name:    "Valid ID",
			url:     "/snippet/view/1",
			expCode: http.StatusOK,
			expBody: "Snippet Content",
		},
		{
			name:    "Negative ID",
			url:     "/snippet/view/-1",
			expCode: http.StatusNotFound,
		},
		{
			name:    "Float ID",
			url:     "/snippet/view/1.123145234",
			expCode: http.StatusNotFound,
		},
		{
			name:    "String ID",
			url:     "/snippet/view/test",
			expCode: http.StatusNotFound,
		},
		{
			name:    "Empty ID",
			url:     "/snippet/view/",
			expCode: http.StatusNotFound,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.url)

			tests.Equal(t, code, tt.expCode)

			if tt.expBody != "" {
				tests.StringContains(t, body, tt.expBody)
			}
		})
	}
}

func Test_SnippetCreate(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("Unauthenticated", func(t *testing.T) {
		status, header, _ := ts.get(t, "/snippet/create")

		tests.Equal(t, status, http.StatusSeeOther)
		tests.Equal(t, header.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")

		csrfToken := extractCsrfToken(t, body)

		const (
			Name     = "User"
			Password = "password"
			Email    = "user@test.com"
			formTag  = "<form action='/user/login' method='POST' novalidate>"
		)

		form := url.Values{}
		form.Add("name", Name)
		form.Add("password", Password)
		form.Add("email", Email)
		form.Add("csrf_token", csrfToken)

		ts.postForm(t, "/user/login", form)

		code, _, body := ts.get(t, "/snippet/create")

		tests.Equal(t, code, http.StatusOK)
		tests.StringContains(t, body, "<form action='/snippet/create' method='POST'>")
	})

}

func Test_UserSignup(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")

	csrfToken := extractCsrfToken(t, body)

	const (
		Name     = "User"
		Password = "password"
		Email    = "user@test.com"
		formTag  = "<form action='/user/signup' method='POST' novalidate>"
	)

	testCases := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		expCode      int
		expFormTag   string
	}{
		{
			name:         "Valid submission",
			userName:     Name,
			userEmail:    Email,
			userPassword: Password,
			csrfToken:    csrfToken,
			expCode:      http.StatusSeeOther,
		},
		{
			name:         "Invalid CSRF Token",
			userName:     Name,
			userEmail:    Email,
			userPassword: Password,
			csrfToken:    "wrongToken",
			expCode:      http.StatusBadRequest,
		},
		{
			name:         "Empty name",
			userName:     "",
			userEmail:    Email,
			userPassword: Password,
			csrfToken:    csrfToken,
			expCode:      http.StatusUnprocessableEntity,
			expFormTag:   formTag,
		},
		{
			name:         "Empty email",
			userName:     Name,
			userEmail:    "",
			userPassword: Password,
			csrfToken:    csrfToken,
			expCode:      http.StatusUnprocessableEntity,
			expFormTag:   formTag,
		},
		{
			name:         "Empty password",
			userName:     Name,
			userEmail:    Email,
			userPassword: "",
			csrfToken:    csrfToken,
			expCode:      http.StatusUnprocessableEntity,
			expFormTag:   formTag,
		},
		{
			name:         "Invalid email",
			userName:     Name,
			userEmail:    "bob@example.",
			userPassword: Password,
			csrfToken:    csrfToken,
			expCode:      http.StatusUnprocessableEntity,
			expFormTag:   formTag,
		},
		{
			name:         "Short password",
			userName:     Name,
			userEmail:    Email,
			userPassword: "pa$$",
			csrfToken:    csrfToken,
			expCode:      http.StatusUnprocessableEntity,
			expFormTag:   formTag,
		},
		{
			name:         "Duplicate email",
			userName:     Name,
			userEmail:    "dupe@example.com",
			userPassword: Password,
			csrfToken:    csrfToken,
			expCode:      http.StatusUnprocessableEntity,
			expFormTag:   formTag,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("password", tt.userPassword)
			form.Add("email", tt.userEmail)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			tests.Equal(t, code, tt.expCode)

			if tt.expFormTag != "" {
				tests.StringContains(t, body, tt.expFormTag)
			}
		})
	}
}
