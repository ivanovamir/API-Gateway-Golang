package proxy

import (
	"gateway/pkg/config"
	"gateway/pkg/logging"
)

type Option func(p *proxy)

func WithProxy(makeProxy bool) Option {
	return func(p *proxy) {
		p.proxy = makeProxy
	}
}

func WithProxyUrl(proxyUrl string) Option {
	return func(p *proxy) {
		p.proxyUrl = proxyUrl
	}
}

func WithRedirectUrl(redirectUrl string) Option {
	return func(p *proxy) {
		p.redirectUrl = redirectUrl
	}
}

func WithLog(log *logging.Logger) Option {
	return func(p *proxy) {
		p.log = log
	}
}

func WithExpectedStatusCodes(expectedStatusCodes []config.ExpectedStatusCodes) Option {
	return func(p *proxy) {
		p.expectedStatusCodes = expectedStatusCodes
	}
}

func WithProxyMethod(proxyMethod string) Option {
	return func(p *proxy) {
		p.proxyMethod = proxyMethod
	}
}
