package main

import (
	"fmt"
	"github.com/bauerd/jqrp/proxy"
	"net/http"
	"net/url"
	"os"
	"time"
)

const usage string = "Usage: jqrp BACKEND"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	rawURL := os.Args[1]
	url, err := url.Parse(rawURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse URL %s", rawURL)
		os.Exit(1)
	}

	config := proxy.NewConfig()
	evaluator, err := config.Evaluator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to allocate compiler")
		os.Exit(1)
	}
	logger := config.Logger()
	frontend := proxy.NewProxy(url, config.Transport(), evaluator, logger)

	logger.Debug(fmt.Sprintf("URL: %s", url))
	logger.Debug(fmt.Sprintf("Port: %d", config.Port))
	logger.Debug(fmt.Sprintf("Query evaluation timeout: %s", config.EvaluationTimeout))
	logger.Debug(fmt.Sprintf("Query cache size: %d", config.CacheSize))
	logger.Debug(fmt.Sprintf("Frontend read timeout: %s", config.ReadTimeout))
	logger.Debug(fmt.Sprintf("Frontend write timeout: %s", config.WriteTimeout))
	logger.Debug(fmt.Sprintf("Backend TCP dial timeout: %s", config.DialTimeout))
	logger.Debug(fmt.Sprintf("Backend TLS handshake timout: %s", config.TLSHandshakeTimeout))
	logger.Debug(fmt.Sprintf("Backend response header timeout: %s", config.ResponseHeaderTimeout))
	logger.Debug(fmt.Sprintf("Backend 100-continue timeout: %s", config.ExpectContinueTimeout))
	logger.Debug(fmt.Sprintf("Log level: %s", logger.Level))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      frontend,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Millisecond,
	}
	logger.Error(server.ListenAndServe().Error())
}
