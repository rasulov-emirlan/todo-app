package resthttp

import "github.com/gin-gonic/gin"

type (
	// stdResponse
	// represents the wrapper that all responses come inside of.
	//
	// swagger:model stdResponse
	stdResponse struct {
		// Error messages. It will not be omited if null
		// required: true
		// in: body
		Errors []string `json:"errors,omitempty"`

		// Actual Data that you expect to get on success. It will be omited if any errors occure
		// required: true
		// in: body
		// type: object
		Data interface{} `json:"data,omitempty"`
	}
)

func respond(ctx *gin.Context, status int, data interface{}, errors []string) {
	ctx.JSON(status, stdResponse{
		Data:   data,
		Errors: errors,
	})
}
