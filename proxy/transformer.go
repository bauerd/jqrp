package proxy

import (
	"github.com/bauerd/jqrp/jq"
	"github.com/bauerd/jqrp/json"
	"mime"
	"net/http"
	"reflect"
)

// RawQueryHTTPHeader is the HTTP request header where the jq query is read
// from.
const RawQueryHTTPHeader string = "JQ"

type contextKey string

// RawQueryContextKey is the context key the jq query is stored under on
// requests.
const RawQueryContextKey contextKey = "RAW_QUERY"

// Transformer transforms some upstream responses.
type Transformer struct {
	evaluator jq.Evaluator
	rewriter  rewriter
}

// NewTransformer returns a new Transformer.
func NewTransformer(evaluator jq.Evaluator, rewriter rewriter) *Transformer {
	return &Transformer{
		evaluator: evaluator,
		rewriter:  rewriter,
	}
}

// ModifyResponse transforms responses to JSON requests that have the jq query
// header set.
func (t *Transformer) ModifyResponse(r *http.Response) error {
	// The context key is set only if (1) the request Accept'ed JSON,
	// and (2) a query was provided in the `JQ` HTTP header.
	// Otherwise, responses are proxied verbatim.
	rawQuery := r.Request.Context().Value(RawQueryContextKey)
	if rawQuery == nil {
		return nil
	}

	// Non-successful responses are proxied verbatim.
	if !(r.StatusCode >= 200 && r.StatusCode <= 299) {
		return nil
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return err
	}
	if mediaType != "application/json" {
		// The backend responded with a Content-Type != application/json,
		// but the client request only Accept'ed application/json.
		return ErrIllegalResponseType
	}

	input, err := json.Parse(r.Body)
	if err != nil {
		return ErrInvalidResponseBody
	}

	results, err := t.evaluator.Evaluate(rawQuery.(string), input)
	if err != nil {
		return err
	}

	if reflect.TypeOf(input).Kind() == reflect.Slice {
		return t.rewriter(results, []byte("[]"), r)
	}
	return t.rewriter(results, []byte("{}"), r)
}
