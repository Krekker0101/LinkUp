package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"LinkUp/internal/models"

	"github.com/gin-gonic/gin"
)

// ==================== РОЛИ И РАЗРЕШЕНИЯ ====================

// @Summary Получить все роли
// @Description Возвращает список всех ролей в системе
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Role
// @Failure 401 {object} ErrorResponse
// @Router /admin/roles [get]
func (h *Handler) GetRoles(c *gin.Context) {
	if !h.hasPermission(c, "admin.roles.read") {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Insufficient permissions"})
		return
	}

	var roles []models.Role
	if err := h.db.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch roles"})
		return
	}

	c.JSON(http.StatusOK, roles)
}

// @Summary Создать роль
// @Description Создает новую роль в системе
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param role body models.Role true "Данные роли"
// @Success 201 {object} models.Role
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/roles [post]
func (h *Handler) CreateRole(c *gin.Context) {
	if !h.hasPermission(c, "admin.roles.create") {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Insufficient permissions"})
		return
	}

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	if err := h.db.Create(&role).Error; err != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Role already exists"})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// @Summary Назначить роль пользователю
// @Description Назначает роль пользователю в комнате или глобально
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param assignment body map[string]interface{} true "Данные назначения"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /admin/assign-role [post]
func (h *Handler) AssignRole(c *gin.Context) {
	if !h.hasPermission(c, "admin.roles.assign") {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Insufficient permissions"})
		return
	}

	var req struct {
		UserID    uint   `json:"userId" binding:"required"`
		RoleID    uint   `json:"roleId" binding:"required"`
		RoomID    *uint  `json:"roomId"`
		ExpiresAt *time.Time `json:"expiresAt"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	userRole := models.UserRole{
		UserID:    req.UserID,
		RoleID:    req.RoleID,
		RoomID:    req.RoomID,
		GrantedBy: uid(c),
		ExpiresAt: req.ExpiresAt,
	}

	if err := h.db.Create(&userRole).Error; err != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Role assignment failed"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{OK: true})
}

// ==================== ДВУХФАКТОРНАЯ АУТЕНТИФИКАЦИЯ ====================

// @Summary Настроить 2FA
// @Description Настраивает двухфакторную аутентификацию для пользователя
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/2fa/setup [post]
func (h *Handler) Setup2FA(c *gin.Context) {
	userID := uid(c)
	
	// Генерируем секрет для 2FA
	secret := generateSecret()
	backupCodes := generateBackupCodes()

	twoFA := models.TwoFactorAuth{
		UserID:      userID,
		Secret:      secret,
		Enabled:     false,
		BackupCodes: backupCodes,
	}

	if err := h.db.Create(&twoFA).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to setup 2FA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret":      secret,
		"backupCodes": backupCodes,
		"qrCode":      generateQRCode(secret, fmt.Sprintf("%d", userID)),
	})
}

// @Summary Подтвердить 2FA
// @Description Подтверждает настройку 2FA с помощью кода
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param code body map[string]string true "Код подтверждения"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/2fa/verify [post]
func (h *Handler) Verify2FA(c *gin.Context) {
	userID := uid(c)
	
	var req struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	var twoFA models.TwoFactorAuth
	if err := h.db.Where("user_id = ?", userID).First(&twoFA).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "2FA not configured"})
		return
	}

	if !verifyTOTPCode(twoFA.Secret, req.Code) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid code"})
		return
	}

	twoFA.Enabled = true
	twoFA.LastUsed = &time.Time{}
	*twoFA.LastUsed = time.Now()

	if err := h.db.Save(&twoFA).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to enable 2FA"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{OK: true})
}

// ==================== АНАЛИТИКА ====================

// @Summary Получить аналитику пользователя
// @Description Возвращает метрики использования пользователя
// @Tags analytics
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Analytics
// @Failure 401 {object} ErrorResponse
// @Router /analytics/me [get]
func (h *Handler) GetUserAnalytics(c *gin.Context) {
	userID := uid(c)

	var analytics models.Analytics
	if err := h.db.Where("user_id = ?", userID).First(&analytics).Error; err != nil {
		// Создаем новую запись аналитики если не существует
		analytics = models.Analytics{
			UserID:        userID,
			MessagesSent:  0,
			TimeSpent:     0,
			LastActive:    time.Now(),
			LoginCount:    0,
			FilesUploaded: 0,
			RoomsJoined:   0,
		}
		h.db.Create(&analytics)
	}

	c.JSON(http.StatusOK, analytics)
}

// @Summary Получить дашборд администратора
// @Description Возвращает общую статистику системы
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.AdminDashboard
// @Failure 401 {object} ErrorResponse
// @Router /admin/dashboard [get]
func (h *Handler) GetAdminDashboard(c *gin.Context) {
	if !h.hasPermission(c, "admin.dashboard.read") {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Insufficient permissions"})
		return
	}

	var dashboard models.AdminDashboard

	// Подсчитываем общее количество пользователей
	h.db.Model(&models.User{}).Count(&dashboard.TotalUsers)

	// Подсчитываем активных пользователей (за последние 24 часа)
	activeSince := time.Now().Add(-24 * time.Hour)
	h.db.Model(&models.User{}).Where("last_seen > ?", activeSince).Count(&dashboard.ActiveUsers)

	// Подсчитываем сообщения за сегодня
	today := time.Now().Truncate(24 * time.Hour)
	h.db.Model(&models.Message{}).Where("created_at >= ?", today).Count(&dashboard.MessagesToday)

	// Подсчитываем комнаты
	h.db.Model(&models.Room{}).Count(&dashboard.RoomsCount)

	// Подсчитываем загруженные файлы
	h.db.Model(&models.FileStorage{}).Count(&dashboard.FilesUploaded)

	// Подсчитываем использованное место
	h.db.Model(&models.FileStorage{}).Select("COALESCE(SUM(size), 0)").Scan(&dashboard.StorageUsed)

	dashboard.LastUpdated = time.Now()

	c.JSON(http.StatusOK, dashboard)
}

// ==================== РАСШИРЕННЫЕ СООБЩЕНИЯ ====================

// @Summary Создать опрос
// @Description Создает опрос в комнате
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID комнаты"
// @Param poll body map[string]interface{} true "Данные опроса"
// @Success 201 {object} models.Poll
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /rooms/{id}/polls [post]
func (h *Handler) CreatePoll(c *gin.Context) {
	roomIDStr := c.Param("id")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid room ID"})
		return
	}

	var req struct {
		Question       string    `json:"question" binding:"required"`
		Options        []string  `json:"options" binding:"required"`
		MultipleChoice bool      `json:"multipleChoice"`
		Anonymous      bool      `json:"anonymous"`
		ExpiresAt      *time.Time `json:"expiresAt"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Создаем сообщение
	message := models.Message{
		RoomID:   uint(roomID),
		UserID:   uid(c),
		Type:     "poll",
		Text:     req.Question,
	}

	if err := h.db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create message"})
		return
	}

	// Создаем опрос
	poll := models.Poll{
		MessageID:      message.ID,
		Question:       req.Question,
		Options:        req.Options,
		MultipleChoice: req.MultipleChoice,
		Anonymous:      req.Anonymous,
		ExpiresAt:      req.ExpiresAt,
		Votes:          make(map[string][]uint),
		TotalVotes:     0,
	}

	if err := h.db.Create(&poll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create poll"})
		return
	}

	// Отправляем через WebSocket
	h.rooms.Emit(uint(roomID), Event{
		Type: "poll_created",
		Payload: gin.H{
			"messageId": message.ID,
			"poll":      poll,
		},
	})

	c.JSON(http.StatusCreated, poll)
}

// @Summary Голосовать в опросе
// @Description Голосует в опросе
// @Tags messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID опроса"
// @Param vote body map[string]string true "Голос"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /polls/{id}/vote [post]
func (h *Handler) VotePoll(c *gin.Context) {
	pollIDStr := c.Param("id")
	pollID, err := strconv.ParseUint(pollIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid poll ID"})
		return
	}

	var req struct {
		Option string `json:"option" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	var poll models.Poll
	if err := h.db.First(&poll, pollID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Poll not found"})
		return
	}

	// Проверяем, не истек ли опрос
	if poll.ExpiresAt != nil && poll.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Poll has expired"})
		return
	}

	// Проверяем, существует ли опция
	optionExists := false
	for _, option := range poll.Options {
		if option == req.Option {
			optionExists = true
			break
		}
	}
	if !optionExists {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid option"})
		return
	}

	userID := uid(c)

	// Удаляем предыдущий голос если есть
	h.db.Where("poll_id = ? AND user_id = ?", pollID, userID).Delete(&models.PollVote{})

	// Добавляем новый голос
	vote := models.PollVote{
		PollID: uint(pollID),
		UserID: userID,
		Option: req.Option,
	}

	if err := h.db.Create(&vote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to vote"})
		return
	}

	// Обновляем статистику опроса
	poll.TotalVotes++
	if poll.Votes[req.Option] == nil {
		poll.Votes[req.Option] = []uint{}
	}
	poll.Votes[req.Option] = append(poll.Votes[req.Option], userID)

	h.db.Save(&poll)

	// Отправляем обновление через WebSocket
	var message models.Message
	h.db.First(&message, poll.MessageID)
	h.rooms.Emit(message.RoomID, Event{
		Type: "poll_updated",
		Payload: gin.H{
			"pollId": poll.ID,
			"votes":  poll.Votes,
			"totalVotes": poll.TotalVotes,
		},
	})

	c.JSON(http.StatusOK, SuccessResponse{OK: true})
}

// ==================== УПОМИНАНИЯ ====================

// @Summary Получить упоминания пользователя
// @Description Возвращает список упоминаний пользователя
// @Tags mentions
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Лимит" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} models.Mention
// @Failure 401 {object} ErrorResponse
// @Router /mentions [get]
func (h *Handler) GetMentions(c *gin.Context) {
	userID := uid(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var mentions []models.Mention
	if err := h.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&mentions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch mentions"})
		return
	}

	c.JSON(http.StatusOK, mentions)
}

// @Summary Отметить упоминание как прочитанное
// @Description Отмечает упоминание как прочитанное
// @Tags mentions
// @Security BearerAuth
// @Produce json
// @Param id path string true "ID упоминания"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /mentions/{id}/read [post]
func (h *Handler) MarkMentionRead(c *gin.Context) {
	mentionIDStr := c.Param("id")
	mentionID, err := strconv.ParseUint(mentionIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid mention ID"})
		return
	}

	userID := uid(c)

	if err := h.db.Model(&models.Mention{}).
		Where("id = ? AND user_id = ?", mentionID, userID).
		Update("read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{OK: true})
}

// ==================== ДОСТИЖЕНИЯ ====================

// @Summary Получить достижения пользователя
// @Description Возвращает список достижений пользователя
// @Tags achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.UserAchievement
// @Failure 401 {object} ErrorResponse
// @Router /achievements [get]
func (h *Handler) GetUserAchievements(c *gin.Context) {
	userID := uid(c)

	var achievements []models.UserAchievement
	if err := h.db.Where("user_id = ?", userID).
		Preload("Achievement").
		Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

// @Summary Получить уровень пользователя
// @Description Возвращает информацию об уровне пользователя
// @Tags achievements
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserLevel
// @Failure 401 {object} ErrorResponse
// @Router /level [get]
func (h *Handler) GetUserLevel(c *gin.Context) {
	userID := uid(c)

	var level models.UserLevel
	if err := h.db.Where("user_id = ?", userID).First(&level).Error; err != nil {
		// Создаем новый уровень если не существует
		level = models.UserLevel{
			UserID:      userID,
			Level:       1,
			Experience:  0,
			Badges:      []string{},
			Title:       "Новичок",
			NextLevelXP: 100,
		}
		h.db.Create(&level)
	}

	c.JSON(http.StatusOK, level)
}

// ==================== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ====================

// hasPermission проверяет, есть ли у пользователя определенное разрешение
func (h *Handler) hasPermission(c *gin.Context, permission string) bool {
	userID := uid(c)
	
	// Проверяем глобальные роли
	var userRoles []models.UserRole
	h.db.Where("user_id = ? AND room_id IS NULL", userID).Find(&userRoles)
	
	for _, userRole := range userRoles {
		var role models.Role
		if err := h.db.First(&role, userRole.RoleID).Error; err != nil {
			continue
		}
		
		for _, perm := range role.Permissions {
			if perm == permission || perm == "admin.*" {
				return true
			}
		}
	}
	
	return false
}

// generateSecret генерирует секрет для 2FA
func generateSecret() string {
	return fmt.Sprintf("%016d", time.Now().UnixNano()%10000000000000000)
}

// generateBackupCodes генерирует резервные коды для 2FA
func generateBackupCodes() []string {
	codes := make([]string, 10)
	for i := range codes {
		codes[i] = fmt.Sprintf("%08d", i*1000000+i*100000+i*10000+i*1000+i*100+i*10+i)
	}
	return codes
}

// generateQRCode генерирует QR код для 2FA
func generateQRCode(secret, userID string) string {
	// В реальном приложении здесь будет генерация QR кода
	return fmt.Sprintf("otpauth://totp/LinkUp:%s?secret=%s&issuer=LinkUp", userID, secret)
}

// verifyTOTPCode проверяет TOTP код
func verifyTOTPCode(secret, code string) bool {
	// В реальном приложении здесь будет проверка TOTP кода
	// Для демонстрации принимаем любой код длиной 6 цифр
	return len(code) == 6 && strings.ContainsAny(code, "0123456789")
}
