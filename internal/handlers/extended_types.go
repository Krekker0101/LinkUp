package handlers

import "time"

// ==================== РОЛИ И РАЗРЕШЕНИЯ ====================

// RoleRequest представляет запрос для создания роли
type RoleRequest struct {
	Name        string   `json:"name" binding:"required" example:"moderator"`
	Description string   `json:"description" example:"Модератор комнаты"`
	Permissions []string `json:"permissions" example:"[\"messages.delete\", \"users.mute\"]"`
	IsDefault   bool     `json:"isDefault" example:"false"`
}

// AssignRoleRequest представляет запрос для назначения роли
type AssignRoleRequest struct {
	UserID    uint       `json:"userId" binding:"required" example:"1"`
	RoleID    uint       `json:"roleId" binding:"required" example:"2"`
	RoomID    *uint      `json:"roomId" example:"1"`
	ExpiresAt *time.Time `json:"expiresAt" example:"2024-12-31T23:59:59Z"`
}

// ==================== ДВУХФАКТОРНАЯ АУТЕНТИФИКАЦИЯ ====================

// Setup2FAResponse представляет ответ при настройке 2FA
type Setup2FAResponse struct {
	Secret      string   `json:"secret" example:"ABCD1234EFGH5678"`
	BackupCodes []string `json:"backupCodes" example:"[\"12345678\", \"87654321\"]"`
	QRCode      string   `json:"qrCode" example:"otpauth://totp/LinkUp:1?secret=ABCD1234EFGH5678&issuer=LinkUp"`
}

// Verify2FARequest представляет запрос для подтверждения 2FA
type Verify2FARequest struct {
	Code string `json:"code" binding:"required" example:"123456"`
}

// ==================== АНАЛИТИКА ====================

// AnalyticsResponse представляет ответ с аналитикой пользователя
type AnalyticsResponse struct {
	UserID        uint      `json:"userId" example:"1"`
	MessagesSent  int64     `json:"messagesSent" example:"150"`
	TimeSpent     int64     `json:"timeSpentMinutes" example:"1200"`
	LastActive    time.Time `json:"lastActive" example:"2024-01-15T10:30:00Z"`
	LoginCount    int64     `json:"loginCount" example:"25"`
	FilesUploaded int64     `json:"filesUploaded" example:"10"`
	RoomsJoined   int64     `json:"roomsJoined" example:"5"`
}

// AdminDashboardResponse представляет ответ дашборда администратора
type AdminDashboardResponse struct {
	TotalUsers     int64     `json:"totalUsers" example:"1000"`
	ActiveUsers    int64     `json:"activeUsers" example:"150"`
	MessagesToday  int64     `json:"messagesToday" example:"5000"`
	RoomsCount     int64     `json:"roomsCount" example:"50"`
	FilesUploaded  int64     `json:"filesUploaded" example:"200"`
	StorageUsed    int64     `json:"storageUsed" example:"1073741824"`
	LastUpdated    time.Time `json:"lastUpdated" example:"2024-01-15T10:30:00Z"`
}

// ==================== ОПРОСЫ ====================

// CreatePollRequest представляет запрос для создания опроса
type CreatePollRequest struct {
	Question       string     `json:"question" binding:"required" example:"Какой язык программирования предпочитаете?"`
	Options        []string   `json:"options" binding:"required" example:"[\"Go\", \"Python\", \"JavaScript\", \"Java\"]"`
	MultipleChoice bool       `json:"multipleChoice" example:"false"`
	Anonymous      bool       `json:"anonymous" example:"false"`
	ExpiresAt      *time.Time `json:"expiresAt" example:"2024-01-20T23:59:59Z"`
}

// PollResponse представляет ответ с данными опроса
type PollResponse struct {
	ID             uint                   `json:"id" example:"1"`
	MessageID      uint                   `json:"messageId" example:"123"`
	Question       string                 `json:"question" example:"Какой язык программирования предпочитаете?"`
	Options        []string               `json:"options" example:"[\"Go\", \"Python\", \"JavaScript\", \"Java\"]"`
	Votes          map[string][]uint      `json:"votes" example:"{\"Go\":[1,2], \"Python\":[3,4,5]}"`
	ExpiresAt      *time.Time             `json:"expiresAt" example:"2024-01-20T23:59:59Z"`
	MultipleChoice bool                   `json:"multipleChoice" example:"false"`
	Anonymous      bool                   `json:"anonymous" example:"false"`
	TotalVotes     int64                  `json:"totalVotes" example:"5"`
}

