package resthttp

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/todos"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

type (
	// reqTodosCreate
	// This is a model used for creating todos and only for that
	// swagger:model
	reqTodosCreate struct {
		// required: true
		// example: Do dishes tomorrow
		// min length: 6
		// max length: 100
		Title string `json:"title"`

		// max length: 2000
		Body string `json:"body"`

		// example: 2022-06-23T22:16:50.782647Z
		Deadline time.Time `json:"deadline"`
	}

	// respTodosCreate
	// This is an id of a newly created todo.
	// swagger:model
	respTodosCreate struct {
		ID string `json:"id"`
	}

	// This is info needed for updating a todo. Its 100% identical to reqTodosCreate
	// swagger:model
	reqTodosUpdate struct {
		// required: true
		// example: Do dishes tomorrow
		// min length: 6
		// max length: 100
		Title string `json:"title"`

		// max length: 2000
		Body string `json:"body"`

		// example: 2022-06-23T22:16:50.782647Z
		Deadline time.Time `json:"deadline"`
	}

	// todo
	// This is the actual model of a todo
	// swagger:model todo
	_ struct {
		// type: string
		// format: uuid
		ID string `json:"id"`

		Author *struct {
			// type: string
			// format: uuid
			ID           string `json:"id"`
			Username     string `json:"username"`
			Email        string `json:"email"`
			PasswordHash string `json:"-"`

			// Unknown: 0
			// Admin: 1
			// User: 2
			// type: int
			Role users.Role `json:"role"`

			CreatedAt time.Time `json:"createdAt"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"author,omitempty"`

		Title string `json:"title"`
		Body  string `json:"body"`

		Completed bool      `json:"completed"`
		Deadline  time.Time `json:"deadline"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}
)

// swagger:route POST /todos todo TodosCreate
//
// Create a todo
//
// This will create a todo. It will use Bearer token to identify
// caller of this endpoint and will use his identity as author of that todo
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//      - Bearer: []
//
//     Parameters:
//       + name: todo info
//         in: body
//         description: Basic info for a todo
//         required: true
//         type: reqTodosCreate
//
//     Responses:
//       default: respTodosCreate
//       200: respTodosCreate
//       400: stdResponse
//       422: stdResponse
func (s *Server) TodosCreate(ctx *gin.Context) {
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
		if errs := s.validator.UnpackErrors(err); errs != nil {
			respond(ctx, http.StatusBadRequest, nil, errs)
			return
		}
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	respond(ctx, http.StatusCreated, respTodosCreate{
		ID: id,
	}, nil)
}

// swagger:route PATCH /todos/{id} todo TodosUpdate
//
// Update a todo
//
// This will update a todo.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//      - Bearer: []
//
//     Parameters:
//       + name: todo info
//         in: body
//         description: Basic info for a todo
//         required: true
//         type: reqTodosUpdate
//       + name: id
//         in: params
//         required: true
//         description: Id of the todo you wish to update
//         type: string
//         example: '89cd8496-07cd-4caf-a9a5-ac3b8e65d05b'
//
//     Responses:
//       200: stdResponse
//       400: stdResponse
//       422: stdResponse
func (s *Server) TodosUpdate(ctx *gin.Context) {
	u, err := getUserData(ctx)
	if err != nil {
		respond(ctx, http.StatusUnauthorized, nil, []string{err.Error()})
		return
	}

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

	err = s.todosService.Update(ctx, u.ID, todos.UpdateInput{
		ID:       id,
		Title:    req.Title,
		Body:     req.Body,
		Deadline: req.Deadline,
	})
	if err != nil {
		if errs := s.validator.UnpackErrors(err); errs != nil {
			respond(ctx, http.StatusBadRequest, nil, errs)
			return
		}
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// swagger:route GET /todos{id} todo TodosGet
//
// Get a todo
//
// This will return a todo
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//      - Bearer: []
//
//     Parameters:
//       + name: id
//         in: params
//         required: true
//         description: Id of the todo you wish to update
//         type: string
//         example: '89cd8496-07cd-4caf-a9a5-ac3b8e65d05b'
//
//     Responses:
//       200: todo
//       400: stdResponse
//       422: stdResponse
func (s *Server) TodosGet(ctx *gin.Context) {
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

// swagger:route GET /todos todo TodosGetAll
//
// Get all todos
//
// This will return a list of your todos
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//      - Bearer: []
//
//     Parameters:
//       + name: pageSize
//         in: query
//         required: false
//         description: Number of todos to get
//         type: integer
//         example: 10
//       + name: page
//         in: query
//         required: false
//         type: integer
//         example: 0
//       + name: onlyCompleted
//         in: query
//         required: false
//         description: If true we will return only completed ones
//         type: boolean
//         example: true
//       + name: sortBy
//         in: query
//         required: false
//         description: How to sort it. Variations: [deadlineDESC, deadlineASC, creationDESC, creationASC]
//         type: string
//         example: deadlineDESC
//
//     Responses:
//       200: []todo
//       400: stdResponse
//       422: stdResponse
func (s *Server) TodosGetAll(ctx *gin.Context) {
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
	log.Println(fSortBy)

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

// swagger:route PUT /todos/{id}/complete todo TodosMakrAsComplete
//
// Mark as complete
//
// This will mark a todo as complete
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//      - Bearer: []
//
//     Parameters:
//       + name: id
//         in: params
//         required: true
//         description: Id for the todo
//         type: string
//
//     Responses:
//       400: stdResponse
//       422: stdResponse
func (s *Server) TodosMarkComplete(ctx *gin.Context) {
	u, err := getUserData(ctx)
	if err != nil {
		respond(ctx, http.StatusUnauthorized, nil, []string{err.Error()})
		return
	}

	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}

	if err := s.todosService.MarkAsComplete(ctx, u.ID, id); err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// swagger:route PUT /todos/{id}/incomplete todo TodosMakrAsNotComplete
//
// Mark as incomplete
//
// This will mark a todo as incomplete
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//      - Bearer: []
//
//     Parameters:
//       + name: id
//         in: params
//         required: true
//         description: Id for the todo
//         type: string
//
//     Responses:
//       400: stdResponse
//       422: stdResponse
func (s *Server) TodosMarkNotComplete(ctx *gin.Context) {
	u, err := getUserData(ctx)
	if err != nil {
		respond(ctx, http.StatusUnauthorized, nil, []string{err.Error()})
		return
	}
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}

	if err := s.todosService.MarkAsNotComplete(ctx, u.ID, id); err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

// swagger:route DELETE /todos/{id} todo TodosDelete
//
// Delete a todo
//
// This will delete a todo forever
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Deprecated: false
//
//     Security:
//      - Bearer: []
//
//     Parameters:
//       + name: id
//         in: params
//         required: true
//         description: Id for the todo
//         type: string
//
//     Responses:
//       400: stdResponse
//       422: stdResponse
func (s *Server) TodosDelete(ctx *gin.Context) {
	u, err := getUserData(ctx)
	if err != nil {
		respond(ctx, http.StatusUnauthorized, nil, []string{err.Error()})
		return
	}
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(ctx, http.StatusBadRequest, nil, []string{ErrParamNotProvided.Error()})
		return
	}

	if err := s.todosService.Delete(ctx, u.ID, id); err != nil {
		respond(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}
