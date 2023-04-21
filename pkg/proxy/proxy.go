package proxy

import (
	"fmt"
	handler_error "gateway/pkg/error"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
)

/*  */

type proxy struct {
	proxy       bool
	proxyUrl    string
	redirectUrl string
}

type Proxy interface {
	makeRequest(r *http.Request, url string) (*http.Response, error)
	prepareRedirect(writer *http.ResponseWriter, resp *http.Response) error
	proxyReq(r *http.Request) (*http.Response, error)
	Redirect() httprouter.Handle
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
			resp, err := p.proxyReq(r)
			// Если ошибка авторизации, возвращаем ошибку
			if err != nil {
				w.WriteHeader(resp.StatusCode)
				w.Write([]byte(err.Error()))
				return
			}
		}

		resp, err := p.makeRequest(r, p.redirectUrl)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(handler_error.ErrorHandler(err))
			return
		}

		if err = p.prepareRedirect(&w, resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(handler_error.ErrorHandler(err))
			return
		}
	}
}

func (p *proxy) makeRequest(r *http.Request, url string) (*http.Response, error) {
	req, err := http.NewRequest(r.Method, url, r.Body)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (p *proxy) prepareRedirect(writer *http.ResponseWriter, resp *http.Response) error {
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

func (p *proxy) proxyReq(r *http.Request) (*http.Response, error) {

	resp, err := p.makeRequest(r, p.proxyUrl)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%v", resp.Body)
	}

	return resp, nil
}
