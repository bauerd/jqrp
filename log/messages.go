package log

import (
	"fmt"
	"net/http"
)

// Request logs a client request.
func Request(logger *Logger, req *http.Request) {
	logger.Info(fmt.Sprintf("[%s] %s %s", requestID(req), req.Method, req.URL.Path))
}

// Query logs a client-supplied query.
func Query(logger *Logger, req *http.Request, rawQuery string) {
	logger.Info(fmt.Sprintf("[%s] Query: %s", requestID(req), rawQuery))
}

// SuccessResponse logs a successfully transformed response.
func SuccessResponse(logger *Logger, req *http.Request) {
	logger.Info(fmt.Sprintf("[%s] Rewriting succeeded", requestID(req)))
}

// FailureResponse logs a rewrite failure.
func FailureResponse(logger *Logger, req *http.Request, err error) {
	logger.Error(fmt.Sprintf("[%s] Rewriting failed: %s", requestID(req), err))
}

func requestID(req *http.Request) string {
	return req.Header.Get("X-Request-ID")
}
