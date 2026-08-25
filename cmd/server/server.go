package server

import (
	"net/http"
	"prasowka/cmd/conf"
	"time"
)

func NewServer(handler http.Handler) *http.Server {

	var cfg conf.Config
	cfg.New()

	port := cfg.Port

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s
}
