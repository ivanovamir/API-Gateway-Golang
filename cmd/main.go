package main

import (
	"fmt"
	"gateway/pkg/config"
	"gateway/pkg/proxy"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
	router := httprouter.New()

	dataReq, err := config.NewConfig(config.WithPath("data.json")).ParseConfig()

	if err != nil {
		log.Fatalf("error parsing data.json: %v", err)
	}

	for _, req := range dataReq.Data {
		router.Handle(req.Method, req.Path, proxy.NewProxy(
			proxy.WithProxy(req.MakeProxy),
			proxy.WithProxyUrl(req.ProxyUrl),
			proxy.WithRedirectUrl(req.Url),
		).Redirect())
	}

	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
