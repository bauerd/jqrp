package proxy

import (
	"net/http"
	"testing"
)

func TestDirector(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Host = "alpha.com"
	req.URL.Host = "beta.com"

	Director(func(*http.Request) {})(req)

	if req.Host != req.URL.Host {
		t.Errorf("Mismatching host %s", req.Host)
	}

	if req.URL.Host != "beta.com" {
		t.Errorf("Mismatching host %s", req.Host)
	}
}
