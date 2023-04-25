package entities

import "time"

type PasswordResetTokens struct {
	Email       string    `gorm:"primaryKey;not null"`
	Token       string    `gorm:"primaryKey;not null"`
	TokenExpiry time.Time `gorm:"not null"`
	Code        string    `gorm:"not null"`
}
