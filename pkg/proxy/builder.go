package proxy

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
