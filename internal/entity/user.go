package entity

import (
	"time"
)

// User represents the user entity stored in the database.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`             // Primary key, auto increment
	Name      string    `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=3,max=100"` // User's full name, required, min 3 chars, max 100 chars
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email" validate:"required,email"`  // Unique email address, required, must be a valid email
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`               // Created timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`               // Updated timestamp
}

// TableName explicitly sets the table name to "users"
func (User) TableName() string {
	return "users"
}
