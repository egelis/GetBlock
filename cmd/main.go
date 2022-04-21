package main

import (
	"context"
	"github.com/egelis/GetBlock/internal/getblock"
	"github.com/egelis/GetBlock/internal/service"
	"github.com/egelis/GetBlock/internal/transport"
	"github.com/egelis/GetBlock/internal/transport/handler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	LogLevel string
	Address  string
	APIKey   string
)

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	LogLevel = viper.GetString("logLevel")
	Address = viper.GetString("address")
	APIKey = viper.GetString("apiKey")

	return nil
}

func initLogger() error {
	level, err := log.ParseLevel(LogLevel)
	if err != nil {
		return err
	}

	log.SetLevel(level)

	log.SetFormatter(new(log.JSONFormatter))

	return nil
}

func main() {
	if err := initConfig(); err != nil {
		log.Fatal("error initializing configs: ", err)
	}

	if err := initLogger(); err != nil {
		log.Fatal("error initializing logger:", err)
	}

	gbClient := getblock.NewClient(APIKey)
	srv := service.NewService(gbClient)
	handlers := handler.NewHandler(srv).InitRoutes()

	server := transport.NewServer(Address, handlers)

	go func() {
		log.Info("server is running")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error running server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	log.Info("server is shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("error occurred on server shutting down:", err)
	}
}
