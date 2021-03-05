package proxy

import (
	"bytes"
	"context"
	"errors"
	"github.com/bauerd/jqrp/jq"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type mockEvaluator struct {
	Result func() ([]interface{}, error)
	Called bool
}

func (m *mockEvaluator) Evaluate(_ string, _ interface{}) ([]interface{}, error) {
	m.Called = true
	return m.Result()
}

type mockRewriter struct {
	Called bool
	Result interface{}
}

func (r *mockRewriter) Rewrite(result []interface{}, _ []byte, _ *http.Response) error {
	r.Called = true
	r.Result = result
	return nil
}

func TestTransformerWithoutRawQuery(t *testing.T) {
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return nil, nil }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{}
	req := http.Request{}
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if evaluator.Called {
		t.Errorf("Evaluator called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}

func TestTransformerUnsuccessfulResponse(t *testing.T) {
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return nil, nil }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 400}
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if evaluator.Called {
		t.Errorf("Evaluator called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}

func TestTransformerWithoutJsonResponse(t *testing.T) {
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return nil, nil }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 200, Header: http.Header{}}
	res.Header.Set("Content-Type", "text/html")
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != ErrIllegalResponseType {
		t.Errorf("Unexpected error")
	}

	if evaluator.Called {
		t.Errorf("Evaluator called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}

func TestTransformerInvalidJsonResponse(t *testing.T) {
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return nil, nil }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 200, Header: http.Header{}}
	res.Header.Set("Content-Type", "application/json")
	body := `{"invalid: json"}`
	res.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != ErrInvalidResponseBody {
		t.Errorf("Unexpected error")
	}

	if evaluator.Called {
		t.Errorf("Evaluator called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}

func TestTransformerInvalidJsonResponsePrimitiveRoot(t *testing.T) {
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return nil, nil }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 200, Header: http.Header{}}
	res.Header.Set("Content-Type", "application/json")
	body := `"foobar"`
	res.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != ErrInvalidResponseBody {
		t.Errorf("Unexpected error")
	}

	if evaluator.Called {
		t.Errorf("Evaluator called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}

func TestTransformerInvalidJsonResponseMultipleRoots(t *testing.T) {
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return nil, nil }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 200, Header: http.Header{}}
	res.Header.Set("Content-Type", "application/json")
	body := `
{"alpha": 1}
{"beta": 2}
`
	res.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != ErrInvalidResponseBody {
		t.Errorf("Unexpected error")
	}

	if evaluator.Called {
		t.Errorf("Evaluator called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}

func TestTransformerEvaluationFailure(t *testing.T) {
	evalErr := errors.New("foobar")
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return nil, evalErr }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 200, Header: http.Header{}}
	res.Header.Set("Content-Type", "application/json")
	body := `{"valid": "json"}`
	res.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != evalErr {
		t.Errorf("Unexpected error %s", err)
	}

	if !evaluator.Called {
		t.Errorf("Evaluator not called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}

func TestTransformerResults(t *testing.T) {
	results := []interface{}{1, 2, 3}
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) { return results, nil }}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 200, Header: http.Header{}}
	res.Header.Set("Content-Type", "application/json")
	body := `{"valid": "json"}`
	res.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != nil {
		t.Errorf("Unexpected error")
	}

	if !evaluator.Called {
		t.Errorf("Evaluator not called")
	}

	if !rewriter.Called {
		t.Errorf("Rewriter not called")
	}

	if !reflect.DeepEqual(results, rewriter.Result) {
		t.Errorf("Rewriter called with unexpected results")
	}
}

func TestTransformerTimeout(t *testing.T) {
	evaluator := mockEvaluator{Result: func() ([]interface{}, error) {
		return nil, jq.ErrEvaluationTimeout
	}}
	rewriter := mockRewriter{}
	transformer := NewTransformer(&evaluator, rewriter.Rewrite)
	res := http.Response{StatusCode: 200, Header: http.Header{}}
	res.Header.Set("Content-Type", "application/json")
	body := `{"valid": "json"}`
	res.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	req := http.Request{}
	req = *req.WithContext(context.WithValue(req.Context(), RawQueryContextKey, "foobar"))
	res.Request = &req
	err := transformer.ModifyResponse(&res)

	if err != jq.ErrEvaluationTimeout {
		t.Errorf("Unexpected error: %s\n", err)
	}

	if !evaluator.Called {
		t.Errorf("Evaluator not called")
	}

	if rewriter.Called {
		t.Errorf("Rewriter called")
	}
}
