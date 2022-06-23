package resthttp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
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
	user, err := getUserData(ctx)
	if err != nil {
		respond(ctx, http.StatusUnauthorized, nil, []string{err.Error()})
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

	id, err := s.todosService.Create(context.Background(), todos.CreateInput{
		UserID:   user.ID,
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

	respond(ctx, http.StatusOK, todo, nil)
}

var sortVariants = map[string]todos.SortBy{
	"creationASC":  todos.SortByCreationASC,
	"creationDESC": todos.SortByCreationDESC,
	"deadlineASC":  todos.SortByDeadlineASC,
	"deadlineDESC": todos.SortByDeadlineDESC,
}

func (s *server) TodosGetAll(ctx *gin.Context) {
	user, err := getUserData(ctx)
	if err != nil {
		respond(ctx, http.StatusUnauthorized, nil, []string{err.Error()})
		return
	}
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
	if sorting, ok := sortVariants[sortBy]; ok {
		fSortBy = sorting
	}

	fOnlyCompleted := false
	if onlyCompleted == "true" {
		fOnlyCompleted = true
	}

	t, err := s.todosService.GetAll(ctx, todos.GetAllInput{
		UserID:            user.ID,
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
