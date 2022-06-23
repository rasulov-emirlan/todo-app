package resthttp

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rasulov-emirlan/todo-app/backends/internal/domain/users"
)

const (
	usersInfoInContext = "userinfo"
)

var (
	ErrNoCredentials = errors.New("could not find Athorization Bearer token in headers")
)

func (s *server) isAdmin(ctx *gin.Context) {
	c, ok := ctx.Get(usersInfoInContext)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, ok := c.(*users.JWTaccess)
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if claims.Role != users.RoleAdmin {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Next()
}

func (s *server) requireAuth(ctx *gin.Context) {
	accessKey := ctx.Request.Header.Get("Authorization")
	if accessKey == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokens := strings.Fields(accessKey)
	if len(tokens) != 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if tokens[0] != "Bearer" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, err := s.usersService.UnpackAccessKey(ctx, tokens[1])
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.Set(usersInfoInContext, &claims)
	ctx.Next()
}

func getUserData(ctx *gin.Context) (users.JWTaccess, error) {
	info, ok := ctx.Get(usersInfoInContext)
	if !ok {
		return users.JWTaccess{}, ErrNoCredentials
	}
	claims, ok := info.(*users.JWTaccess)
	if !ok {
		return users.JWTaccess{}, errors.New("could not map users info into jwt")
	}
	return *claims, nil
}
