package proxy

import (
	"github.com/bauerd/jqrp/jq"
	"github.com/bauerd/jqrp/log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy is a mutating reverse proxy.
type Proxy struct {
	backend *httputil.ReverseProxy
	logger  *log.Logger
}

// NewProxy returns a new proxy that mutates upstream responses by using the
// given compiler
func NewProxy(url *url.URL, transport http.RoundTripper, evaluator jq.Evaluator, logger *log.Logger) *Proxy {
	backend := httputil.NewSingleHostReverseProxy(url)
	transformer := NewTransformer(evaluator, Rewriter(logger))

	// Preserve the director set by NewSingleHostReverseProxy
	backend.Director = Director(backend.Director, logger)
	backend.ModifyResponse = transformer.ModifyResponse
	backend.ErrorHandler = ErrorHandler(logger)
	backend.Transport = transport

	return &Proxy{
		backend: backend,
		logger:  logger,
	}
}

// ServeHTTP serves the proxy.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RequestID(HeaderParser(p.backend.ServeHTTP, p.logger)).ServeHTTP(w, r)
}
