package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const (
	LoggedInDevicesOSAndroid = 1
	LoggedInDevicesOSiOS     = 2
)

type LoggedInDevice struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	UserID    int
	User      User `gorm:"constraint:OnDelete:CASCADE;not null"`
	OS        int
	Token     string `gorm:"not null;index"`
	FCMToken  string
	Version   float64
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

func (d *LoggedInDevice) IsExpired() bool {
	return d.ExpiresAt.Before(time.Now())
}

func (d *LoggedInDevice) SetToken() string {
	if d.Token != "" {
		return d.Token
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	d.Token = newUUID.String()

	return d.Token
}
