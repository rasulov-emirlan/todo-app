package resthttp

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

type (
	// reqUsersSignUp
	// represents all info needed for a user to sign up
	// It is used only for signing up
	//
	// swagger:model
	reqUsersSignUp struct {
		// the email address for this user
		//
		// required: true
		// example: user@example.com
		Email string `json:"email"`

		// the password for this user
		//
		// required: true
		// min length: 6
		// max length: 128
		// example: password
		Password string `json:"password"`

		// the name for this user
		//
		// required: true
		// min length: 6
		// max length: 20
		// example: John Doe
		Username string `json:"username"`
	}

	// reqUsersSignIn
	// represents all info needed for a user to sign in
	//
	// swagger:model
	reqUsersSignIn struct {
		// the email address for this user
		//
		// required: true
		// example: user@example.com
		Email string `json:"email"`

		// password for this user
		// required: true
		// example: password
		Password string `json:"password"`
	}

	// usersKeys represents a pair of keys used for authorization in our api
	//
	// swagger:response usersKeys
	respUsersKeys struct {
		// accessKey. Its life span is 10 minutes
		// in: body
		AccessKey string `json:"accessKey"`

		// If your client is browser then you should not use
		// this field. We will set up a cookie for you, so relax bro.
		// Its life span is 7 days.
		//
		// in: body
		RefreshKey string `json:"refreshKey"`

		// This cookie will contain refresh key
		//
		// in: cookie
		refresh_key string
	}

	// user represents a response for /users/me endpoint
	// swagger:model
	usersMeResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`

		Role users.Role `json:"role"`

		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	// reqUsersRefresh is used for mobile clients. They should send their refresh keys in this model to refresh endpoint for updating their keys
	//
	// swagger:model
	reqUsersRefresh struct {
		RefreshKey string `json:"refreshKey"`
	}
)

const (
	cookieNameRefreshKey = "refresh_key"

	accessLifeTime  = time.Minute * 10
	refreshLifeTime = time.Hour * 24 * 7
)

var (
	ErrNoRefreshProvided      = errors.New("refresh key has to be provided via cookie or request body")
	ErrRequestBodyNotProvided = errors.New("request body has to be provided")
	ErrParamNotProvided       = errors.New("parameter has to be provided")
)

// swagger:route POST /users/auth/signup auth UsersSignUp
//
// Sign up a user
//
// This will create a user in our database IF AND ONLY
// if he doesnt exist yet. After creating him it will automaticaly
// sign him in
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
//     Parameters:
//       + name: user info
//         in: body
//         description: Basic info for user to sign up
//         required: true
//         type: reqUsersSignUp
//
//     Responses:
//       default: usersKeys
//       200: usersKeys
//       422: stdResponse
func (s *Server) UsersSignUp(ctx *gin.Context) {
	var inp reqUsersSignUp
	if err := ctx.ShouldBindJSON(&inp); err != nil {
		if errors.Is(err, io.EOF) {
			err = ErrRequestBodyNotProvided
		}
		respond(
			ctx,
			http.StatusBadRequest,
			nil,
			[]string{err.Error()},
		)
		return
	}

	out, err := s.usersService.SignUp(
		ctx,
		users.SignUpInput{
			Email:    inp.Email,
			Username: inp.Username,
			Password: inp.Password,
		},
	)
	if err != nil {
		errs := s.validator.UnpackErrors(err)
		resp := []string{err.Error()}
		if errs != nil {
			resp = errs
		}
		respond(
			ctx,
			http.StatusInternalServerError,
			nil,
			resp,
		)
		return
	}

	ctx.SetCookie(
		cookieNameRefreshKey,
		out.RefreshKey,
		int(refreshLifeTime.Seconds()),
		"/",
		"",
		false,
		true,
	)

	respond(ctx, http.StatusOK, out, nil)
}

// swagger:route POST /users/auth/signin auth UsersSignIn
//
// Sign in a user
//
// This should return a pair of keys for the user, if user info provided is valid
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
//     Parameters:
//       + name: user info
//         in: body
//         description: Basic info for user to sign up
//         required: true
//         type: reqUsersSignIn
//
//     Responses:
//       default: usersKeys
//       200: usersKeys
//       422: stdResponse
func (s *Server) UsersSignIn(ctx *gin.Context) {
	var inp reqUsersSignIn

	if err := ctx.ShouldBindJSON(&inp); err != nil {
		respond(
			ctx,
			http.StatusBadRequest,
			nil,
			[]string{err.Error()},
		)
		return
	}

	out, err := s.usersService.SignIn(
		ctx,
		inp.Email,
		inp.Password,
	)
	if err != nil {
		respond(
			ctx,
			http.StatusInternalServerError,
			nil,
			[]string{err.Error()},
		)
		return
	}

	ctx.SetCookie(
		cookieNameRefreshKey,
		out.RefreshKey,
		int(refreshLifeTime.Seconds()),
		"/",
		"",
		false,
		true,
	)

	respond(ctx, http.StatusOK, out, nil)
}

// swagger:route POST /users/auth/refresh auth UsersRefresh
//
// Refresh keys
//
// This is supposed to return a new pair of keys. It will check the body for refresh key.
// If it wont find it in body, then it will check refresh_key Cookie. If both are empty then
// ur mad bro.
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
//     Parameters:
//       + name: refresh
//         in: body
//         description: Refresh key
//         required: false
//         type: reqUsersRefresh
//		 + name: refreshCookie
//         in: cookie
//         description: Refresh Key
//         type: string
//
//     Responses:
//       default: usersKeys
//       200: usersKeys
//       422: stdResponse
func (s *Server) UsersRefresh(ctx *gin.Context) {
	var inp reqUsersRefresh
	if err := ctx.ShouldBindJSON(&inp); err != nil {
		inp.RefreshKey, err = ctx.Cookie(cookieNameRefreshKey)
		if err != nil {
			respond(
				ctx,
				http.StatusBadRequest,
				nil,
				[]string{ErrNoRefreshProvided.Error()},
			)
			return
		}
	}

	out, err := s.usersService.Refresh(
		ctx,
		inp.RefreshKey,
	)
	if err != nil {
		respond(
			ctx,
			http.StatusInternalServerError,
			nil,
			[]string{err.Error()},
		)
		return
	}

	ctx.SetCookie(
		cookieNameRefreshKey,
		out.RefreshKey,
		int(refreshLifeTime.Seconds()),
		"/",
		"",
		false,
		true,
	)

	respond(ctx, http.StatusOK, out, nil)
}

// swagger:route DELETE /users/auth/logout auth UsersLogout
//
// Logout
//
// This will delete refresh_key cookie.
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
//     Responses:
//       default: stdResponse
func (s *Server) UsersLogout(ctx *gin.Context) {
	ctx.SetCookie(
		cookieNameRefreshKey,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	respond(ctx, http.StatusOK, nil, nil)
}

// swagger:route DELETE /users{id} users UsersDelete
//
// Delete a user
//
// This will delete delete user.
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
//         description: Id of the user you wish to delete
//         type: string
//         example: '89cd8496-07cd-4caf-a9a5-ac3b8e65d05b'
//
//     Responses:
//       default: stdResponse
func (s *Server) UsersDelete(ctx *gin.Context) {
	id := ctx.Param("id")
	if len(id) == 0 {
		respond(
			ctx,
			http.StatusBadRequest,
			nil,
			[]string{ErrParamNotProvided.Error()},
		)
		return
	}

	if err := s.usersService.Delete(ctx, id); err != nil {
		respond(
			ctx,
			http.StatusInternalServerError,
			nil,
			[]string{err.Error()},
		)
		return
	}

	respond(ctx, http.StatusOK, nil, nil)
}

// swagger:route GET /users/me users UsersMe
//
// Get current user
//
// This will return current user.
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
//     Responses:
//       default: usersMeResponse
func (s *Server) usersMe(ctx *gin.Context) {
	d, err := getUserData(ctx)
	if err != nil {
		respond(
			ctx,
			http.StatusUnauthorized,
			nil,
			[]string{err.Error()},
		)
		return
	}

	out, err := s.usersService.Me(ctx, d.ID)
	if err != nil {
		respond(
			ctx,
			http.StatusInternalServerError,
			nil,
			[]string{err.Error()},
		)
		return
	}

	respond(ctx, http.StatusOK, out, nil)
}
