package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"LinkUp/internal/auth"
	"LinkUp/internal/handlers"
	"LinkUp/internal/storage"

	_ "LinkUp/docs" // Импорт для swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title LinkUp API
// @version 1.0
// @description API для профессионального чат-приложения LinkUp.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// Auto-generated swagger comments for Run
// @Summary Auto-generated summary for Run
// @Description Auto-generated description for Run — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := storage.OpenDefault()
	if err != nil {
		return err
	}
	if err := storage.AutoMigrate(db); err != nil {
		return err
	}

	r := gin.Default()


	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	
	r.Use(func(c *gin.Context) {
		origins := os.Getenv("CORS_ORIGINS")
		if origins == "" {
			origins = "*"
		}
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	})

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	_ = os.MkdirAll(uploadDir, 0755)

	staticBase := os.Getenv("STATIC_BASE_URL")
	if staticBase == "" {
		staticBase = fmt.Sprintf("http://localhost:%s", port)
	}
	r.Static("/uploads", uploadDir)

	h := handlers.New(db, uploadDir, staticBase)

	
	// @Summary Проверка здоровья сервера
	// @Description Возвращает статус сервера
	// @Tags system
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true, "time": time.Now()})
	})
	// @Summary Регистрация пользователя
	// @Description Создает нового пользователя в системе
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param user body handlers.RegisterRequest true "Данные пользователя"
	// @Success 201 {object} handlers.AuthResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 409 {object} handlers.ErrorResponse
	// @Router /register [post]
	r.POST("/register", h.Register)

	// @Summary Авторизация пользователя
	// @Description Вход в систему с логином и паролем
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param credentials body handlers.LoginRequest true "Данные для входа"
	// @Success 200 {object} handlers.AuthResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /login [post]
	r.POST("/login", h.Login)

	pr := r.Group("")
	pr.Use(auth.JWTMiddleware())

	// @Summary Получить профиль текущего пользователя
	// @Description Возвращает информацию о текущем авторизованном пользователе
	// @Tags user
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} handlers.UserResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /user/me [get]
	pr.GET("/user/me", h.Me)

	// @Summary Обновить профиль пользователя
	// @Description Обновляет информацию профиля текущего пользователя
	// @Tags user
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param user body handlers.UpdateProfileRequest true "Данные для обновления"
	// @Success 200 {object} handlers.UserResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /user/me [put]
	pr.PUT("/user/me", h.UpdateProfile)

	
	// @Summary Получить список комнат
	// @Description Возвращает список всех доступных комнат с количеством непрочитанных сообщений
	// @Tags rooms
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {array} handlers.RoomResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /rooms [get]
	pr.GET("/rooms", h.ListRooms)

	// @Summary Создать комнату
	// @Description Создает новую комнату для общения
	// @Tags rooms
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param room body handlers.CreateRoomRequest true "Данные комнаты"
	// @Success 201 {object} handlers.RoomResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /rooms [post]
	pr.POST("/rooms", h.CreateRoom)

	// @Summary Присоединиться к комнате
	// @Description Добавляет текущего пользователя в указанную комнату
	// @Tags rooms
	// @Security BearerAuth
	// @Produce json
	// @Param id path string true "ID комнаты"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 404 {object} handlers.ErrorResponse
	// @Router /rooms/{id}/join [post]
	pr.POST("/rooms/:id/join", h.JoinRoom)

	// @Summary Покинуть комнату
	// @Tags rooms
	// @Security BearerAuth
	// @Produce json
	// @Param id path string true "ID комнаты"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 404 {object} handlers.ErrorResponse
	// @Router /rooms/{id}/leave [post]
	pr.POST("/rooms/:id/leave", h.LeaveRoom)

	// @Summary Список пользователей в комнате
	// @Tags rooms
	// @Security BearerAuth
	// @Produce json
	// @Param id path string true "ID комнаты"
	// @Success 200 {array} handlers.UserResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 404 {object} handlers.ErrorResponse
	// @Router /rooms/{id}/users [get]
	pr.GET("/rooms/:id/users", h.RoomMembers)

	// @Summary История сообщений комнаты
	// @Tags messages
	// @Security BearerAuth
	// @Produce json
	// @Param id path string true "ID комнаты"
	// @Param limit query int false "Лимит сообщений" default(50)
	// @Param offset query int false "Смещение" default(0)
	// @Success 200 {array} handlers.MessageResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 404 {object} handlers.ErrorResponse
	// @Router /rooms/{id}/history [get]
	pr.GET("/rooms/:id/history", h.MessageHistory)

	// @Summary Отправить сообщение
	// @Tags messages
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param id path string true "ID комнаты"
	// @Param message body handlers.SendMessageRequest true "Текст сообщения"
	// @Success 201 {object} handlers.MessageResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 404 {object} handlers.ErrorResponse
	// @Router /rooms/{id}/messages [post]
	pr.POST("/rooms/:id/messages", h.SendMessageREST)

	// @Summary Добавить реакцию к сообщению
	// @Tags messages
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param id path string true "ID сообщения"
	// @Param reaction body handlers.AddReactionRequest true "Реакция"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 404 {object} handlers.ErrorResponse
	// @Router /messages/{id}/reactions [post]
	pr.POST("/messages/:id/reactions", h.AddReaction)

	// @Summary Удалить реакцию
	// @Tags messages
	// @Security BearerAuth
	// @Produce json
	// @Param id path string true "ID сообщения"
	// @Param reaction path string true "Тип реакции"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 404 {object} handlers.ErrorResponse
	// @Router /messages/{id}/reactions/{reaction} [delete]
	pr.DELETE("/messages/:id/reactions/:reaction", h.RemoveReaction)

	
	// @Summary Загрузить файл
	// @Tags upload
	// @Security BearerAuth
	// @Accept multipart/form-data
	// @Produce json
	// @Param file formData file true "Файл для загрузки"
	// @Success 200 {object} handlers.UploadResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Failure 413 {object} handlers.ErrorResponse
	// @Router /upload [post]
	pr.POST("/upload", h.Upload)

	
	// @Summary Поиск сообщений
	// @Tags search
	// @Security BearerAuth
	// @Produce json
	// @Param q query string true "Поисковый запрос"
	// @Param room_id query string false "ID комнаты"
	// @Param limit query int false "Лимит" default(20)
	// @Param offset query int false "Смещение" default(0)
	// @Success 200 {array} handlers.MessageResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /search [get]
	pr.GET("/search", h.SearchMessages)

	// ==================== РОЛИ И РАЗРЕШЕНИЯ ====================
	// @Summary Получить все роли
	// @Tags admin
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {array} models.Role
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /admin/roles [get]
	pr.GET("/admin/roles", h.GetRoles)

	// @Summary Создать роль
	// @Tags admin
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param role body handlers.RoleRequest true "Данные роли"
	// @Success 201 {object} models.Role
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /admin/roles [post]
	pr.POST("/admin/roles", h.CreateRole)

	// @Summary Назначить роль пользователю
	// @Tags admin
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param assignment body handlers.AssignRoleRequest true "Данные назначения"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /admin/assign-role [post]
	pr.POST("/admin/assign-role", h.AssignRole)

	// @Summary Получить дашборд администратора
	// @Tags admin
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} handlers.AdminDashboardResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /admin/dashboard [get]
	pr.GET("/admin/dashboard", h.GetAdminDashboard)

	// ==================== ДВУХФАКТОРНАЯ АУТЕНТИФИКАЦИЯ ====================
	// @Summary Настроить 2FA
	// @Tags auth
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Success 200 {object} handlers.Setup2FAResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /auth/2fa/setup [post]
	pr.POST("/auth/2fa/setup", h.Setup2FA)

	// @Summary Подтвердить 2FA
	// @Tags auth
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param code body handlers.Verify2FARequest true "Код подтверждения"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /auth/2fa/verify [post]
	pr.POST("/auth/2fa/verify", h.Verify2FA)

	// ==================== АНАЛИТИКА ====================
	// @Summary Получить аналитику пользователя
	// @Tags analytics
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} handlers.AnalyticsResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /analytics/me [get]
	pr.GET("/analytics/me", h.GetUserAnalytics)

	// ==================== ОПРОСЫ ====================
	// @Summary Создать опрос
	// @Tags messages
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param id path string true "ID комнаты"
	// @Param poll body handlers.CreatePollRequest true "Данные опроса"
	// @Success 201 {object} handlers.PollResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /rooms/{id}/polls [post]
	pr.POST("/rooms/:id/polls", h.CreatePoll)

	// @Summary Голосовать в опросе
	// @Tags messages
	// @Security BearerAuth
	// @Accept json
	// @Produce json
	// @Param id path string true "ID опроса"
	// @Param vote body handlers.VotePollRequest true "Голос"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 400 {object} handlers.ErrorResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /polls/{id}/vote [post]
	pr.POST("/polls/:id/vote", h.VotePoll)

	// ==================== УПОМИНАНИЯ ====================
	// @Summary Получить упоминания пользователя
	// @Tags mentions
	// @Security BearerAuth
	// @Produce json
	// @Param limit query int false "Лимит" default(20)
	// @Param offset query int false "Смещение" default(0)
	// @Success 200 {array} handlers.MentionResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /mentions [get]
	pr.GET("/mentions", h.GetMentions)

	// @Summary Отметить упоминание как прочитанное
	// @Tags mentions
	// @Security BearerAuth
	// @Produce json
	// @Param id path string true "ID упоминания"
	// @Success 200 {object} handlers.SuccessResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /mentions/{id}/read [post]
	pr.POST("/mentions/:id/read", h.MarkMentionRead)

	// ==================== ДОСТИЖЕНИЯ ====================
	// @Summary Получить достижения пользователя
	// @Tags achievements
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {array} handlers.AchievementResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /achievements [get]
	pr.GET("/achievements", h.GetUserAchievements)

	// @Summary Получить уровень пользователя
	// @Tags achievements
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} handlers.UserLevelResponse
	// @Failure 401 {object} handlers.ErrorResponse
	// @Router /level [get]
	pr.GET("/level", h.GetUserLevel)


	r.GET("/ws/rooms/:id", auth.UpgradeWithJWT(h.RoomWebSocket))

	addr := ":" + port
	log.Printf("🔥 LinkUp API listening on %s (CORS: %s)\n", addr, strings.TrimSpace(os.Getenv("CORS_ORIGINS")))
	return r.Run(addr)
}
