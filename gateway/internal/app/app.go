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
		mltPlexer    = httprouter.New()
		dataReq, err = config.NewConfig(config.WithPath("config_proxy/config.json")).ParseConfig()
	)

	if err != nil {
		log.Fatalf("error parsing config.json: %v", err)
	}

	for _, service := range dataReq.Services {
		for _, request := range service.Requests {
			ttt := "/" + service.Service + request.Path
			mltPlexer.Handler(request.Method, ttt, proxy.NewProxy(
				proxy.WithLog(log),
				proxy.WithRequest(request),
				proxy.WithServices(service.Service),
			).Redirect())
		}
	}

	readTimeOut, err := time.ParseDuration(viper.GetString("gateway_service.read_timeout"))

	if err != nil {
		log.Fatalf("errror occured parsing read_timeout from config: %s", err.Error())
	}

	writeTimeOut, err := time.ParseDuration(viper.GetString("gateway_service.write_timeout"))

	if err != nil {
		log.Fatalf("errror occured parsing write_timeout from config: %s", err.Error())
	}

	ln, err := net.Listen(viper.GetString("gateway_service.con_type"), fmt.Sprintf(":%s", viper.GetString("gateway_service.port")))

	if err != nil {
		log.Fatalf("error occured listener: %s", err.Error())
	}

	var (
		srv = server.NewServer(
			server.WithListener(&ln),
			server.WithSrv(&http.Server{
				Handler:        mltPlexer,
				MaxHeaderBytes: viper.GetInt("gateway_service.max_header_bytes"),
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

	log.Info("Gateway service successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Gateway service shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
