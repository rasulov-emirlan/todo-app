package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rasulov-emirlan/todo-app/backends/config"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
	"github.com/rasulov-emirlan/todo-app/backends/wire"
)

var (
	flagConfigName     = flag.String("config", "", "This flag accepts a path to .env file. If not provided we will get our configs from enviorment variables or we will use default values.")
	flagWithMigrations = flag.Bool("migrations", false, "If 'true' is given then migrations will be ran automaticaly on start of the app")
	// TODO: this flag should be used to start our server in debug mode
	// and enable panics in our services. Also make jwts live longer???
	flagIsDevMode = flag.Bool("isDev", false, "If 'true' all of our services will start in development mode. Our keys will live longer. And our logs will be more informative")
)

func main() {
	flag.Parse()
	config, err := config.LoadConfigs(*flagConfigName)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := wire.InitializeLogger(*config)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := wire.InitializeRepo(*config, logger)
	if err != nil {
		log.Fatal(err)
	}

	validator, err := wire.InitializeValidator()
	if err != nil {
		log.Fatal(err)
	}

	server, err := wire.InitializeRestApi(*config, logger, validator, repo)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Run()
		if err != nil {
			logger.Fatal("Server stopped", logging.String("error", err.Error()))
		}
	}()
	logger.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Gracefully stopping server")
	if err := repo.Close(); err != nil {
		logger.Fatal("Error closing store", logging.String("error", err.Error()))
	}
}
