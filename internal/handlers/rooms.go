package handlers

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"LinkUp/internal/models"
)

type createRoomReq struct {
	Slug      string `json:"slug" binding:"required"`
	Name      string `json:"name" binding:"required"`
	IsPrivate bool   `json:"isPrivate"`
}

// @Summary Создать комнату
// @Description Создает новую комнату для чата
// @Tags rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param room body CreateRoomRequest true "Данные комнаты"
// @Success 201 {object} RoomResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /rooms [post]
func (h *Handler) CreateRoom(c *gin.Context) {
	var req createRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondErr(c, 400, "invalid body")
		return
	}
	r := models.Room{Slug: strings.ToLower(req.Slug), Name: req.Name, IsPrivate: req.IsPrivate, OwnerID: uid(c)}
	if err := h.db.Create(&r).Error; err != nil {
		respondErr(c, 409, "room slug exists?")
		return
	}

	h.db.Where(models.RoomMember{RoomID: r.ID, UserID: uid(c)}).FirstOrCreate(&models.RoomMember{})
	c.JSON(201, r)
}

// @Summary Получить список комнат
// @Description Возвращает список всех доступных комнат
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Success 200 {array} RoomResponse
// @Failure 401 {object} ErrorResponse
// @Router /rooms [get]
func (h *Handler) ListRooms(c *gin.Context) {
	var rooms []models.Room
	h.db.Order("created_at asc").Find(&rooms)

	res := []gin.H{}
	for _, r := range rooms {
		cnt := h.unreadCount(c, r.ID, uid(c))
		res = append(res, gin.H{"id": r.ID, "slug": r.Slug, "name": r.Name, "isPrivate": r.IsPrivate, "unread": cnt})
	}
	c.JSON(200, res)
}

// Auto-generated swagger comments for unreadCount
// @Summary Auto-generated summary for unreadCount
// @Description Auto-generated description for unreadCount — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func (h *Handler) unreadCount(c *gin.Context, roomID, userID uint) int64 {
	var m models.RoomMember
	if err := h.db.Where("room_id = ? AND user_id = ?", roomID, userID).First(&m).Error; err != nil {
		return 0
	}
	var count int64
	if m.LastReadAt == nil {
		h.db.Model(&models.Message{}).Where("room_id = ?", roomID).Count(&count)
	} else {
		h.db.Model(&models.Message{}).Where("room_id = ? AND created_at > ?", roomID, m.LastReadAt).Count(&count)
	}
	return count
}

// @Summary Присоединиться к комнате
// @Description Присоединяет текущего пользователя к указанной комнате
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID комнаты"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /rooms/{id}/join [post]
func (h *Handler) JoinRoom(c *gin.Context) {
	roomID := c.Param("id")
	var r models.Room
	if err := h.db.First(&r, roomID).Error; err != nil {
		respondErr(c, 404, "room not found")
		return
	}

	h.db.Where(models.RoomMember{RoomID: r.ID, UserID: uid(c)}).FirstOrCreate(&models.RoomMember{})
	c.JSON(200, gin.H{"ok": true})
}

// @Summary Покинуть комнату
// @Description Удаляет текущего пользователя из указанной комнаты
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID комнаты"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /rooms/{id}/leave [post]
func (h *Handler) LeaveRoom(c *gin.Context) {
	roomID := c.Param("id")
	h.db.Where("room_id = ? AND user_id = ?", roomID, uid(c)).Delete(&models.RoomMember{})
	c.JSON(200, gin.H{"ok": true})
}

// @Summary Список пользователей в комнате
// @Description Возвращает список всех пользователей в указанной комнате
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID комнаты"
// @Success 200 {array} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /rooms/{id}/users [get]
func (h *Handler) RoomMembers(c *gin.Context) {
	roomID := c.Param("id")
	var rms []models.RoomMember
	if err := h.db.Where("room_id = ?", roomID).Find(&rms).Error; err != nil {
		respondErr(c, 404, "room not found")
		return
	}
	userIDs := []uint{}
	for _, m := range rms {
		userIDs = append(userIDs, m.UserID)
	}
	var users []models.User
	h.db.Where("id IN ?", userIDs).Find(&users)
	res := []gin.H{}
	for _, u := range users {
		u.Online = h.presence.IsOnline(u.ID)
		res = append(res, gin.H{"id": u.ID, "name": u.Name, "login": u.Login, "avatarUrl": u.AvatarURL, "online": u.Online, "lastSeen": u.LastSeen})
	}
	c.JSON(200, res)
}

// @Summary Отметить комнату как прочитанную
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID комнаты"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /rooms/{id}/read [post]
func (h *Handler) MarkRoomRead(c *gin.Context) {
	// [COPILOT-END]
	roomID := c.Param("id")
	var m models.RoomMember
	if err := h.db.Where("room_id = ? AND user_id = ?", roomID, uid(c)).First(&m).Error; err != nil {
		respondErr(c, 404, "not a member")
		return
	}
	now := time.Now()
	h.db.Model(&m).Update("last_read_at", &now)
	c.JSON(200, gin.H{"ok": true})
}
