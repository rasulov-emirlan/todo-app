package resthttp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

type (
	reqTodosCreate struct {
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		Deadline time.Time `json:"deadline"`
	}
	respTodosCreate struct {
		ID string `json:"id"`
	}
)

func (s *server) TodosCreate(ctx *gin.Context) {
	user, ok := ctx.Get(usersInfoInContext)
	if !ok {
		panic("not implemented middleware")
	}

	userinfo, ok := user.(*users.JWTaccess)
	if !ok {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrNoCredentials.Error()})
		return
	}

	req := reqTodosCreate{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		respond(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	id, err := s.todosService.Create(ctx, todos.CreateInput{
		UserID:   userinfo.ID,
		Title:    req.Title,
		Body:     req.Body,
		Deadline: req.Deadline,
	})
	if err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	respond(ctx, http.StatusCreated, id, nil)
}
