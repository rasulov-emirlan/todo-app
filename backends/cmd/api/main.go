package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rasulov-emirlan/todo-app/backends/config"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres"
	"github.com/rasulov-emirlan/todo-app/backends/internal/transport/resthttp"
	customLog "github.com/rasulov-emirlan/todo-app/backends/pkg/log"
	"go.uber.org/zap"
)

func main() {
	var cfgfilename string
	if len(os.Args) > 1 {
		cfgfilename = os.Args[1]
	}
	config, err := config.LoadConfigs(cfgfilename)
	if err != nil {
		log.Fatal(err)
	}

	logger := customLog.NewLogger(
		customLog.ParseLevel(config.Log.Level),
		config.Log.Output,
	)
	defer logger.Sync()

	log.Println("config:", config)
	logger.Info("Logger initialized")

	url := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.Database.User, config.Database.Pass),
		Host:   config.Database.Host + ":" + config.Database.Port,
		Path:   config.Database.Name,
	}
	store, err := postgres.NewRepository(url.String()+"?sslmode=disable", true)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Store initialized")

	usersService := users.NewService(
		store.Users(),
		logger,
		[]byte(config.JWTsecret))

	logger.Info("Services initialized")

	srvr := resthttp.NewServer([]string{"*"},
		config.Port, time.Second*15, time.Second*15,
		logger, usersService, nil)

	logger.Info("Server initialized")

	go func() {
		err := srvr.Run()
		if err != nil {
			logger.Fatal("Server stopped", zap.Error(err))
		}
	}()
	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Gracefully stopping server")
	if err := store.Close(); err != nil {
		logger.Fatal("Error closing store", zap.Error(err))
	}
}
