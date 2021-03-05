package proxy

import (
	"net/http"
	"testing"
)

func TestRequestIDWithoutHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	handler := RequestID(func(_ http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			t.Error("X-Request-ID missing")
		}
	})
	handler.ServeHTTP(nil, req)
}

func TestRequestIDWithHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "foobar")
	handler := RequestID(func(_ http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") != "foobar" {
			t.Error("Unexpected X-Request-ID")
		}
	})
	handler.ServeHTTP(nil, req)
}
