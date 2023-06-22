package config

type Services struct {
	Service  string    `json:"service"`
	Requests []Request `json:"requests"`
}

type Request struct {
	Path                     string                `json:"path"`
	Url                      string                `json:"url"`
	Method                   string                `json:"method"`
	MakeProxy                bool                  `json:"make_proxy"`
	ProxyUrl                 string                `json:"proxy_url"`
	ProxyMethod              string                `json:"proxy_method"`
	ExpectedProxyStatusCodes []ExpectedStatusCodes `json:"expected_proxy_status_codes"`
}

type Data struct {
	Services []Services `json:"services"`
}

type ExpectedStatusCodes struct {
	StatusCode string `json:"status_code"`
}
