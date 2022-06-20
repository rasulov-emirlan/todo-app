package resthttp

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"go.uber.org/zap"
)

type server struct {
	router  *gin.Engine
	address string
	logger  *zap.Logger

	usersService users.Service
	todosService todos.Service
}

func NewServer(
	corsOrigins []string,
	address string,
	rTimeout time.Duration,
	wTimeout time.Duration,
	logger *zap.Logger,
	usersService users.Service,
	todosService todos.Service,
) *server {
	r := gin.New()
	return &server{
		router:       r,
		address:      address,
		logger:       logger,
		usersService: usersService,
		todosService: todosService,
	}
}

func (s *server) Run() error {
	s.setRoutes()
	return s.router.Run(s.address)
}

func (s *server) setRoutes() {
	s.router.Use(gin.Recovery())
	s.router.Use(func(ctx *gin.Context) {
		// TODO: add cors handling
	})

	s.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, stdResponse{
			Status: 200,
			Data: gin.H{
				"message": "pong",
			},
			Errors: nil,
		})
	})

	usersGroup := s.router.Group("/users")
	{
		// TODO: separate users logic and auth
		usersGroup.POST("/auth/signup", s.UsersSignUp)
		usersGroup.POST("/auth/signin", s.UsersSignIn)
		usersGroup.GET("/auth/refresh", s.UsersRefresh)
	}
}
