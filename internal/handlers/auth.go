package handlers

import (
	"strings"
	"time"

	"LinkUp/internal/auth"
	apiErrors "LinkUp/internal/err"
	"LinkUp/internal/models"
	"LinkUp/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db         *gorm.DB
	uploadDir  string
	staticBase string
	presence   *Presence
	rooms      *RoomHubs
}

// Auto-generated swagger comments for New
// @Summary Auto-generated summary for New
// @Description Auto-generated description for New — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func New(db *gorm.DB, uploadDir, staticBase string) *Handler {
	return &Handler{db: db, uploadDir: uploadDir, staticBase: staticBase, presence: NewPresence(), rooms: NewRoomHubs()}
}

type registerReq struct {
	Login     string `json:"login" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	AvatarURL string `json:"avatarUrl"`
}

// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Данные пользователя"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	var req registerReq
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		apiErr := apiErrors.NewAPIError("Register.BindJSON", bindErr, "invalid body", 400)
		apiErrors.LogAndRespondAPI(c, apiErr, "Invalid request body.")
		return
	}

	req.Login = strings.TrimSpace(req.Login)
	if req.Login == "" || req.Password == "" || req.Name == "" {
		apiErr := apiErrors.NewAPIError("Register.Validate", nil, "login, password, name required", 400)
		apiErrors.LogAndRespondAPI(c, apiErr, "Login, password, and name are required.")
		return
	}

	hash, hashErr := utils.HashPassword(req.Password)
	if hashErr != nil {
		apiErr := apiErrors.NewAPIError("Register.HashPassword", hashErr, "failed to hash password", 500)
		apiErrors.LogAndRespondAPI(c, apiErr, "Internal server error.")
		return
	}

	u := models.User{
		Login:     req.Login,
		Password:  hash,
		Name:      req.Name,
		AvatarURL: req.AvatarURL,
	}

	if createErr := h.db.Create(&u).Error; createErr != nil {
		apiErr := apiErrors.NewAPIError("Register.CreateUser", createErr, "user exists?", 409)
		apiErrors.LogAndRespondAPI(c, apiErr, "User already exists.")
		return
	}

	tok, tokenErr := auth.GenerateToken(u.ID)
	if tokenErr != nil {
		apiErr := apiErrors.NewAPIError("Register.GenerateToken", tokenErr, "failed to generate token", 500)
		apiErrors.LogAndRespondAPI(c, apiErr, "Internal server error.")
		return
	}

	c.JSON(201, gin.H{
		"token": tok,
		"user": gin.H{
			"id":        u.ID,
			"login":     u.Login,
			"name":      u.Name,
			"avatarUrl": u.AvatarURL,
		},
	})
}

type loginReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Авторизация пользователя
// @Description Вход в систему с логином и паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Данные для входа"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {

	var req loginReq
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		apiErr := apiErrors.NewAPIError("Login.BindJSON", bindErr, "invalid body", 400)
		apiErrors.LogAndRespondAPI(c, apiErr, "Invalid request body.")
		return
	}
	var u models.User
	if findErr := h.db.Where("login = ?", req.Login).First(&u).Error; findErr != nil {
		apiErr := apiErrors.NewAPIError("Login.FindUser", findErr, "invalid credentials", 401)
		apiErrors.LogAndRespondAPI(c, apiErr, "Invalid credentials.")
		return
	}
	if !utils.CheckPassword(u.Password, req.Password) {
		apiErr := apiErrors.NewAPIError("Login.CheckPassword", nil, "invalid credentials", 401)
		apiErrors.LogAndRespondAPI(c, apiErr, "Invalid credentials.")
		return
	}
	now := time.Now()
	u.LastSeen = &now
	if updateErr := h.db.Model(&u).Update("last_seen", u.LastSeen).Error; updateErr != nil {
		apiErr := apiErrors.NewAPIError("Login.UpdateLastSeen", updateErr, "failed to update last seen", 500)
		apiErrors.LogAndRespondAPI(c, apiErr, "Internal server error.")
		return
	}
	tok, tokenErr := auth.GenerateToken(u.ID)
	if tokenErr != nil {
		apiErr := apiErrors.NewAPIError("Login.GenerateToken", tokenErr, "failed to generate token", 500)
		apiErrors.LogAndRespondAPI(c, apiErr, "Internal server error.")
		return
	}
	c.JSON(200, gin.H{"token": tok, "user": sanitizeUser(u)})

}

// Auto-generated swagger comments for sanitizeUser
// @Summary Auto-generated summary for sanitizeUser
// @Description Auto-generated description for sanitizeUser — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func sanitizeUser(u models.User) gin.H {
	return gin.H{"id": u.ID, "login": u.Login, "name": u.Name, "avatarUrl": u.AvatarURL, "lastSeen": u.LastSeen}
}

// @Summary Получить профиль текущего пользователя
// @Description Возвращает информацию о текущем авторизованном пользователе
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user/me [get]
func (h *Handler) Me(c *gin.Context) {
	var u models.User
	if findErr := h.db.First(&u, uid(c)).Error; findErr != nil {
		apiErr := apiErrors.NewAPIError("Me.FindUser", findErr, "not found", 404)
		apiErrors.LogAndRespondAPI(c, apiErr, "User not found.")
		return
	}
	u.Online = h.presence.IsOnline(u.ID)
	c.JSON(200, sanitizeUser(u))
}

type updateProfileReq struct {
	Name      string `json:"name"`
	AvatarURL string `json:"avatarUrl"`
}

// @Summary Обновить профиль пользователя
// @Description Обновляет информацию профиля текущего пользователя
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body UpdateProfileRequest true "Данные для обновления"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /user/me [put]
func (h *Handler) UpdateProfile(c *gin.Context) {
	// [COPILOT-END]
	var req updateProfileReq
	if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
		apiErr := apiErrors.NewAPIError("UpdateProfile.BindJSON", bindErr, "invalid body", 400)
		apiErrors.LogAndRespondAPI(c, apiErr, "Invalid request body.")
		return
	}
	updates := map[string]interface{}{}
	if strings.TrimSpace(req.Name) != "" {
		updates["name"] = strings.TrimSpace(req.Name)
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}
	if len(updates) == 0 {
		c.JSON(200, gin.H{"ok": true})
		return
	}
	if updateErr := h.db.Model(&models.User{}).Where("id = ?", uid(c)).Updates(updates).Error; updateErr != nil {
		apiErr := apiErrors.NewAPIError("UpdateProfile.Updates", updateErr, "update failed", 400)
		apiErrors.LogAndRespondAPI(c, apiErr, "Failed to update profile.")
		return
	}
	var u models.User
	if findErr := h.db.First(&u, uid(c)).Error; findErr != nil {
		apiErr := apiErrors.NewAPIError("UpdateProfile.FindUser", findErr, "user not found after update", 404)
		apiErrors.LogAndRespondAPI(c, apiErr, "User not found after update.")
		return
	}
	c.JSON(200, sanitizeUser(u))
}
