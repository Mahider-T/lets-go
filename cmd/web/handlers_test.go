package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"oogway/first/snippetbox/internal/assert"
	"testing"
)

func TestPing(t *testing.T) {
	record := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	ping(record, r)

	rs := record.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}

func TestPingEndToEnd(t *testing.T) {

	app := newTestApplication(t)

	ts := newTestServer(t, app.routers())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
