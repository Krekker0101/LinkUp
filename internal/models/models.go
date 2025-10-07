package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Login     string     `gorm:"uniqueIndex;size:64" json:"login"`
	Password  string     `json:"-"`
	Name      string     `gorm:"size:120" json:"name"`
	AvatarURL string     `gorm:"size:255" json:"avatarUrl"`
	Online    bool       `gorm:"-" json:"online"`
	LastSeen  *time.Time `json:"lastSeen"`
}

type Room struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Slug      string `gorm:"uniqueIndex;size:64" json:"slug"`
	Name      string `gorm:"size:120" json:"name"`
	IsPrivate bool   `json:"isPrivate"`
	OwnerID   uint   `json:"ownerId"`
}

type RoomMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	RoomID     uint       `gorm:"index;uniqueIndex:uniq_room_user" json:"roomId"`
	UserID     uint       `gorm:"index;uniqueIndex:uniq_room_user" json:"userId"`
	LastReadAt *time.Time `json:"lastReadAt"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	RoomID   uint   `gorm:"index" json:"roomId"`
	UserID   uint   `gorm:"index" json:"userId"`
	Type     string `gorm:"size:16" json:"type"` // "text" | "image" | "system"
	Text     string `gorm:"size:4000" json:"text"`
	ImageURL string `gorm:"size:255" json:"imageUrl"`
}

type Reaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	MessageID uint      `gorm:"index;uniqueIndex:uniq_msg_user_react" json:"messageId"`
	UserID    uint      `gorm:"index;uniqueIndex:uniq_msg_user_react" json:"userId"`
	Reaction  string    `gorm:"size:32;uniqueIndex:uniq_msg_user_react" json:"reaction"`
}
