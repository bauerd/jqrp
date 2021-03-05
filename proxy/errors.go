package proxy

import "errors"

// Upstream errors.
var (
	// ErrInvalidResponseBody signals that an upstream response contained invalid
	// JSON.
	ErrInvalidResponseBody = errors.New("upstream response is invalid JSON")

	// ErrIllegalResponseType signals that an upstream response's Content-Type was
	// not application/json.
	ErrIllegalResponseType = errors.New("upstream response is not application/json")

	// ErrIllegalQueryResult signals that a query resulted in a result type that
	// has no JSON representation on its own.
	ErrIllegalQueryResult = errors.New("query resulted in primitive type")
)
