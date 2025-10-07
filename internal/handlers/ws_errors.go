// [COPILOT-BEGIN]
package handlers

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
)

type WebSocketError struct {
	Op    string
	Cause error
	Msg   string
}
// Auto-generated swagger comments for Error
// @Summary Auto-generated summary for Error
// @Description Auto-generated description for Error — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (e *WebSocketError) Error() string {
	return fmt.Sprintf("ws error in %s: %s | cause: %v", e.Op, e.Msg, e.Cause)
}

func NewWebSocketError(op string, cause error, msg string) *WebSocketError {
	return &WebSocketError{Op: op, Cause: cause, Msg: msg}
}
// Auto-generated swagger comments for LogAndRespondWS
// @Summary Auto-generated summary for LogAndRespondWS
// @Description Auto-generated description for LogAndRespondWS — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func LogAndRespondWS(c *gin.Context, status int, err error, userMsg string) {
	log.Printf("[WS][%d] %v", status, err)
	c.AbortWithStatusJSON(status, gin.H{"error": userMsg})
}