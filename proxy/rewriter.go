package proxy

import (
	"bytes"
	"encoding/json"
	"github.com/bauerd/jqrp/log"
	"io/ioutil"
	"net/http"
	"reflect"
)

type rewriter = func([]interface{}, []byte, *http.Response) error

// Rewriter rewrites response bodies depending on query results.
func Rewriter(logger *log.Logger) rewriter {
	return func(results []interface{}, fallbackBody []byte, response *http.Response) error {
		resultsLen := len(results)

		if resultsLen == 0 {
			return writeRawBody(fallbackBody, response, logger)
		}

		if resultsLen > 1 {
			return writeJSONBody(results, response, logger)
		}

		// If the only result is null, it has no JSON representation.
		if results[0] == nil {
			return ErrIllegalQueryResult
		}

		// If the only result is of another primitive type, it has no JSON
		// representation.
		switch reflect.TypeOf(results[0]).Kind() {
		case reflect.Slice:
		case reflect.Map:
			break
		default:
			return ErrIllegalQueryResult
		}

		return writeJSONBody(results[0], response, logger)
	}
}

func writeJSONBody(payload interface{}, response *http.Response, logger *log.Logger) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return writeRawBody(body, response, logger)
}

func writeRawBody(payload []byte, response *http.Response, logger *log.Logger) error {
	response.Body.Close()
	response.Body = ioutil.NopCloser(bytes.NewReader(payload))
	defer response.Body.Close() // no-op
	response.ContentLength = int64(len(payload))
	response.StatusCode = 203
	log.SuccessResponse(logger, response.Request)
	return nil
}
