package models

import (
	"time"
)

type Users struct {
	ID                string              `gorm:"type:varchar(36);primaryKey"`
	FirebaseUID       string              `gorm:"type:varchar(36);uniqueIndex"`
	Email             string              `gorm:"uniqueIndex;not null"`
	PasswordHash      string              `gorm:"type:varchar(255);not null"`
	DisplayName       string              `gorm:"type:varchar(255)"`
	PhotoURL          string              `gorm:"type:text"`
	Provider          string              `gorm:"type:varchar(50);not null"`
	IsVerified        bool                `gorm:"default:false"`
	CreatedAt         time.Time           `gorm:"autoCreateTime"`
	UpdatedAt         time.Time           `gorm:"autoUpdateTime"`
	VerificationCodes []VerificationCodes `gorm:"foreignKey:UserID"`
	Subscriptions     []UserSubscription  `gorm:"foreignKey:UserID"`
}
