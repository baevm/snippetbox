package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
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
