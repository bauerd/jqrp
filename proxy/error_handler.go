package proxy

import (
	"errors"
	"github.com/bauerd/jqrp/jq"
	"github.com/bauerd/jqrp/log"
	"net/http"
)

// ErrorHandler writes the response status code in case of errors.
func ErrorHandler(responseWriter http.ResponseWriter, req *http.Request, err error) {
	var e *jq.QueryEvaluationError
	if errors.As(err, &e) {
		responseWriter.WriteHeader(400)
		log.FailureResponse(req, err)
		return
	}

	switch err {
	case ErrInvalidResponseBody:
		responseWriter.WriteHeader(502)
		log.FailureResponse(req, err)
	case ErrIllegalResponseType:
		responseWriter.WriteHeader(502)
		log.FailureResponse(req, err)
	case ErrIllegalQueryResult:
		responseWriter.WriteHeader(422)
		log.FailureResponse(req, err)
	case jq.ErrEvaluationTimeout:
		responseWriter.WriteHeader(408)
		log.FailureResponse(req, err)
	default:
		responseWriter.WriteHeader(500)
		log.FailureResponse(req, err)
	}
}
