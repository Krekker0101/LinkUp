package handlers

import (
	"net/http"
	"strings"

	"LinkUp/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary Поиск сообщений
// @Description Ищет сообщения по тексту в указанной комнате или во всех комнатах
// @Tags search
// @Security BearerAuth
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Param room_id query string false "ID комнаты"
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /search [get]

func (h *Handler) SearchMessages(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if q == "" {
		respondErr(c, 400, "q is required")
		return
	}
	roomID := c.Query("roomId")
	query := h.db.Model(&models.Message{}).Where("type = ?", "text").Where("text LIKE ?", "%"+q+"%")
	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}
	var msgs []models.Message
	query.Order("created_at desc").Limit(100).Find(&msgs)
	c.JSON(http.StatusOK, msgs)
}
