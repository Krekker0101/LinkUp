package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"LinkUp/internal/models"
)

// Auto-generated swagger comments for OpenDefault
// @Summary Auto-generated summary for OpenDefault
// @Description Auto-generated description for OpenDefault — review and improve
// @Tags internal
// (internal function — not necessarily an HTTP handler)

func OpenDefault() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=5432 dbname=chat port=5432 sslmode=disable"
	}
	var db *gorm.DB
	var err error
	if strings.HasPrefix(dsn, "sqlite://") {
		path := strings.TrimPrefix(dsn, "sqlite://")
		os.MkdirAll(filepath.Dir(path), 0755)
		db, err = gorm.Open(sqlite.Open(path), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") || strings.Contains(dsn, "host=") {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	} else {
		return nil, errors.New("unsupported DATABASE_URL, use sqlite:// or postgres:// or DSN with host=")
	}
	if err != nil {
		return nil, err
	}

	// Set sane pool
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	return db, nil
}
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.RoomMember{},
		&models.Message{},
		&models.Reaction{},
		// Новые модели
		&models.Role{},
		&models.UserRole{},
		&models.Permission{},
		&models.TwoFactorAuth{},
		&models.Analytics{},
		&models.RichMessage{},
		&models.Poll{},
		&models.PollVote{},
		&models.NotificationSettings{},
		&models.Mention{},
		&models.FileStorage{},
		&models.FileVersion{},
		&models.AIAssistant{},
		&models.ContentModeration{},
		&models.Achievement{},
		&models.UserAchievement{},
		&models.UserLevel{},
		&models.GitHubIntegration{},
		&models.CalendarEvent{},
		&models.MusicRoom{},
		&models.ChatGame{},
		&models.PushSubscription{},
		&models.OfflineMessage{},
	)
}
