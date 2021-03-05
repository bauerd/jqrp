package log

import (
	zerolog "github.com/rs/zerolog/log"
	"net/http"
)

// Request logs a client request.
func ConfigValue(key string, val string) {
	zerolog.Debug().Str(key, val)
}

func ServerError(err error) {
	zerolog.Error().
		Err(err).
		Msg("Server error")
}

// Request logs a client request.
func Request(req *http.Request) {
	zerolog.Debug().
		Str("RequestID", requestID(req)).
		Str("Method", req.Method).
		Str("Path", req.URL.Path).
		Msg("Client request")
}

// Query logs a client-supplied query.
func Query(req *http.Request, rawQuery string) {
	zerolog.Debug().
		Str("RequestID", requestID(req)).
		Str("Query", rawQuery).
		Msg("Query")
}

// SuccessResponse logs a successfully transformed response.
func SuccessResponse(req *http.Request) {
	zerolog.Info().
		Str("RequestID", requestID(req)).
		Msg("Rewriting successful")
}

// FailureResponse logs a rewrite failure.
func FailureResponse(req *http.Request, err error) {
	zerolog.Info().
		Str("RequestID", requestID(req)).
		Err(err).
		Msg("Rewriting failed")
}

func requestID(req *http.Request) string {
	return req.Header.Get("X-Request-ID")
}
