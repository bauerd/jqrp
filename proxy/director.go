package proxy

import (
	"github.com/bauerd/jqrp/log"
	"net/http"
)

// Director modifies incoming client requests.
func Director(super func(*http.Request)) func(*http.Request) {
	return func(r *http.Request) {
		super(r)
		r.Host = r.URL.Host
		log.Request(r)
	}
}
