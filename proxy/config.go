package proxy

import (
	"github.com/bauerd/jqrp/jq"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Config is the runtime configuration of jqrp.
type Config struct {
	Port                  int
	CacheSize             int
	EvaluationTimeout     time.Duration
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	DialTimeout           time.Duration
	DialKeepAlive         time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
}

// NewConfig returns a configuration read from environment variables.
func NewConfig() *Config {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	return &Config{
		Port:                  intFromEnvironment("PORT", 8989),
		CacheSize:             intFromEnvironment("CACHE_SIZE", 512),
		EvaluationTimeout:     durationFromEnvironment("EVAL_TIMEOUT", 0),
		ReadTimeout:           durationFromEnvironment("READ_TIMEOUT", 0),
		WriteTimeout:          durationFromEnvironment("WRITE_TIMEOUT", 0),
		DialTimeout:           durationFromEnvironment("DIAL_TIMEOUT", 0),
		DialKeepAlive:         durationFromEnvironment("DIAL_KEEPALIVE", 0),
		TLSHandshakeTimeout:   durationFromEnvironment("TLS_HANDSHAKE_TIMEOUT", 0),
		ResponseHeaderTimeout: durationFromEnvironment("RESPONSE_HEADER_TIMEOUT", 0),
		ExpectContinueTimeout: durationFromEnvironment("EXPECT_CONTINUE_TIMEOUT", 0),
	}
}

func intFromEnvironment(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fallback
		}
		return int(i)
	}
	return fallback
}

func durationFromEnvironment(key string, fallback int) time.Duration {
	return time.Duration(intFromEnvironment(key, fallback)) * time.Millisecond
}

// Transport returns an HTTP transport with timeouts set.
func (c *Config) Transport() *http.Transport {
	return &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   c.DialTimeout,
			KeepAlive: c.DialKeepAlive,
		}).Dial,
		TLSHandshakeTimeout:   c.TLSHandshakeTimeout,
		ResponseHeaderTimeout: c.ResponseHeaderTimeout,
		ExpectContinueTimeout: c.ExpectContinueTimeout,
	}
}

// Evaluator returns a configured evaluator.
func (c *Config) Evaluator() (jq.Evaluator, error) {
	compiler, err := c.compiler()
	if err != nil {
		return nil, err
	}
	if c.EvaluationTimeout <= 0 {
		return jq.NewQueryEvaluator(compiler), nil
	}
	return jq.NewTimeoutEvaluator(jq.NewQueryEvaluator(compiler), c.EvaluationTimeout), nil
}

func (c *Config) compiler() (jq.Compiler, error) {
	if c.CacheSize <= 0 {
		return jq.QueryCompiler, nil
	}
	cachedCompiler, err := jq.NewCachedCompiler(jq.QueryCompiler, c.CacheSize)
	if err != nil {
		return nil, err
	}
	return cachedCompiler.Compiler, nil
}
