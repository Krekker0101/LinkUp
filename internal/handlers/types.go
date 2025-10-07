package handlers

import "time"

// Swagger response types

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Login     string `json:"login" binding:"required" example:"john_doe"`
	Password  string `json:"password" binding:"required" example:"securepassword123"`
	Name      string `json:"name" binding:"required" example:"John Doe"`
	AvatarURL string `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Login    string `json:"login" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// UpdateProfileRequest represents the request body for profile update
type UpdateProfileRequest struct {
	Name      string `json:"name" example:"John Smith"`
	AvatarURL string `json:"avatarUrl" example:"https://example.com/new-avatar.jpg"`
}

// CreateRoomRequest represents the request body for room creation
type CreateRoomRequest struct {
	Slug      string `json:"slug" binding:"required" example:"general"`
	Name      string `json:"name" binding:"required" example:"General Chat"`
	IsPrivate bool   `json:"isPrivate" example:"false"`
}

// SendMessageRequest represents the request body for sending messages
type SendMessageRequest struct {
	Type     string `json:"type" example:"text"`
	Text     string `json:"text" example:"Hello everyone!"`
	ImageURL string `json:"imageUrl" example:"https://example.com/image.jpg"`
}

// AddReactionRequest represents the request body for adding reactions
type AddReactionRequest struct {
	Reaction string `json:"reaction" binding:"required" example:"👍"`
}

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponse `json:"user"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID        uint       `json:"id" example:"1"`
	Login     string     `json:"login" example:"john_doe"`
	Name      string     `json:"name" example:"John Doe"`
	AvatarURL string     `json:"avatarUrl" example:"https://example.com/avatar.jpg"`
	Online    bool       `json:"online" example:"true"`
	LastSeen  *time.Time `json:"lastSeen" example:"2024-01-15T10:30:00Z"`
}

// RoomResponse represents room data in responses
type RoomResponse struct {
	ID        uint   `json:"id" example:"1"`
	Slug      string `json:"slug" example:"general"`
	Name      string `json:"name" example:"General Chat"`
	IsPrivate bool   `json:"isPrivate" example:"false"`
	OwnerID   uint   `json:"ownerId" example:"1"`
	Unread    int64  `json:"unread" example:"5"`
}

// MessageResponse represents message data in responses
type MessageResponse struct {
	ID        uint              `json:"id" example:"1"`
	RoomID    uint              `json:"roomId" example:"1"`
	UserID    uint              `json:"userId" example:"1"`
	Type      string            `json:"type" example:"text"`
	Text      string            `json:"text" example:"Hello everyone!"`
	ImageURL  string            `json:"imageUrl" example:"https://example.com/image.jpg"`
	CreatedAt time.Time         `json:"createdAt" example:"2024-01-15T10:30:00Z"`
	Reactions map[string][]uint `json:"reactions"`
}

// UploadResponse represents the response for file upload
type UploadResponse struct {
	URL string `json:"url" example:"http://localhost:8080/uploads/1_1642234567890.jpg"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	OK bool `json:"ok" example:"true"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}
