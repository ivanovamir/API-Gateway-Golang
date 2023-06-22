package proxy

import (
	"bytes"
	"fmt"
	"gateway/pkg/config"
	handler_error "gateway/pkg/error"
	"gateway/pkg/logging"
	"io"
	"net/http"
	"strconv"
	"strings"
)

/*  */

type proxy struct {
	log         *logging.Logger
	ServiceName string
	requestData config.Request
	proxyHeader string
}

type Proxy interface {
	request(r *http.Request, url string, method string, body io.Reader) (*http.Response, error)
	prepareResponse(writer *http.ResponseWriter, resp *http.Response) error
	Redirect() http.HandlerFunc
	validateStatusCodes(receivedStatusCode int) error
}

func NewProxy(options ...Option) Proxy {
	p := &proxy{}
	for _, opt := range options {
		opt(p)
	}
	return p
}

func (p *proxy) Redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var body bytes.Buffer
		// Copy request body in buffer before original request is closed
		if _, err := io.Copy(&body, r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(handler_error.ErrorHandler(err))
			return
		}
		defer r.Body.Close()

		if p.requestData.MakeProxy {

			/*
				We give nil as request body, because we don't need it in proxy request.
				In these cases we use auth service as proxy.
				Proxy service need only headers.
			*/

			resp, err := p.request(r, p.requestData.ProxyUrl, p.requestData.ProxyMethod, nil)
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

			r.Header.Set("X-User-ID", resp.Header.Get("X-User-Id"))
		}

		resp, err := p.request(r, p.createUrl(strings.TrimPrefix(r.URL.String(), fmt.Sprintf("/%s", p.ServiceName))), r.Method, &body)

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

func (p *proxy) request(r *http.Request, url string, method string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)

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
		return fmt.Errorf("error copying response body: %s", err.Error())
	}
	return nil
}

func (p *proxy) validateStatusCodes(receivedStatusCode int) error {
	for _, val := range p.requestData.ExpectedProxyStatusCodes {
		sc, err := strconv.Atoi(val.StatusCode)
		if err != nil {
			return err
		}
		if sc == receivedStatusCode {
			return nil
		}
	}
	return fmt.Errorf("expected status codes: %v, but got: %d", p.requestData.ExpectedProxyStatusCodes, receivedStatusCode)
}

func (p *proxy) createUrl(oldUrl string) string {

	return p.requestData.Url + oldUrl
}
