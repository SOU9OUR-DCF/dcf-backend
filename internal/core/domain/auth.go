package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuthToken struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type BaseRegisterRequest struct {
	Username string   `json:"username" binding:"required"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	UserType UserType `json:"user_type" binding:"required,oneof=regular restaurant volunteer"`
}

type RestaurantRegisterRequest struct {
	BaseRegisterRequest
	Name          string `json:"name" binding:"required"`
	OwnerName     string `json:"owner_name" binding:"required"`
	BusinessEmail string `json:"business_email" binding:"required,email"`
	ContactNumber string `json:"contact_number" binding:"required"`
	Address       string `json:"address" binding:"required"`
}

type VolunteerRegisterRequest struct {
	BaseRegisterRequest
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Address     string `json:"address" binding:"required"`
}

type AuthResponse struct {
	ExpiresAt time.Time   `json:"expires_at"`
	User      User        `json:"user"`
	Profile   interface{} `json:"profile,omitempty"` // Will contain the specific profile data
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type Token string

func (t Token) String() string {
	return string(t)
}

func NewToken(token string) Token {
	return Token(token)
}
