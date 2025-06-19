package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"unique;not null" json:"email"`
	FullName     string `gorm:"not null" json:"full_name"`
	PasswordHash string `gorm:"not null" json:"-"`
}

type Quote struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Text      string    `json:"text"`
	Author    string    `json:"author"`
	Votes     int       `json:"votes"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Vote struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	QuoteID   uint      `json:"quote_id"`
	CreatedAt time.Time `json:"created_at"`
}
