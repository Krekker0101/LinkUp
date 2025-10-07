package handlers

import (
	"net/http"
	"strconv"

	"LinkUp/internal/models"

	"github.com/gin-gonic/gin"
)

// @Summary История сообщений комнаты
// @Description Возвращает историю сообщений для указанной комнаты
// @Tags messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID комнаты"
// @Param limit query int false "Лимит сообщений" default(50)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} MessageResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /rooms/{id}/history [get]

func (h *Handler) MessageHistory(c *gin.Context) {
	roomID := c.Param("id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	var msgs []models.Message
	if err := h.db.Where("room_id = ?", roomID).Order("created_at desc").Limit(limit).Offset(offset).Find(&msgs).Error; err != nil {
		respondErr(c, 400, "load failed")
		return
	}

	for i, j := 0, len(msgs)-1; i < j; i, j = i+1, j-1 {
		msgs[i], msgs[j] = msgs[j], msgs[i]
	}

	var ids []uint
	for _, m := range msgs {
		ids = append(ids, m.ID)
	}
	var reacts []models.Reaction
	if len(ids) > 0 {
		h.db.Where("message_id IN ?", ids).Find(&reacts)
	}
	reactMap := map[uint]map[string][]uint{}
	for _, r := range reacts {
		if _, ok := reactMap[r.MessageID]; !ok {
			reactMap[r.MessageID] = map[string][]uint{}
		}
		reactMap[r.MessageID][r.Reaction] = append(reactMap[r.MessageID][r.Reaction], r.UserID)
	}
	res := []gin.H{}
	for _, m := range msgs {
		res = append(res, gin.H{"id": m.ID, "roomId": m.RoomID, "userId": m.UserID, "type": m.Type, "text": m.Text, "imageUrl": m.ImageURL, "createdAt": m.CreatedAt, "reactions": reactMap[m.ID]})
	}
	c.JSON(http.StatusOK, res)
}

// @Summary Отправить сообщение
// @Description Отправляет новое сообщение в указанную комнату
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID комнаты"
// @Param message body SendMessageRequest true "Текст сообщения"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /rooms/{id}/messages [post]
func (h *Handler) SendMessageREST(c *gin.Context) {
	roomID64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	roomID := uint(roomID64)

	var req struct {
		Type     string `json:"type"`
		Text     string `json:"text"`
		ImageURL string `json:"imageUrl"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondErr(c, 400, "invalid payload")
		return
	}

	msg := models.Message{
		RoomID:   roomID,
		UserID:   uid(c),
		Type:     req.Type,
		Text:     req.Text,
		ImageURL: req.ImageURL,
	}
	if msg.Type == "" {
		msg.Type = "text"
	}

	if err := h.db.Create(&msg).Error; err != nil {
		respondErr(c, 500, "db error")
		return
	}

	h.rooms.Emit(msg.RoomID, Event{
		Type: "message",
		Payload: gin.H{
			"id":        msg.ID,
			"roomId":    msg.RoomID,
			"userId":    msg.UserID,
			"type":      msg.Type,
			"text":      msg.Text,
			"imageUrl":  msg.ImageURL,
			"createdAt": msg.CreatedAt,
		},
	})

	c.JSON(200, gin.H{"ok": true})
}

type reactionReq struct {
	Reaction string `json:"reaction" binding:"required"`
}

// @Summary Добавить реакцию к сообщению
// @Description Добавляет реакцию к указанному сообщению
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID сообщения"
// @Param reaction body AddReactionRequest true "Реакция"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /messages/{id}/reactions [post]

func (h *Handler) AddReaction(c *gin.Context) {
	mid := c.Param("id")
	var msg models.Message
	if err := h.db.First(&msg, mid).Error; err != nil {
		respondErr(c, 404, "сооооообщение не найдено")
		return
	}
	var req reactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondErr(c, 400, "Инвалид пейлоад")
		return
	}
	r := models.Reaction{MessageID: msg.ID, UserID: uid(c), Reaction: req.Reaction}
	if err := h.db.Where("message_id = ? AND user_id = ? AND reaction = ?", msg.ID, uid(c), req.Reaction).FirstOrCreate(&r).Error; err != nil {
		respondErr(c, 400, "Не удалоооось добавить реакцию")
		return
	}

	h.rooms.Emit(msg.RoomID, Event{Type: "reaction", Payload: gin.H{"messageId": msg.ID, "userId": uid(c), "reaction": req.Reaction}})
	c.JSON(200, gin.H{"ok": true})
}

// @Summary Удалить реакцию
// @Description Удаляет реакцию с указанного сообщения
// @Tags messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID сообщения"
// @Param reaction path string true "Тип реакции"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /messages/{id}/reactions/{reaction} [delete]

func (h *Handler) RemoveReaction(c *gin.Context) {
	mid := c.Param("id")
	reac := c.Param("reaction")
	h.db.Where("message_id = ? AND user_id = ? AND reaction = ?", mid, uid(c), reac).Delete(&models.Reaction{})
	var msg models.Message
	if err := h.db.First(&msg, mid).Error; err == nil {
		h.rooms.Emit(msg.RoomID, Event{Type: "reaction_removed", Payload: gin.H{"messageId": msg.ID, "userId": uid(c), "reaction": reac}})
	}
	c.JSON(200, gin.H{"ok": true})
}
