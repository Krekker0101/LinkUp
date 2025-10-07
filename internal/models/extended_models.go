package models

import (
	"time"
)

// Role представляет роль пользователя в системе
type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string   `gorm:"uniqueIndex;size:64" json:"name"` // admin, moderator, user
	Description string   `gorm:"size:255" json:"description"`
	Permissions []string `gorm:"serializer:json" json:"permissions"`
	IsDefault   bool     `json:"isDefault"`
}

// UserRole связывает пользователей с ролями
type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID    uint       `gorm:"index;uniqueIndex:uniq_user_role" json:"userId"`
	RoleID    uint       `gorm:"index;uniqueIndex:uniq_user_role" json:"roleId"`
	RoomID    *uint      `gorm:"index" json:"roomId"` // null = глобальная роль
	GrantedBy uint       `json:"grantedBy"`
	ExpiresAt *time.Time `json:"expiresAt"`
}

// Permission представляет разрешение в системе
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string `gorm:"uniqueIndex;size:64" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Category    string `gorm:"size:32" json:"category"`
}

// TwoFactorAuth представляет настройки двухфакторной аутентификации
type TwoFactorAuth struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID      uint       `gorm:"uniqueIndex" json:"userId"`
	Secret      string     `gorm:"size:32" json:"secret"`
	Enabled     bool       `json:"enabled"`
	BackupCodes []string   `gorm:"serializer:json" json:"backupCodes"`
	LastUsed    *time.Time `json:"lastUsed"`
}

// Analytics представляет метрики использования пользователя
type Analytics struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID        uint      `gorm:"uniqueIndex" json:"userId"`
	MessagesSent  int64     `json:"messagesSent"`
	TimeSpent     int64     `json:"timeSpentMinutes"`
	LastActive    time.Time `json:"lastActive"`
	LoginCount    int64     `json:"loginCount"`
	FilesUploaded int64     `json:"filesUploaded"`
	RoomsJoined   int64     `json:"roomsJoined"`
}

// AdminDashboard представляет данные для дашборда администратора
type AdminDashboard struct {
	TotalUsers    int64     `json:"totalUsers"`
	ActiveUsers   int64     `json:"activeUsers"`
	MessagesToday int64     `json:"messagesToday"`
	RoomsCount    int64     `json:"roomsCount"`
	FilesUploaded int64     `json:"filesUploaded"`
	StorageUsed   int64     `json:"storageUsed"`
	LastUpdated   time.Time `json:"lastUpdated"`
}

// RichMessage представляет расширенное сообщение с форматированием
type RichMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	MessageID  uint                   `gorm:"uniqueIndex" json:"messageId"`
	Type       string                 `gorm:"size:16" json:"type"` // markdown, html, code, poll
	Content    string                 `gorm:"size:8000" json:"content"`
	Formatting string                 `gorm:"size:32" json:"formatting"`
	Metadata   map[string]interface{} `gorm:"serializer:json" json:"metadata"`
}

// Poll представляет опрос в сообщении
type Poll struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	MessageID      uint              `gorm:"uniqueIndex" json:"messageId"`
	Question       string            `gorm:"size:500" json:"question"`
	Options        []string          `gorm:"serializer:json" json:"options"`
	Votes          map[string][]uint `gorm:"serializer:json" json:"votes"`
	ExpiresAt      *time.Time        `json:"expiresAt"`
	MultipleChoice bool              `json:"multipleChoice"`
	Anonymous      bool              `json:"anonymous"`
	TotalVotes     int64             `json:"totalVotes"`
}

// PollVote представляет голос в опросе
type PollVote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	PollID uint   `gorm:"index;uniqueIndex:uniq_poll_user" json:"pollId"`
	UserID uint   `gorm:"index;uniqueIndex:uniq_poll_user" json:"userId"`
	Option string `gorm:"size:100" json:"option"`
}

// NotificationSettings представляет настройки уведомлений пользователя
type NotificationSettings struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID       uint                              `gorm:"uniqueIndex" json:"userId"`
	EmailEnabled bool                              `json:"emailEnabled"`
	PushEnabled  bool                              `json:"pushEnabled"`
	MentionOnly  bool                              `json:"mentionOnly"`
	QuietHours   string                            `gorm:"size:20" json:"quietHours"` // "22:00-08:00"
	Keywords     []string                          `gorm:"serializer:json" json:"keywords"`
	RoomSettings map[uint]RoomNotificationSettings `gorm:"serializer:json" json:"roomSettings"`
}

// RoomNotificationSettings представляет настройки уведомлений для конкретной комнаты
type RoomNotificationSettings struct {
	RoomID      uint     `json:"roomId"`
	Muted       bool     `json:"muted"`
	MentionOnly bool     `json:"mentionOnly"`
	Keywords    []string `json:"keywords"`
}

// Mention представляет упоминание пользователя в сообщении
type Mention struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	MessageID   uint `gorm:"index" json:"messageId"`
	UserID      uint `gorm:"index" json:"userId"`
	MentionedBy uint `gorm:"index" json:"mentionedBy"`
	Read        bool `json:"read"`
}

// FileStorage представляет расширенную систему файлов
type FileStorage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID        uint     `gorm:"index" json:"userId"`
	RoomID        *uint    `gorm:"index" json:"roomId"`
	Filename      string   `gorm:"size:255" json:"filename"`
	OriginalName  string   `gorm:"size:255" json:"originalName"`
	Size          int64    `json:"size"`
	MimeType      string   `gorm:"size:100" json:"mimeType"`
	StorageType   string   `gorm:"size:32" json:"storageType"` // local, s3, gcs
	URL           string   `gorm:"size:500" json:"url"`
	Thumbnail     string   `gorm:"size:500" json:"thumbnail"`
	Hash          string   `gorm:"size:64" json:"hash"`
	Tags          []string `gorm:"serializer:json" json:"tags"`
	IsPublic      bool     `json:"isPublic"`
	DownloadCount int64    `json:"downloadCount"`
}

// FileVersion представляет версии файлов
type FileVersion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	FileID    uint   `gorm:"index" json:"fileId"`
	Version   int    `json:"version"`
	Filename  string `gorm:"size:255" json:"filename"`
	Size      int64  `json:"size"`
	URL       string `gorm:"size:500" json:"url"`
	Hash      string `gorm:"size:64" json:"hash"`
	CreatedBy uint   `json:"createdBy"`
}

// AIAssistant представляет AI-помощника для комнат
type AIAssistant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	RoomID      uint                   `gorm:"uniqueIndex" json:"roomId"`
	Name        string                 `gorm:"size:100" json:"name"`
	Personality string                 `gorm:"size:32" json:"personality"`
	Commands    []string               `gorm:"serializer:json" json:"commands"`
	Enabled     bool                   `json:"enabled"`
	Model       string                 `gorm:"size:50" json:"model"`
	APIKey      string                 `gorm:"size:255" json:"apiKey"`
	Settings    map[string]interface{} `gorm:"serializer:json" json:"settings"`
}

// ContentModeration представляет автоматическую модерацию контента
type ContentModeration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	MessageID     uint     `gorm:"uniqueIndex" json:"messageId"`
	ToxicityScore float64  `json:"toxicityScore"`
	Categories    []string `gorm:"serializer:json" json:"categories"`
	Action        string   `gorm:"size:32" json:"action"` // warn, delete, mute
	Reviewed      bool     `json:"reviewed"`
	ReviewedBy    *uint    `json:"reviewedBy"`
	ReviewNote    string   `gorm:"size:500" json:"reviewNote"`
}

// Achievement представляет достижение пользователя
type Achievement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Name        string                 `gorm:"uniqueIndex;size:100" json:"name"`
	Description string                 `gorm:"size:500" json:"description"`
	Icon        string                 `gorm:"size:100" json:"icon"`
	Points      int                    `json:"points"`
	Category    string                 `gorm:"size:32" json:"category"`
	Rarity      string                 `gorm:"size:16" json:"rarity"` // common, rare, epic, legendary
	Condition   map[string]interface{} `gorm:"serializer:json" json:"condition"`
	Enabled     bool                   `json:"enabled"`
}

// UserAchievement представляет достижение пользователя
type UserAchievement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID        uint       `gorm:"index;uniqueIndex:uniq_user_achievement" json:"userId"`
	AchievementID uint       `gorm:"index;uniqueIndex:uniq_user_achievement" json:"achievementId"`
	Progress      int        `json:"progress"`
	Completed     bool       `json:"completed"`
	CompletedAt   *time.Time `json:"completedAt"`
}

// UserLevel представляет уровень пользователя
type UserLevel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID      uint     `gorm:"uniqueIndex" json:"userId"`
	Level       int      `json:"level"`
	Experience  int64    `json:"experience"`
	Badges      []string `gorm:"serializer:json" json:"badges"`
	Title       string   `gorm:"size:100" json:"title"`
	NextLevelXP int64    `json:"nextLevelXp"`
}

// GitHubIntegration представляет интеграцию с GitHub
type GitHubIntegration struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID        uint       `gorm:"uniqueIndex" json:"userId"`
	AccessToken   string     `gorm:"size:255" json:"accessToken"`
	RefreshToken  string     `gorm:"size:255" json:"refreshToken"`
	Repositories  []string   `gorm:"serializer:json" json:"repositories"`
	Notifications bool       `json:"notifications"`
	WebhookSecret string     `gorm:"size:255" json:"webhookSecret"`
	LastSync      *time.Time `json:"lastSync"`
}

// CalendarEvent представляет событие календаря
type CalendarEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	RoomID         uint      `gorm:"index" json:"roomId"`
	UserID         uint      `gorm:"index" json:"userId"`
	Title          string    `gorm:"size:200" json:"title"`
	Description    string    `gorm:"size:1000" json:"description"`
	StartTime      time.Time `json:"startTime"`
	EndTime        time.Time `json:"endTime"`
	Attendees      []uint    `gorm:"serializer:json" json:"attendees"`
	Location       string    `gorm:"size:200" json:"location"`
	Recurring      bool      `json:"recurring"`
	RecurrenceRule string    `gorm:"size:100" json:"recurrenceRule"`
}

// MusicRoom представляет музыкальную комнату
type MusicRoom struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	RoomID      uint         `gorm:"uniqueIndex" json:"roomId"`
	CurrentSong string       `gorm:"size:200" json:"currentSong"`
	Queue       []string     `gorm:"serializer:json" json:"queue"`
	Playing     bool         `json:"playing"`
	Volume      int          `json:"volume"`
	DJ          *uint        `json:"dj"`
	Playlist    []MusicTrack `gorm:"serializer:json" json:"playlist"`
}

// MusicTrack представляет трек в музыкальной комнате
type MusicTrack struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Artist    string `json:"artist"`
	Duration  int    `json:"duration"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

