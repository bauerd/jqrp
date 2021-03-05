package proxy

import (
	"net/http"
	"testing"
)

func TestHeaderParserWithoutAcceptHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set(RawQueryHTTPHeader, "foobar")
	handler := HeaderParser(func(_ http.ResponseWriter, r *http.Request) {
		rawQuery := r.Context().Value(RawQueryContextKey)
		if rawQuery != nil {
			t.Errorf("Context value is %s", rawQuery)
		}
	})
	handler.ServeHTTP(nil, req)
}

func TestHeaderParserWithQueryHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set(RawQueryHTTPHeader, "foobar")
	handler := HeaderParser(func(_ http.ResponseWriter, r *http.Request) {
		rawQuery := r.Context().Value(RawQueryContextKey)
		if rawQuery != "foobar" {
			t.Errorf("Context value is %s", rawQuery)
		}
	})
	handler.ServeHTTP(nil, req)
}

func TestHeaderParserWithoutQueryHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	handler := HeaderParser(func(_ http.ResponseWriter, r *http.Request) {
		rawQuery := r.Context().Value(RawQueryContextKey)
		if rawQuery != nil {
			t.Errorf("Context value is %s", rawQuery)
		}
	})
	handler.ServeHTTP(nil, req)
}
