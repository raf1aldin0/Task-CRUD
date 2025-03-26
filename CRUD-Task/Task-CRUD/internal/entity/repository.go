package entity

import (
	"time"
)

type Repository struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
	URL       string    `json:"url"`
	AIEnabled bool      `json:"ai_enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
