package resthttp

import (
	"errors"
	"io"
	"net/http"
	"strconv"
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

	reqTodosUpdate struct {
		Title    string    `json:"title"`
		Body     string    `json:"body"`
		Deadline time.Time `json:"deadline"`
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
		if errors.Is(err, io.EOF) {
			respond(ctx, http.StatusBadRequest, nil, []string{ErrRequestBodyNotProvided.Error()})
			return
		}
		respond(ctx, http.StatusBadRequest, nil, []string{err.Error()})
		return
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

	respond(ctx, http.StatusCreated, respTodosCreate{
		ID: id,
	}, nil)
}

func (s *server) TodosUpdate(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}
	req := reqTodosUpdate{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			respond(ctx, http.StatusBadRequest, nil, []string{ErrRequestBodyNotProvided.Error()})
			return
		}
		respond(ctx, http.StatusBadRequest, nil, []string{err.Error()})
		return
	}

	err := s.todosService.Update(ctx, todos.UpdateInput{
		ID:       id,
		Title:    req.Title,
		Body:     req.Body,
		Deadline: req.Deadline,
	})
	if err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (s *server) TodosGet(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}

	todo, err := s.todosService.Get(ctx, id)
	if err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	respond(ctx, http.StatusOK, todo, []string{err.Error()})
}

func (s *server) TodosGetAll(ctx *gin.Context) {
	pageSize := ctx.Query("pageSize")
	page := ctx.Query("page")
	onlyCompleted := ctx.Query("onlyCompleted")
	sortBy := ctx.Query("sortBy")

	fPageSize := 10
	if len(pageSize) != 0 {
		n, err := strconv.Atoi(pageSize)
		if err != nil {
			respond(ctx, http.StatusBadRequest, nil, []string{err.Error()})
			return
		}
		if n <= 0 || n > 50 {
			respond(ctx, http.StatusBadRequest, nil, []string{"pageSize cannot be more than 50 or less than 1"})
		}
		fPageSize = n
	}

	fPage := 0
	if len(page) != 0 {
		n, err := strconv.Atoi(page)
		if err != nil {
			respond(ctx, http.StatusBadRequest, nil, []string{err.Error()})
			return
		}
		if n >= 0 {
			fPage = n
		}
	}

	fSortBy := todos.SortByCreationASC

	switch sortBy {
	case "creationDESC":
		fSortBy = todos.SortByCreationDESC
	case "deadlineASC":
		fSortBy = todos.SortByDeadlineASC
	case "deadlineDESC":
		fSortBy = todos.SortByCreationDESC
	}

	fOnlyCompleted := false
	if onlyCompleted == "true" {
		fOnlyCompleted = true
	}

	t, err := s.todosService.GetAll(ctx, todos.GetAllInput{
		PageSize:          fPageSize,
		Page:              fPage,
		ShowOnlyCompleted: fOnlyCompleted,
		SortBy:            fSortBy,
	})
	if err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	respond(ctx, http.StatusOK, t, nil)
}

func (s *server) TodosMarkComplete(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}

	if err := s.todosService.MarkAsComplete(ctx, id); err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (s *server) TodosMarkNotComplete(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}

	if err := s.todosService.MarkAsNotComplete(ctx, id); err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (s *server) TodosDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}

	if err := s.todosService.Delete(ctx, id); err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
