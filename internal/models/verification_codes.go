package models

import "time"

type VerificationCodes struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;type:bigint"`
	UserID    string    `gorm:"type:varchar(36);not null"`
	Type      string    `gorm:"type:varchar(20);not null"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Code      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiryAt  time.Time `gorm:"not null"`
	IsUsed    bool      `gorm:"type:boolean"`
}
