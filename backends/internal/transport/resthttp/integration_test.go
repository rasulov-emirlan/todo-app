package resthttp_test

// IMPORTANT these tests DO NOT WORK
// I AM TOO LAZY TO WRITE THEM

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rasulov-emirlan/todo-app/backends/config"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
	"github.com/rasulov-emirlan/todo-app/backends/wire"
)

const (
	postgresUsername = "postgres"
	postgresPassword = "postgres"
	postgresDbname   = "todo_test"
	apiUrl           = "http://localhost:8080"
)

var (
	postgresHostPort = "database:5432"

	logger  *logging.Logger
	store   *postgres.Repository
	cleanup func()
)

func TestMain(m *testing.M) {
	if testMain(m) != 0 {
		os.Exit(1)
	}
	err := error(nil)
	config := &config.Config{}
	logger, err = wire.InitializeLogger(*config)
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
		if err = server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("could not start server due to error: ", err)
			cleanup()
			os.Exit(1)
		}
	}()
	code := m.Run()
	if err = server.Shutdown(context.Background()); err != nil {
		log.Println("could not stop server due to error:", err)
		code = 1
	}
	cleanup()
	os.Exit(code)
}

func testMain(m *testing.M) int {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Printf("could not create docker pool: %v\n", err)
		return 1
	}
	cleanup = setupDB(pool)
	return 0
}

func setupDB(pool *dockertest.Pool) func() {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_PASSWORD=" + postgresPassword,
			"POSTGRES_USER=" + postgresUsername,
			"POSTGRES_DB=" + postgresDbname,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatal("Could not start postgres resource due to error: ", err)
	}

	postgresHostPort = resource.GetHostPort("5432/tcp")
	resource.Expire(200) // Tell docker to hard kill the container in 200 seconds
	pool.MaxWait = 200 * time.Second

	if err = pool.Retry(func() error {
		err := error(nil)
		store, err = postgres.NewRepository(config.Config{}, logger)
		if err != nil {
			return err
		}
		return store.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return func() {
		if err := store.Close(); err != nil {
			log.Fatal(err)
		}
		if err := pool.Purge(resource); err != nil {
			log.Fatal(err)
		}
	}
}

func TestServer(t *testing.T) {
	type Response struct {
		Errors []string `json:"errors"`
		Data   struct {
			Message string `json:"message"`
		} `json:"data"`
	}
	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Error(err)
	}
	r := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		resp.Body.Close()
		t.Errorf("response is incorrect")
	}
	if err := resp.Body.Close(); err != nil {
		t.Error(err)
	}
}
