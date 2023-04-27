package server

import (
	"net"
	"net/http"
)

type Option func(*server)

func WithListener(ln *net.Listener) Option {
	return func(s *server) {
		s.ln = ln
	}
}

func WithSrv(srv *http.Server) Option {
	return func(s *server) {
		s.srv = srv
	}
}

func WithHandler(handler http.Handler) Option {
	return func(s *server) {
		s.handler = handler
	}
}
