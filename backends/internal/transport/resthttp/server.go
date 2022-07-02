package resthttp

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/log"
)

type server struct {
	server  *http.Server
	address string
	logger  *log.Logger

	usersService users.Service
	todosService todos.Service
}

func NewServer(
	corsOrigins []string,
	address string,
	rTimeout time.Duration,
	wTimeout time.Duration,
	logger *log.Logger,
	usersService users.Service,
	todosService todos.Service,
) *server {
	return &server{
		server: &http.Server{
			Addr:         address,
			ReadTimeout:  rTimeout,
			WriteTimeout: wTimeout,
		},
		address:      address,
		logger:       logger,
		usersService: usersService,
		todosService: todosService,
	}
}

func (s *server) Run() error {
	router := gin.New()
	s.setRoutes(router)
	s.server.Handler = router
	return s.server.ListenAndServe()
}

func (s *server) setRoutes(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(func(ctx *gin.Context) {
		// TODO: add cors handling
	})

	router.GET("/ping", func(c *gin.Context) {
		respond(c, http.StatusOK, gin.H{"message": "pong"}, nil)
	})

	usersGroup := router.Group("/users")
	{
		// TODO: separate users logic and auth
		usersGroup.POST("/auth/signup", s.UsersSignUp)
		usersGroup.POST("/auth/signin", s.UsersSignIn)
		usersGroup.POST("/auth/refresh", s.UsersRefresh)
		usersGroup.DELETE("/auth/logout", s.UsersLogout)

		usersGroup.DELETE("/:id", s.UsersDelete, s.isAdmin, s.requireAuth)
	}

	todosGroup := router.Group("/todos", s.requireAuth)
	{
		todosGroup.POST("", s.TodosCreate)
		todosGroup.GET("/:id", s.TodosGet)
		todosGroup.GET("", s.TodosGetAll)
		todosGroup.PATCH("/:id", s.TodosUpdate)
		todosGroup.DELETE("/:id", s.TodosDelete)
	}
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) Handler() http.Handler {
	return s.server.Handler
}
