package proxy

import (
	"testing"
	"github.com/bauerd/jqrp/jq"
)

func TestConfigEvaluatorWithoutTimeout(t *testing.T) {
	config := Config{EvaluationTimeout: 0}
	evaluator, _ := config.Evaluator()
	switch evaluator.(type) {
	case *jq.QueryEvaluator:
		return
	default:
		t.Error("Unexpected evaluator")
	}
}

func TestConfigEvaluatorWithTimeout(t *testing.T) {
	config := Config{EvaluationTimeout: 1}
	evaluator, _ := config.Evaluator()
	switch evaluator.(type) {
	case *jq.TimeoutEvaluator:
		return
	default:
		t.Error("Unexpected evaluator")
	}
}
