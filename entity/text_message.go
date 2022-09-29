package entity

import (
	"gorm.io/gorm"
	"time"
)

const (
	TextMessageTypeOTPLogin   = 1
	TextMessageWelcomeMessage = 2
)

type TextMessage struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Name      string
	IsActive  bool
	Type      int    `gorm:"index:typePhone,priority:1"`
	Phone     string `gorm:"index:typePhone,priority:2"`
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