// ChatGame представляет игру в чате
type ChatGame struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	RoomID   uint                   `gorm:"index" json:"roomId"`
	Type     string                 `gorm:"size:32" json:"type"` // trivia, word, drawing, quiz
	Players  []uint                 `gorm:"serializer:json" json:"players"`
	Status   string                 `gorm:"size:32" json:"status"` // waiting, playing, finished
	Score    map[uint]int           `gorm:"serializer:json" json:"score"`
	Settings map[string]interface{} `gorm:"serializer:json" json:"settings"`
	Winner   *uint                  `json:"winner"`
	EndedAt  *time.Time             `json:"endedAt"`
}

// PushSubscription представляет подписку на push уведомления
type PushSubscription struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID     uint              `gorm:"index" json:"userId"`
	Endpoint   string            `gorm:"size:500" json:"endpoint"`
	Keys       map[string]string `gorm:"serializer:json" json:"keys"`
	DeviceInfo string            `gorm:"size:200" json:"deviceInfo"`
	Active     bool              `json:"active"`
	LastUsed   *time.Time        `json:"lastUsed"`
}

// OfflineMessage представляет сообщение для офлайн режима
type OfflineMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	UserID    uint                   `gorm:"index" json:"userId"`
	MessageID uint                   `gorm:"index" json:"messageId"`
	Synced    bool                   `json:"synced"`
	SyncedAt  *time.Time             `json:"syncedAt"`
	Data      map[string]interface{} `gorm:"serializer:json" json:"data"`
}
