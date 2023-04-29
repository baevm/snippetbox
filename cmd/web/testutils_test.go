package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"snippetbox/internal/models/mocks"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

type testServer struct {
	*httptest.Server
}

var csrfTokenMock = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

// loggers needed for middlewares
// other way it will result in panic
func newTestApp(t *testing.T) *app {
	templateCache, err := newTemplateCache()

	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &app{
		errLogger:      log.New(io.Discard, "", 0),
		infoLogger:     log.New(io.Discard, "", 0),
		users:          &mocks.UserModel{},
		snippets:       &mocks.SnippetModel{},
		templaceCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// mock tls server
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// cookie jar
	jar, err := cookiejar.New(nil)

	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	// disable redirect for test server
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func extractCsrfToken(t *testing.T, body string) string {
	matches := csrfTokenMock.FindStringSubmatch(body)

	if len(matches) < 2 {
		t.Fatal("No csrf token found")
	}

	return html.UnescapeString(string(matches[1]))
}

// mock GET
func (ts *testServer) get(t *testing.T, url string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + url)

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

// mock POST for forms
func (ts *testServer) postForm(t *testing.T, url string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+url, form)

	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
