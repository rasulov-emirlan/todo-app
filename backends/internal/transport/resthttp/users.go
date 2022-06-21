package resthttp

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type (
	reqUsersSignUp struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	reqUsersSignIn struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
	ErrNoRefreshProvided = errors.New("refresh key has to be provided via cookie or request body")
)

func (s *server) UsersSignUp(ctx *gin.Context) {
	var inp reqUsersSignUp
	if err := ctx.ShouldBindJSON(&inp); err != nil {
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
		inp.Email,
		inp.Password,
		inp.Username,
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

func (s *server) UsersSignIn(ctx *gin.Context) {
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

func (s *server) UsersRefresh(ctx *gin.Context) {
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

func (s *server) UsersLogout(ctx *gin.Context) {
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