// VotePollRequest представляет запрос для голосования
type VotePollRequest struct {
	Option string `json:"option" binding:"required" example:"Go"`
}

// ==================== УПОМИНАНИЯ ====================

// MentionResponse представляет ответ с упоминанием
type MentionResponse struct {
	ID         uint      `json:"id" example:"1"`
	MessageID  uint      `json:"messageId" example:"123"`
	UserID     uint      `json:"userId" example:"1"`
	MentionedBy uint     `json:"mentionedBy" example:"2"`
	Read       bool      `json:"read" example:"false"`
	CreatedAt  time.Time `json:"createdAt" example:"2024-01-15T10:30:00Z"`
}

// ==================== ДОСТИЖЕНИЯ ====================

// AchievementResponse представляет ответ с достижением
type AchievementResponse struct {
	ID          uint      `json:"id" example:"1"`
	Name        string    `json:"name" example:"Первое сообщение"`
	Description string    `json:"description" example:"Отправить первое сообщение"`
	Icon        string    `json:"icon" example:"💬"`
	Points      int       `json:"points" example:"10"`
	Category    string    `json:"category" example:"messages"`
	Rarity      string    `json:"rarity" example:"common"`
	Completed   bool      `json:"completed" example:"true"`
	CompletedAt *time.Time `json:"completedAt" example:"2024-01-15T10:30:00Z"`
}

// UserLevelResponse представляет ответ с уровнем пользователя
type UserLevelResponse struct {
	UserID      uint      `json:"userId" example:"1"`
	Level       int       `json:"level" example:"5"`
	Experience  int64     `json:"experience" example:"2500"`
	Badges      []string  `json:"badges" example:"[\"💬\", \"🏆\", \"⭐\"]"`
	Title       string    `json:"title" example:"Опытный пользователь"`
	NextLevelXP int64     `json:"nextLevelXp" example:"3000"`
}

// ==================== ФАЙЛЫ ====================

// FileUploadResponse представляет ответ при загрузке файла
type FileUploadResponse struct {
	ID           uint   `json:"id" example:"1"`
	Filename     string `json:"filename" example:"document.pdf"`
	OriginalName string `json:"originalName" example:"My Document.pdf"`
	Size         int64  `json:"size" example:"1048576"`
	MimeType     string `json:"mimeType" example:"application/pdf"`
	URL          string `json:"url" example:"http://localhost:8080/uploads/1_1642234567890.pdf"`
	Thumbnail    string `json:"thumbnail" example:"http://localhost:8080/uploads/thumbs/1_1642234567890.jpg"`
	Hash         string `json:"hash" example:"a1b2c3d4e5f6"`
	Tags         []string `json:"tags" example:"[\"document\", \"pdf\"]"`
	IsPublic     bool   `json:"isPublic" example:"false"`
}

// ==================== AI ПОМОЩНИК ====================

// AIAssistantResponse представляет ответ с AI помощником
type AIAssistantResponse struct {
	ID          uint                   `json:"id" example:"1"`
	RoomID      uint                   `json:"roomId" example:"1"`
	Name        string                 `json:"name" example:"Ассистент"`
	Personality string                 `json:"personality" example:"helpful"`
	Commands    []string               `json:"commands" example:"[\"help\", \"weather\", \"translate\"]"`
	Enabled     bool                   `json:"enabled" example:"true"`
	Model       string                 `json:"model" example:"gpt-3.5-turbo"`
	Settings    map[string]interface{} `json:"settings" example:"{\"temperature\": 0.7}"`
}

// ==================== ИГРЫ ====================

// ChatGameResponse представляет ответ с игрой в чате
type ChatGameResponse struct {
	ID       uint                   `json:"id" example:"1"`
	RoomID   uint                   `json:"roomId" example:"1"`
	Type     string                 `json:"type" example:"trivia"`
	Players  []uint                 `json:"players" example:"[1, 2, 3]"`
	Status   string                 `json:"status" example:"playing"`
	Score    map[uint]int           `json:"score" example:"{\"1\": 10, \"2\": 8, \"3\": 5}"`
	Settings map[string]interface{} `json:"settings" example:"{\"timeLimit\": 30}"`
	Winner   *uint                  `json:"winner" example:"1"`
}

// ==================== МУЗЫКА ====================

// MusicRoomResponse представляет ответ с музыкальной комнатой
type MusicRoomResponse struct {
	RoomID      uint         `json:"roomId" example:"1"`
	CurrentSong string       `json:"currentSong" example:"Song Title - Artist"`
	Queue       []string     `json:"queue" example:"[\"Song 1\", \"Song 2\"]"`
	Playing     bool         `json:"playing" example:"true"`
	Volume      int          `json:"volume" example:"50"`
	DJ          *uint        `json:"dj" example:"1"`
	Playlist    []MusicTrack `json:"playlist"`
}

