package resthttp

import "github.com/gin-gonic/gin"

type (
	stdResponse struct {
		Status int         `json:"status"`
		Errors []string    `json:"errors"`
		Data   interface{} `json:"data"`
	}
)

func respond(ctx *gin.Context, status int, data interface{}, errors []string) {
	ctx.JSON(status, stdResponse{
		Status: status,
		Data:   data,
		Errors: errors,
	})
}
