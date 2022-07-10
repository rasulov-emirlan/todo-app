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
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
	"github.com/rasulov-emirlan/todo-app/backends/pkg/logging"
)

type server struct {
	server  *http.Server
	address string
	logger  *logging.Logger

	usersService users.Service
	todosService todos.Service
}

func NewServer(
	corsOrigins []string,
	address string,
	rTimeout time.Duration,
	wTimeout time.Duration,
	logger *logging.Logger,
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	s.setRoutes(router)
	s.server.Handler = router
	return s.server.ListenAndServe()
}

//go:embed swaggerui
var swagger embed.FS

func (s *server) setRoutes(router *gin.Engine) {
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gzip.Gzip(gzip.BestCompression))
	api := router.Group("api")
	api.Use(CORSMiddleware())

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

func (s *server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *server) Handler() http.Handler {
	return s.server.Handler
}
