package handlers

import (
	"github.com/gin-gonic/gin"
)
// Auto-generated swagger comments for respondErr
// @Summary Auto-generated summary for respondErr
// @Description Auto-generated description for respondErr — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func respondErr(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

func uid(c *gin.Context) uint {
	v, _ := c.Get("userID")
	if v == nil {
		return 0
	}
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}
