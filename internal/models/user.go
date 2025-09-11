package models

import "time"

type User struct {
	ID           string                 `json:"id"`
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone"`
	FullName     string                 `json:"full_name"`
	Role         string                 `json:"role"`
	Status       string                 `json:"status"`
	PasswordHash string                 `json:"-"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
