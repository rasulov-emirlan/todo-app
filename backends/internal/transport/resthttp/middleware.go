package resthttp

import (
	"errors"
	"net/http"

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
	accessKey := ctx.Request.Header.Get("Authorization")
	if accessKey == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, err := s.usersService.UnpackAccessKey(ctx, accessKey)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims.Role != users.RoleAdmin {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Next()
}
