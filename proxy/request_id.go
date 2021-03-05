package proxy

import (
	"github.com/google/uuid"
	"net/http"
)

// RequestID sets the X-Request-ID header to a UUID if not already set.
var RequestID = func(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			r.Header.Set("X-Request-ID", uuid.NewString())
		}
		f(w, r)
	}
}
