package app

import (
	"context"
	"fmt"
	"gateway/pkg/config"
	"gateway/pkg/logging"
	"gateway/pkg/proxy"
	"gateway/pkg/server"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	log := logging.NewLogger()
	log.Info("Starting Gateway...")

	if err := initConfig(); err != nil {
		log.Errorf("error initializing configs: %s", err.Error())
	}

	router := httprouter.New()

	dataReq, err := config.NewConfig(config.WithPath("config_proxy/data.json")).ParseConfig()

	if err != nil {
		log.Fatalf("error parsing data.json: %v", err)
	}

	for _, req := range dataReq.Data {
		router.Handle(req.Method, req.Path, proxy.NewProxy(
			proxy.WithProxy(req.MakeProxy),
			proxy.WithProxyUrl(req.ProxyUrl),
			proxy.WithRedirectUrl(req.Url),
			proxy.WithLog(log),
			proxy.WithExpectedStatusCodes(req.ExpectedProxyStatusCodes),
			proxy.WithProxyMethod(req.ProxyMethod),
		).Redirect())
	}

	log.Info("Starting server on port 8080...")

	readTimeOut, err := time.ParseDuration(viper.GetString("server.read_timeout"))

	if err != nil {
		fmt.Println(err.Error())
	}

	writeTimeOut, err := time.ParseDuration(viper.GetString("server.write_timeout"))

	if err != nil {
		fmt.Println(err.Error())
	}

	srv := server.NewServer(
		server.WithSrv(&http.Server{
			Addr:           fmt.Sprintf(":%s", viper.GetString("server.port")),
			Handler:        router,
			MaxHeaderBytes: viper.GetInt("server.max_header_bytes"),
			ReadTimeout:    readTimeOut,
			WriteTimeout:   writeTimeOut,
		}),
		server.WithHandler(router))

	go func() {
		if err = srv.Run(); err != nil {
			log.Errorf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Info("Gateway successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Gateway server shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
