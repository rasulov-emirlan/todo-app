// Package classification Todo App API
//
// the purpose of this application is to learn more about REST and swagger
//
// This should demonstrate how to write clean code in go
// and communicate with it using http
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     BasePath: /api
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Emirlan Rasulov<rasulov.emirlan@gmail.com> https://github.com/rasulov-emirlan
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//     - application/xml
//
//     Security:
//     - Bearer: []
//
//	   securityDefinitions:
//       Bearer:
//         type: apiKey
//         in: Header
//         name: Authorization
//
// swagger:meta
package resthttp

//go:generate swagger generate spec -o ./swaggerui/swagger.yaml --scan-models
//go:generate swagger generate spec -o ./swaggerui/swagger.json --scan-models

import (
	"context"
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/rasulov-emirlan/todo-app/backends/config"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/validation"
)

type Server struct {
	server  *http.Server
	address string

	// utility dependencies
	logger    *logging.Logger
	validator *validation.Validator

	// domain logic dependencies
	usersService users.Service
	todosService todos.Service
}

func NewServer(
	cfg config.Config,
	logger *logging.Logger,
	validator *validation.Validator,
	usersService users.Service,
	todosService todos.Service,
) *Server {
	return &Server{
		server: &http.Server{
			Addr:         cfg.Port,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		address:      cfg.Port,
		logger:       logger,
		validator:    validator,
		usersService: usersService,
		todosService: todosService,
	}
}

func (s *Server) Run() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	s.setRoutes(router)
	s.server.Handler = router
	return s.server.ListenAndServe()
}

//go:embed swaggerui
var swagger embed.FS

func (s *Server) setRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*", "Content-type", "Authorization"},
		AllowCredentials: true,
		AllowWildcard:    true,
	}))

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gzip.Gzip(gzip.BestCompression))
	api := router.Group("api")

	dir, err := fs.Sub(swagger, "swaggerui")
	if err != nil {
		s.logger.Fatal(err.Error())
	}
	api.StaticFS("/swagger", http.FS(dir))

	api.GET("/ping", func(c *gin.Context) {
		respond(c, http.StatusOK, gin.H{"message": "pong"}, nil)
	})

	api.GET("/health", s.requireAuth, s.isAdmin, healthCheck)

	usersGroup := api.Group("users")
	{
		// TODO: separate users logic and auth
		usersGroup.POST("/auth/signup", s.UsersSignUp)
		usersGroup.POST("/auth/signin", s.UsersSignIn)
		usersGroup.POST("/auth/refresh", s.UsersRefresh)
		usersGroup.DELETE("/auth/logout", s.UsersLogout)

		usersGroup.DELETE("/:id", s.requireAuth, s.isAdmin, s.UsersDelete)
	}

	todosGroup := api.Group("todos", s.requireAuth)
	{
		todosGroup.POST("", s.TodosCreate)
		todosGroup.GET("/:id", s.TodosGet)
		todosGroup.GET("", s.TodosGetAll)
		todosGroup.PATCH("/:id", s.TodosUpdate)
		todosGroup.PUT("/:id/complete", s.TodosMarkComplete)
		todosGroup.PUT("/:id/incomplete", s.TodosMarkNotComplete)

		todosGroup.DELETE("/:id", s.TodosDelete)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) Handler() http.Handler {
	return s.server.Handler
}
