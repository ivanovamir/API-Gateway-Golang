package app

import (
	"context"
	"fmt"
	"gateway/pkg/config"
	"gateway/pkg/logging"
	"gateway/pkg/proxy"
	"gateway/pkg/server"
	"github.com/spf13/viper"
	"net"
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

	var (
		mltPlexer    = http.NewServeMux()
		dataReq, err = config.NewConfig(config.WithPath("config_proxy/data.json")).ParseConfig()
	)

	if err != nil {
		log.Fatalf("error parsing data.json: %v", err)
	}

	for _, req := range dataReq.Data {
		mltPlexer.Handle(req.Path, proxy.NewProxy(
			proxy.WithProxy(req.MakeProxy),
			proxy.WithProxyUrl(req.ProxyUrl),
			proxy.WithRedirectUrl(req.Url),
			proxy.WithLog(log),
			proxy.WithExpectedStatusCodes(req.ExpectedProxyStatusCodes),
			proxy.WithProxyMethod(req.ProxyMethod),
			proxy.WithRequestMethod(req.Method),
		).Redirect())
	}

	log.Info("Starting server on port 8080...")

	readTimeOut, err := time.ParseDuration(viper.GetString("server.read_timeout"))

	if err != nil {
		log.Fatalf("errror occured parsing read_timeout from config: %s", err.Error())
	}

	writeTimeOut, err := time.ParseDuration(viper.GetString("server.write_timeout"))

	if err != nil {
		log.Fatalf("errror occured parsing write_timeout from config: %s", err.Error())
	}

	ln, err := net.Listen(viper.GetString("server.server_protocol"), fmt.Sprintf(":%s", viper.GetString("server.port")))

	if err != nil {
		log.Fatalf("error occured listening on port: %s", viper.GetString("server.port"))
	}

	var (
		srv = server.NewServer(
			server.WithListener(&ln),
			server.WithSrv(&http.Server{
				Handler:        mltPlexer,
				MaxHeaderBytes: viper.GetInt("server.max_header_bytes"),
				ReadTimeout:    readTimeOut,
				WriteTimeout:   writeTimeOut,
			}),
			server.WithHandler(mltPlexer))
	)
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

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
