package entity

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
