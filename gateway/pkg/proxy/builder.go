package proxy

import (
	"gateway/pkg/config"
	"gateway/pkg/logging"
)

type Option func(p *proxy)

func WithRequest(req config.Request) Option {
	return func(p *proxy) {
		p.requestData = req
	}
}

func WithLog(log *logging.Logger) Option {
	return func(p *proxy) {
		p.log = log
	}
}

func WithServices(service string) Option {
	return func(p *proxy) {
		p.ServiceName = service
	}
}
