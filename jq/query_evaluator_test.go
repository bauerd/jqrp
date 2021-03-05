package jq

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestEvaluatorMultipleResults(t *testing.T) {
	evaluator := NewQueryEvaluator(QueryCompiler)
	rawQuery := ".[] .id"
	rawInput := `
[
  {
    "id": 1,
    "name" :"alpha"
  },
  {
    "id": 2,
    "name": "beta"
  }
]
`
	var input interface{}
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBufferString(rawInput)))
	_ = decoder.Decode(&input)
	results, err := evaluator.Evaluate(rawQuery, input)
	if err != nil {
		t.Errorf("Evaluation failed")
	}
	eq := reflect.DeepEqual(results, []interface{}{1.0, 2.0})
	if !eq {
		t.Errorf("Unexpected results returned: %q\n", results)
	}
}

func TestEvaluatorSingleResult(t *testing.T) {
	evaluator := NewQueryEvaluator(QueryCompiler)
	rawQuery := "{id}"
	rawInput := `
{
  "id": 1,
  "name": "alpha"
}
`
	var input interface{}
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBufferString(rawInput)))
	_ = decoder.Decode(&input)
	results, err := evaluator.Evaluate(rawQuery, input)
	if err != nil {
		t.Errorf("Evaluation failed")
	}
	eq := reflect.DeepEqual(results, []interface{}{map[string]interface{}{"id": 1.0}})
	if !eq {
		t.Errorf("Unexpected results returned: %q\n", results)
	}
}

func TestEvaluatorInvalidQuery(t *testing.T) {
	evaluator := NewQueryEvaluator(QueryCompiler)
	rawQuery := "{!id}"
	rawInput := `
{
  "id": 1,
  "name": "alpha"
}
`
	var input interface{}
	decoder := json.NewDecoder(ioutil.NopCloser(bytes.NewBufferString(rawInput)))
	_ = decoder.Decode(&input)
	results, err := evaluator.Evaluate(rawQuery, input)
	if err == nil {
		t.Errorf("Evaluation did not fail")
	}
	if results != nil {
		t.Errorf("Unexpected results returned: %q\n", results)
	}
}
