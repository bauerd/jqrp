package proxy

import (
	"github.com/bauerd/jqrp/jq"
	"github.com/bauerd/jqrp/log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestProxyWithoutQueryHeader(t *testing.T) {
	const backendResponse = `{"valid": "json"}`
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, backendStatus; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if actual, expected := string(bodyBytes), backendResponse; actual != expected {
		t.Fatalf("Unexpected response body")
	}
}

func TestProxyWithoutAcceptHeader(t *testing.T) {
	const backendResponse = `{"valid": "json"}`
	const backendStatus = 201
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("JQ", "{valid}")
	res, err := frontendClient.Do(req)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if actual, expected := res.StatusCode, backendStatus; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if actual, expected := string(bodyBytes), backendResponse; actual != expected {
		t.Fatalf("Unexpected response body")
	}
}

func TestProxyInvalidQuery(t *testing.T) {
	const backendResponse = `
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
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", "!")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 400; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	expectedBody := ""
	if actual, expected := string(bodyBytes), expectedBody; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyArrayResult(t *testing.T) {
	const backendResponse = `
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
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".[] .id")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 203; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	expectedBody := "[1,2]"
	if actual, expected := string(bodyBytes), expectedBody; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyObjectResult(t *testing.T) {
	const backendResponse = `
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
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".[0] | {id}")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 203; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	expectedBody := `{"id":1}`
	if actual, expected := string(bodyBytes), expectedBody; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyPrimitiveResult(t *testing.T) {
	const backendResponse = `
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
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".[0] | .id")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 422; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	expectedBody := ""
	if actual, expected := string(bodyBytes), expectedBody; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyNullResult(t *testing.T) {
	const backendResponse = `
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
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".[0] | .foobar")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 422; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	expectedBody := ""
	if actual, expected := string(bodyBytes), expectedBody; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyNullArrayResult(t *testing.T) {
	const backendResponse = `
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
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".[] | .foobar")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 203; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	expectedBody := "[null,null]"
	if actual, expected := string(bodyBytes), expectedBody; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyEvaluationFailure(t *testing.T) {
	const backendResponse = `
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
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".foobar")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 400; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	expectedBody := ""
	if actual, expected := string(bodyBytes), expectedBody; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyInvalidContentType(t *testing.T) {
	const backendResponse = `{"valid": "json"}`
	const backendStatus = 201
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 502; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if actual, expected := string(bodyBytes), ""; actual != expected {
		t.Fatalf("Unexpected response body")
	}
}

func TestProxyMalformedResponseBody(t *testing.T) {
	const backendResponse = `{"broken: json"}`
	const backendStatus = 201
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", ".")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 502; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if actual, expected := string(bodyBytes), ""; actual != expected {
		t.Fatalf("Unexpected response body")
	}
}

func TestProxyEmptyResultObjectResponse(t *testing.T) {
	const backendResponse = `{"valid": "json"}`
	const backendStatus = 201
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", "empty")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 203; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if actual, expected := string(bodyBytes), "{}"; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyEmptyResultArrayResponse(t *testing.T) {
	const backendResponse = "[1, 2]"
	const backendStatus = 200
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(backendStatus)
		w.Write([]byte(backendResponse))
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("JQ", "empty")
	res, _ := frontendClient.Do(req)
	if actual, expected := res.StatusCode, 203; actual != expected {
		t.Errorf("Unexpected status code %d; expected %d", actual, expected)
	}
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	if actual, expected := string(bodyBytes), "[]"; actual != expected {
		t.Fatalf("Unexpected response body: %s", bodyBytes)
	}
}

func TestProxyXForwardedForHeader(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-For") == "" {
			t.Error("X-Forwarded-For header missing")
		}
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	frontendClient.Do(req)
}

func TestProxyXRequestIDHeader(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			t.Error("X-Request-ID header missing")
		}
	}))
	defer backend.Close()
	backendURL, _ := url.Parse(backend.URL)
	logger := log.New(log.Error)
	frontend := httptest.NewServer(NewProxy(backendURL, http.DefaultTransport, jq.NewQueryEvaluator(jq.QueryCompiler), logger))
	defer frontend.Close()
	frontendClient := frontend.Client()
	req, _ := http.NewRequest("GET", frontend.URL, nil)
	req.Header.Set("Accept", "application/json")
	frontendClient.Do(req)
}
