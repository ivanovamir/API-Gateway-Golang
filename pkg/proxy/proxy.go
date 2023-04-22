package proxy

import (
	"fmt"
	"gateway/pkg/config"
	handler_error "gateway/pkg/error"
	"gateway/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
)

/*  */

type proxy struct {
	proxy               bool
	proxyUrl            string
	redirectUrl         string
	log                 *logging.Logger
	expectedStatusCodes []config.ExpectedStatusCodes
	proxyMethod         string
}

type Response struct {
	ProxyResponse   []byte `json:"proxy_response"`
	ServiceResponse []byte `json:"service_response"`
}

type Proxy interface {
	request(r *http.Request, url string, method string) (*http.Response, error)
	prepareResponse(writer *http.ResponseWriter, resp *http.Response) error
	Redirect() httprouter.Handle
	validateStatusCodes(receivedStatusCode int) error
}

func NewProxy(options ...Option) Proxy {
	p := &proxy{}
	for _, opt := range options {
		opt(p)
	}
	return p
}

func (p *proxy) Redirect() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		if p.proxy {
			resp, err := p.request(r, p.proxyUrl, p.proxyMethod)
			// If err from proxy, return error
			if err != nil {
				w.WriteHeader(resp.StatusCode)
				w.Write(handler_error.ErrorHandler(err))
				return
			}

			body, err := io.ReadAll(resp.Body)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(handler_error.ErrorHandler(err))
				return
			}

			if err = p.validateStatusCodes(resp.StatusCode); err != nil {
				w.WriteHeader(resp.StatusCode)
				w.Write([]byte(fmt.Sprintf("%s", string(body))))
				return
			}
		}

		resp, err := p.request(r, p.redirectUrl, r.Method)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(handler_error.ErrorHandler(err))
			return
		}

		if err = p.prepareResponse(&w, resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(handler_error.ErrorHandler(err))
			return
		}
	}
}

func (p *proxy) request(r *http.Request, url string, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, r.Body)

	if err != nil {
		return nil, err
	}

	for header, values := range r.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *proxy) prepareResponse(writer *http.ResponseWriter, resp *http.Response) error {
	defer resp.Body.Close()
	// Возвращаем все хедеры и их значения от сервиса
	for key, value := range resp.Header {
		for _, val := range value {
			(*writer).Header().Set(key, val)
		}
	}

	// Возвращаем статус код ответа от сервиса
	(*writer).WriteHeader(resp.StatusCode)

	// Копируем тело ответа сервиса в тело ответа клиенту
	if _, err := io.Copy(*writer, resp.Body); err != nil {
		return fmt.Errorf("error copying response body: ", err)
	}
	return nil
}

func (p *proxy) validateStatusCodes(receivedStatusCode int) error {
	for _, val := range p.expectedStatusCodes {
		sc, err := strconv.Atoi(val.StatusCode)
		if err != nil {
			return err
		}
		if sc == receivedStatusCode {
			return nil
		}
	}
	return fmt.Errorf("expected status codes: %v, but got: %d", p.expectedStatusCodes, receivedStatusCode)
}
