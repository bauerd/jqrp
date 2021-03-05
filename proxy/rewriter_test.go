package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestRewriterEmptyMultipleResults(t *testing.T) {
	url, _ := url.Parse("https://example.com")
	req := http.Request{URL: url}
	res := http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{})), Request: &req}
	results := []interface{}{1, 2}
	err := Rewriter(results, []byte{}, &res)
	if err != nil {
		t.Errorf("Rewriting failed")
	}
	if res.ContentLength != 5 {
		t.Errorf("Wrong Content-Length")
	}
	expectedBody := []byte("[1,2]")
	body, _ := ioutil.ReadAll(res.Body)
	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Wrong response body")
	}
}

func TestRewriterEmptySingleComplexResult(t *testing.T) {
	url, _ := url.Parse("https://example.com")
	req := http.Request{URL: url}
	res := http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{})), Request: &req}
	results := []interface{}{map[string]string{"hello": "world"}}
	err := Rewriter(results, []byte{}, &res)
	if err != nil {
		t.Errorf("Rewriting failed")
	}
	if res.ContentLength != 17 {
		t.Errorf("Wrong Content-Length")
	}
	expectedBody := []byte(`{"hello":"world"}`)
	body, _ := ioutil.ReadAll(res.Body)
	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Wrong response body")
	}
}

func TestRewriterEmptySinglePrimitiveResult(t *testing.T) {
	url, _ := url.Parse("https://example.com")
	req := http.Request{URL: url}
	res := http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{})), Request: &req}
	results := []interface{}{1}
	err := Rewriter(results, []byte{}, &res)
	if err != ErrIllegalQueryResult {
		t.Errorf("Rewriting succeeded")
	}
	if res.ContentLength != 0 {
		t.Errorf("Wrong Content-Length")
	}
	expectedBody := []byte{}
	body, _ := ioutil.ReadAll(res.Body)
	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Wrong response body")
	}
}

func TestRewriterEmptyWithoutResult(t *testing.T) {
	url, _ := url.Parse("https://example.com")
	req := http.Request{URL: url}
	res := http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte{})), Request: &req}
	results := []interface{}{}
	err := Rewriter(results, []byte("[]"), &res)
	if err != nil {
		t.Errorf("Rewriting failed")
	}
	if res.ContentLength != 2 {
		t.Errorf("Wrong Content-Length")
	}
	expectedBody := []byte("[]")
	body, _ := ioutil.ReadAll(res.Body)
	if !reflect.DeepEqual(body, expectedBody) {
		t.Errorf("Wrong response body: %s", body)
	}
}
