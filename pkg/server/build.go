package server

import (
	"net"
	"net/http"
)

type Option func(*server)

func WithSrv(srv *http.Server) Option {
	return func(s *server) {
		s.srv = srv
	}
}

func WithListener(listener net.Listener) Option {
	return func(s *server) {
		s.listener = listener
	}
}
