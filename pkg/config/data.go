package config

type Requests struct {
	Path      string `json:"path"`
	Url       string `json:"url"`
	Method    string `json:"method"`
	MakeProxy bool   `json:"make_proxy"`
	ProxyUrl  string `json:"proxy"`
}

type Data struct {
	Data []Requests `json:"requests"`
}
