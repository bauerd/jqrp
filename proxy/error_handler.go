package proxy

import (
	"errors"
	"github.com/bauerd/jqrp/jq"
	"github.com/bauerd/jqrp/log"
	"net/http"
)

// ErrorHandler writes the response status code in case of errors.
func ErrorHandler(logger *log.Logger) func(http.ResponseWriter, *http.Request, error) {
	return func(responseWriter http.ResponseWriter, req *http.Request, err error) {
		var e *jq.QueryEvaluationError
		if errors.As(err, &e) {
			responseWriter.WriteHeader(400)
			log.FailureResponse(logger, req, err)
			return
		}

		switch err {
		case ErrInvalidResponseBody:
			responseWriter.WriteHeader(502)
			log.FailureResponse(logger, req, err)
		case ErrIllegalResponseType:
			responseWriter.WriteHeader(502)
			log.FailureResponse(logger, req, err)
		case ErrIllegalQueryResult:
			responseWriter.WriteHeader(422)
			log.FailureResponse(logger, req, err)
		case jq.ErrEvaluationTimeout:
			responseWriter.WriteHeader(408)
			log.FailureResponse(logger, req, err)
		default:
			responseWriter.WriteHeader(500)
			log.FailureResponse(logger, req, err)
		}
	}
}
