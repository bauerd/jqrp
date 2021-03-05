package main

import (
	"fmt"
	"github.com/bauerd/jqrp/log"
	"github.com/bauerd/jqrp/proxy"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	frontend := proxy.NewProxy(url, config.Transport(), evaluator)

	log.ConfigValue("URL", url.String())
	log.ConfigValue("Port", strconv.Itoa(config.Port))
	log.ConfigValue("Query evaluation timeout", config.EvaluationTimeout.String())
	log.ConfigValue("Query cache size", strconv.Itoa(config.CacheSize))
	log.ConfigValue("Frontend read timeout", config.ReadTimeout.String())
	log.ConfigValue("Frontend write timeout", config.WriteTimeout.String())
	log.ConfigValue("Backend TCP dial timeout", config.DialTimeout.String())
	log.ConfigValue("Backend TLS handshake timout", config.TLSHandshakeTimeout.String())
	log.ConfigValue("Backend response header timeout", config.ResponseHeaderTimeout.String())
	log.ConfigValue("Backend 100-continue timeout", config.ExpectContinueTimeout.String())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      frontend,
		ReadTimeout:  time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.WriteTimeout) * time.Millisecond,
	}
	server.ListenAndServe().Error()
}
