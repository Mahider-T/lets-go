package main

import (
	"bytes"
	"io"
	"log"
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

	app := &application{
		errorLogger: log.New(io.Discard, "", 0),
		infoLogger:  log.New(io.Discard, "", 0),
	}

	ts := httptest.NewTLSServer(app.routers())
	defer ts.Close()

	rs, err := ts.Client().Get((ts.URL + "/ping"))

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
