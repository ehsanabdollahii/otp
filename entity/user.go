package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Name      string
	IsActive  bool
	Phone     string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
