package jq

import (
	"errors"
	"fmt"
)

// Evaluation errors.
var (
	// ErrEvaluationTimeout signals that evaluating a query exceeded the
	// evaluation timeout.
	ErrEvaluationTimeout = errors.New("query evaluation timed out")
)

// QueryEvaluationError signals that a query was malformed.
type QueryEvaluationError struct {
	Err error
}

// Error returns the wrapped error message.
func (e *QueryEvaluationError) Error() string {
	return fmt.Sprintf("query evaluation failed: %s", e.Err.Error())
}
