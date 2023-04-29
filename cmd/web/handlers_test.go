package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"snippetbox/internal/assert"
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
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}

func Test_PingE2E(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/ping")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func Test_SnippetView(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.url)

			assert.Equal(t, code, tt.expCode)

			if tt.expBody != "" {
				assert.StringContains(t, body, tt.expBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
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

	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("password", tt.userPassword)
			form.Add("email", tt.userEmail)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.expCode)

			if tt.expFormTag != "" {
				assert.StringContains(t, body, tt.expFormTag)
			}
		})
	}
}
