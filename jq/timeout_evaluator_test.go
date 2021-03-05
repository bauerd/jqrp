package jq

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type mockEvaluator struct {
	results  []interface{}
	error    error
	duration time.Duration
}

func (m *mockEvaluator) Evaluate(_ string, _ interface{}) ([]interface{}, error) {
	time.Sleep(m.duration)
	if m.error != nil {
		return nil, m.error
	}
	return m.results, nil
}

func TestTimeoutEvaluatorNoTimeout(t *testing.T) {
	results := []interface{}{1, 2}
	mockEvaluator := mockEvaluator{results: results, duration: 1 * time.Second}
	evaluator := NewTimeoutEvaluator(&mockEvaluator, 2*time.Second)
	res, err := evaluator.Evaluate("", nil)
	if !reflect.DeepEqual(res, results) {
		t.Errorf("Unexpected results")
	}
	if err != nil {
		t.Errorf("Unexpected error")
	}
}

func TestTimeoutEvaluatorNoTimeoutError(t *testing.T) {
	error := errors.New("foobar")
	mockEvaluator := mockEvaluator{error: error, duration: 1 * time.Second}
	evaluator := NewTimeoutEvaluator(&mockEvaluator, 2*time.Second)
	res, err := evaluator.Evaluate("", nil)
	if res != nil {
		t.Errorf("Unexpected results")
	}
	if err != error {
		t.Errorf("Unexpected error")
	}
}

func TestTimeoutEvaluatorTimeout(t *testing.T) {
	results := []interface{}{1, 2}
	mockEvaluator := mockEvaluator{results: results, duration: 2 * time.Second}
	evaluator := NewTimeoutEvaluator(&mockEvaluator, 1*time.Second)
	res, err := evaluator.Evaluate("", nil)
	if reflect.DeepEqual(res, results) {
		t.Errorf("Unexpected results")
	}
	if err != ErrEvaluationTimeout {
		t.Errorf("Unexpected error")
	}
}
