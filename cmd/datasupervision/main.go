package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"datasupervision/internal/controller"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	configPath := flag.String("c", "config.yaml", "specify path to a config.yaml")
	flag.Parse()

	cfg, err := configure(*configPath)
	if err != nil {
		log.WithError(err).Fatal("can't read config")

		return
	}

	ctx := context.Background()

	server, err := controller.NewServer(&cfg.Server)
	if err != nil {
		log.WithError(err).Fatal("can't create server")
	}

	log.Info("server create successfully")

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer stop()

	go func() {
		log.Info("starting http server")

		serverErr := server.Server.ListenAndServe()
		if serverErr != nil {
			log.Error(serverErr)
		}
	}()
	<-ctx.Done()
	stop()
	log.Info("shutting down gracefully, press Ctrl+C again to force")

	//nolint:gomnd // timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = server.Server.Shutdown(timeoutCtx); err != nil {
		log.WithError(err).Error("error on shutdown")
	}

}
