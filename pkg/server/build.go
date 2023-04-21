package server

import (
	"net/http"
)

type Option func(*server)

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
