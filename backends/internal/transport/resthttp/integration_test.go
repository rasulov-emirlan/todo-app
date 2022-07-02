package resthttp_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/internal/storage/postgres"
	"github.com/rasulov-emirlan/todo-app/backends/internal/transport/resthttp"
	customLogger "github.com/rasulov-emirlan/todo-app/backends/pkg/log"
)

const (
	postgresUsername = "postgres"
	postgresPassword = "postgres"
	postgresDbname   = "todo_test"
	apiUrl           = "http://localhost:8080"
)

var (
	postgresHostPort = "database:5432"

	store   *postgres.Repository
	cleanup func()
)

func TestMain(m *testing.M) {
	if testMain(m) != 0 {
		os.Exit(1)
	}
	logger, err := customLogger.NewLogger("debug", "stdout")
	if err != nil {
		log.Fatal(err)
	}
	uService := users.NewService(store.Users(), logger, []byte("secretkey"))
	tService := todos.NewService(store.Todos(), logger)
	srvr := resthttp.NewServer(
		[]string{"*"}, ":8080", time.Second*15, time.Second*15, logger, uService, tService)
	go func() {
		if err = srvr.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("could not start server due to error: ", err)
			cleanup()
			os.Exit(1)
		}
	}()
	code := m.Run()
	if err = srvr.Shutdown(context.Background()); err != nil {
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
	dburl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		postgresUsername, postgresPassword, postgresHostPort, postgresDbname)
	resource.Expire(200) // Tell docker to hard kill the container in 200 seconds
	pool.MaxWait = 200 * time.Second

	if err = pool.Retry(func() error {
		err := error(nil)
		store, err = postgres.NewRepository(dburl, true)
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
