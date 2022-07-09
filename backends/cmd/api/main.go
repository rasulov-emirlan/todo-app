package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rasulov-emirlan/todo-app/backends/config"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres"
	"github.com/rasulov-emirlan/todo-app/backends/internal/transport/resthttp"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
	"go.uber.org/zap"
)

var (
	flagConfigName     = flag.String("config", "", "This flag accepts a path to .env file. If not provided we will get our configs from enviorment variables or we will use default values.")
	flagWithMigrations = flag.Bool("migrations", false, "If 'true' is given then migrations will be ran automaticaly on start of the app")
	flagIsDevMode      = flag.Bool("isDev", false, "If 'true' all of our services will start in development mode. Our keys will live longer. And our logs will be more informative")
)

func main() {
	flag.Parse()
	config, err := config.LoadConfigs(*flagConfigName)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.NewLogger(
		config.Log.Level,
		config.Log.Output,
	)
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatal(err)
		}
		if err := logger.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	logger.Info("Logger initialized")

	url := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.Database.User, config.Database.Pass),
		Host:   config.Database.Host + ":" + config.Database.Port,
		Path:   config.Database.Name,
	}
	store, err := postgres.NewRepository(url.String()+"?sslmode=disable", *flagWithMigrations)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Store initialized")

	usersService, err := users.NewService(
		store.Users(),
		logger,
		[]byte(config.JWTsecret))
	if err != nil {
		logger.Fatal("Could not initialize usersService", logging.String("error", err.Error()))
	}
	todosService := todos.NewService(
		store.Todos(),
		store.Users(),
		logger,
	)

	logger.Info("Services initialized")

	srvr := resthttp.NewServer([]string{"*"},
		config.Port, time.Second*15, time.Second*15,
		logger, usersService, todosService)

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
