package server

import (
	"context"
	"net/http"
)

type server struct {
	srv     *http.Server
	handler http.Handler
}

func NewServer(opts ...Option) *server {
	s := &server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
