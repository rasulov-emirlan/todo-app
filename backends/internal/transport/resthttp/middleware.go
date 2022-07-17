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

func (s *Server) isAdmin(ctx *gin.Context) {
	c, ok := ctx.Get(usersInfoInContext)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims, ok := c.(*users.JWTaccess)
	if !ok {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	if claims.Role != users.RoleAdmin {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Next()
}

func (s *Server) requireAuth(ctx *gin.Context) {
	accessKey := ctx.Request.Header.Get("Authorization")
	if accessKey == "" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	tokens := strings.Fields(accessKey)
	if len(tokens) != 2 {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	if tokens[0] != "Bearer" {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	claims, err := s.usersService.UnpackAccessKey(ctx, tokens[1])
	if err != nil {
		ctx.AbortWithStatus(http.StatusForbidden)
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

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		if method == "OPTIONS" {
			ctx.Header("Access-Control-Max-Age", "1728000")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type,Cookie,Authorization,Access-Control-Request-Headers,Access-Control-Request-Method,Origin,Referer,Sec-Fetch-Dest,Accept-Language,Accept-Encoding,Sec-Fetch-Mode,Sec-Fetch-Site,User-Agent,Pragma,Host,Connection,Cache-Control,Accept-Language,Accept-Encoding,X-Requested-With,X-Forwarded-For,X-Forwarded-Host,X-Forwarded-Proto,X-Forwarded-Port,X-Forwarded-Prefix,X-Real-IP,Accept")
			ctx.Header("Access-Control-Allow-Origin", "* http://localhost:3000")
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}
