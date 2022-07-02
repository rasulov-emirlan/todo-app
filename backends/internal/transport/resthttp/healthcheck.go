package resthttp

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type (
	// This is all info we could get about this servers process
	// swagger:model
	healthCheckResponse struct {
		MemoryUsage      runtime.MemStats
		ActiveGoRoutines int
	}
)

// swagger:route GET /health healthcheck healthCheck
//
// Checkup server
//
// This will return info about memory usage and goroutines
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
//       default: healthCheckResponse
func healthCheck(c *gin.Context) {
	resp := healthCheckResponse{}
	runtime.ReadMemStats(&resp.MemoryUsage)
	resp.ActiveGoRoutines = runtime.NumGoroutine()
	respond(c, http.StatusTeapot, resp, nil)
}
