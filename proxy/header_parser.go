package proxy

import (
	"context"
	"github.com/bauerd/jqrp/log"
	"mime"
	"net/http"
)

// HeaderParser attaches the JQ request header value as a context key on client
// requests.
var HeaderParser = func(f http.HandlerFunc, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
		if err != nil || mediaType != "application/json" {
			f(w, r)
			return
		}

		rawQuery := r.Header.Get(string(RawQueryHTTPHeader))
		if rawQuery == "" {
			f(w, r)
			return
		}
		log.Query(logger, r, rawQuery)
		f(w, r.WithContext(context.WithValue(r.Context(), RawQueryContextKey, rawQuery)))
	}
}
