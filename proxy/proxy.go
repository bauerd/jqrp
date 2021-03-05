package proxy

import (
	"github.com/bauerd/jqrp/jq"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy is a mutating reverse proxy.
type Proxy struct {
	backend *httputil.ReverseProxy
}

// NewProxy returns a new proxy that mutates upstream responses by using the
// given compiler
func NewProxy(url *url.URL, transport http.RoundTripper, evaluator jq.Evaluator) *Proxy {
	backend := httputil.NewSingleHostReverseProxy(url)
	transformer := NewTransformer(evaluator, Rewriter)

	// Preserve the director set by NewSingleHostReverseProxy
	backend.Director = Director(backend.Director)
	backend.ModifyResponse = transformer.ModifyResponse
	backend.ErrorHandler = ErrorHandler
	backend.Transport = transport

	return &Proxy{
		backend: backend,
	}
}

// ServeHTTP serves the proxy.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RequestID(HeaderParser(p.backend.ServeHTTP)).ServeHTTP(w, r)
}
