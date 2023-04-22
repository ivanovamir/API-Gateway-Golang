package config

type Requests struct {
	Path                     string                `json:"path"`
	Url                      string                `json:"url"`
	Method                   string                `json:"method"`
	MakeProxy                bool                  `json:"make_proxy"`
	ProxyUrl                 string                `json:"proxy"`
	ProxyMethod              string                `json:"proxy_method"`
	ExpectedProxyStatusCodes []ExpectedStatusCodes `json:"expected_proxy_status_codes"`
}

type Data struct {
	Data []Requests `json:"requests"`
}

type ExpectedStatusCodes struct {
	StatusCode string `json:"status_code"`
}
