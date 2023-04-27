package server

import (
	"context"
	"net"
	"net/http"
)

type server struct {
	ln      *net.Listener
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
	return s.srv.Serve(*s.ln)
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
