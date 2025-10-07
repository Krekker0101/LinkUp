package err

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Op    string
	Cause error
	Msg   string
	Code  int
}
// Auto-generated swagger comments for Error
// @Summary Auto-generated summary for Error
// @Description Auto-generated description for Error — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (e *APIError) Error() string {
	return fmt.Sprintf("api error in %s: %s | cause: %v", e.Op, e.Msg, e.Cause)
}

func NewAPIError(op string, cause error, msg string, code int) *APIError {
	return &APIError{Op: op, Cause: cause, Msg: msg, Code: code}
}
// Auto-generated swagger comments for LogAndRespondAPI
// @Summary Auto-generated summary for LogAndRespondAPI
// @Description Auto-generated description for LogAndRespondAPI — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func LogAndRespondAPI(c *gin.Context, err *APIError, userMsg string) {
	log.Printf("[API][%d] %v", err.Code, err)
	c.AbortWithStatusJSON(err.Code, gin.H{"error": userMsg})
}
