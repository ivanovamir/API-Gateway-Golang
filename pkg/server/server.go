package server

import (
	"net"
	"net/http"
)

type server struct {
	srv      *http.Server
	listener net.Listener
}

func NewServer(options ...Option) *server {
	srv := &server{}

	for _, opt := range options {
		opt(srv)
	}
	return srv
}

func (s *server) Run() error {
	if err := s.srv.Serve(s.listener); err != nil {
		return err
	}
	return nil
}
