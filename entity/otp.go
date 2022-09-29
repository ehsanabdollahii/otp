package entity

import (
	"database/sql/driver"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type OTPType string

const (
	OTPTypeLogin OTPType = "login"
)

func (e *OTPType) Scan(value interface{}) error {
	*e = OTPType(value.([]byte))
	return nil
}

func (e OTPType) Value() (driver.Value, error) {
	return string(e), nil
}

type OneTimePassword struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	UserID    *int
	Source    string  `gorm:"not null;index"` // for now phone number but might add emails later
	User      *User   `gorm:"constraint:OnDelete:SET NULL"`
	Type      OTPType `sql:"type:ENUM('login')"`
	Code      string  `gorm:"not null;index"`
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (otp *OneTimePassword) SetCode(length int, digitOnly bool) string {
	if otp.Code != "" {
		return otp.Code
	}

	rand.Seed(time.Now().UnixNano())

	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	if digitOnly {
		letterRunes = []rune("1234567890")
	}

	b := make([]rune, length)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	otp.Code = string(b)

	return otp.Code
}

func (otp *OneTimePassword) IsExpired() bool {
	return time.Now().After(otp.ExpiresAt)
}