// MusicTrack представляет трек в музыкальной комнате
type MusicTrack struct {
	ID        string `json:"id" example:"track_123"`
	Title     string `json:"title" example:"Song Title"`
	Artist    string `json:"artist" example:"Artist Name"`
	Duration  int    `json:"duration" example:"180"`
	URL       string `json:"url" example:"https://example.com/song.mp3"`
	Thumbnail string `json:"thumbnail" example:"https://example.com/thumb.jpg"`
}

// ==================== КАЛЕНДАРЬ ====================

// CalendarEventResponse представляет ответ с событием календаря
type CalendarEventResponse struct {
	ID            uint      `json:"id" example:"1"`
	RoomID        uint      `json:"roomId" example:"1"`
	UserID        uint      `json:"userId" example:"1"`
	Title         string    `json:"title" example:"Встреча команды"`
	Description   string    `json:"description" example:"Еженедельная встреча команды разработки"`
	StartTime     time.Time `json:"startTime" example:"2024-01-15T14:00:00Z"`
	EndTime       time.Time `json:"endTime" example:"2024-01-15T15:00:00Z"`
	Attendees     []uint    `json:"attendees" example:"[1, 2, 3]"`
	Location      string    `json:"location" example:"Конференц-зал A"`
	Recurring     bool      `json:"recurring" example:"true"`
	RecurrenceRule string   `json:"recurrenceRule" example:"FREQ=WEEKLY"`
}

// ==================== ИНТЕГРАЦИИ ====================

// GitHubIntegrationResponse представляет ответ с интеграцией GitHub
type GitHubIntegrationResponse struct {
	ID            uint       `json:"id" example:"1"`
	UserID        uint       `json:"userId" example:"1"`
	Repositories  []string   `json:"repositories" example:"[\"user/repo1\", \"user/repo2\"]"`
	Notifications bool       `json:"notifications" example:"true"`
	LastSync      *time.Time `json:"lastSync" example:"2024-01-15T10:30:00Z"`
}

// ==================== УВЕДОМЛЕНИЯ ====================

// NotificationSettingsResponse представляет ответ с настройками уведомлений
type NotificationSettingsResponse struct {
	UserID       uint                        `json:"userId" example:"1"`
	EmailEnabled bool                        `json:"emailEnabled" example:"true"`
	PushEnabled  bool                        `json:"pushEnabled" example:"true"`
	MentionOnly  bool                        `json:"mentionOnly" example:"false"`
	QuietHours   string                      `json:"quietHours" example:"22:00-08:00"`
	Keywords     []string                    `json:"keywords" example:"[\"важно\", \"срочно\"]"`
	RoomSettings map[uint]RoomNotificationSettings `json:"roomSettings"`
}

// RoomNotificationSettings представляет настройки уведомлений для комнаты
type RoomNotificationSettings struct {
	RoomID      uint     `json:"roomId" example:"1"`
	Muted       bool     `json:"muted" example:"false"`
	MentionOnly bool     `json:"mentionOnly" example:"false"`
	Keywords    []string `json:"keywords" example:"[\"важно\"]"`
}

// ==================== PUSH УВЕДОМЛЕНИЯ ====================

// PushSubscriptionRequest представляет запрос для подписки на push уведомления
type PushSubscriptionRequest struct {
	Endpoint   string            `json:"endpoint" binding:"required" example:"https://fcm.googleapis.com/fcm/send/..."`
	Keys       map[string]string `json:"keys" binding:"required" example:"{\"p256dh\":\"...\", \"auth\":\"...\"}"`
	DeviceInfo string            `json:"deviceInfo" example:"Chrome 120.0.0.0"`
}

// ==================== МОДЕРАЦИЯ ====================

// ContentModerationResponse представляет ответ с модерацией контента
type ContentModerationResponse struct {
	ID            uint      `json:"id" example:"1"`
	MessageID     uint      `json:"messageId" example:"123"`
	ToxicityScore float64   `json:"toxicityScore" example:"0.85"`
	Categories    []string  `json:"categories" example:"[\"harassment\", \"spam\"]"`
	Action        string    `json:"action" example:"warn"`
	Reviewed      bool      `json:"reviewed" example:"false"`
	ReviewedBy    *uint     `json:"reviewedBy" example:"1"`
	ReviewNote    string    `json:"reviewNote" example:"Автоматическое предупреждение"`
}
