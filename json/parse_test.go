package json

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseValidJson(t *testing.T) {
	input := `{"valid": "json"}`
	reader := strings.NewReader(input)
	result, err := Parse(reader)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if !reflect.DeepEqual(result, map[string]interface{}{"valid": "json"}) {
		t.Errorf("Unexpected result: %v\n", result)
	}
}

func TestParseMalformedJson(t *testing.T) {
	input := `{"invalid: json"}`
	reader := strings.NewReader(input)
	result, err := Parse(reader)
	if err == nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if result != nil {
		t.Errorf("Unexpected result: %v\n", result)
	}
}

func TestParseMultipleRoot(t *testing.T) {
	input := `{"valid": "json"}[1, 2]`
	reader := strings.NewReader(input)
	result, err := Parse(reader)
	if err != ErrMultipleRoots {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if result != nil {
		t.Errorf("Unexpected result: %v\n", result)
	}
}

func TestParsePrimitiveRoot(t *testing.T) {
	input := `"foobar"`
	reader := strings.NewReader(input)
	result, err := Parse(reader)
	if err != ErrPrimitiveRootType {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if result != nil {
		t.Errorf("Unexpected result: %v\n", result)
	}
}
