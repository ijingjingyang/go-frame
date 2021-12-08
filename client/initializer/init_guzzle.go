package initializer

import (
	"github.com/jjonline/go-lib-backend/guzzle"
	"net/http"
	"time"
)

//go:noinline
func initGuzzle() *guzzle.Client {
	return guzzle.New(&http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 10 * time.Second,
			DisableCompression:  true,
			MaxIdleConns:        400,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     50,
			IdleConnTimeout:     120 * time.Second,
		},
	})
}
