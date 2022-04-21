package transport

import (
	"context"
	"net/http"
	"time"
)

const (
	readTimeout  = time.Minute
	writeTimeout = time.Minute
)

type Server struct {
	httpServer *http.Server
}

func NewServer(address string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         address,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
