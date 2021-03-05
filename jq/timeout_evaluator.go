package jq

import (
	"time"
)

// TimeoutEvaluator wraps an evaluator and sets a timeout on its execution.
type TimeoutEvaluator struct {
	evaluator Evaluator // TODO: Use embedding
	timeout   time.Duration
}

type timeoutEvaluationResult struct {
	results []interface{}
	error   error
}

// NewTimeoutEvaluator returns a new evaluator with a timeout on its execution.
func NewTimeoutEvaluator(evaluator Evaluator, timeout time.Duration) *TimeoutEvaluator {
	return &TimeoutEvaluator{
		evaluator: evaluator,
		timeout:   timeout,
	}
}

// Evaluate evaluates a raw query, erroring if a time limit is exceeded.
func (e *TimeoutEvaluator) Evaluate(rawQuery string, input interface{}) ([]interface{}, error) {
	evalResult := make(chan timeoutEvaluationResult, 1)
	go func() {
		results, err := e.evaluator.Evaluate(rawQuery, input)
		if err != nil {
			evalResult <- timeoutEvaluationResult{error: err}
		} else {
			evalResult <- timeoutEvaluationResult{results: results}
		}
	}()

	var results []interface{}
	select {
	case result := <-evalResult:
		if result.error != nil {
			return nil, result.error
		}
		results = result.results
	case <-time.After(e.timeout):
		return nil, ErrEvaluationTimeout
	}

	return results, nil
}
