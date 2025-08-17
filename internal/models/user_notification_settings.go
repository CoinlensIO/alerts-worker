package models

import (
	"time"
)

type UserNotificationSettings struct {
	ID              string    `gorm:"type:varchar(36);primaryKey"`
	UserID          string    `gorm:"type:varchar(36);not null;index;unique"`
	EmailEnabled    bool      `gorm:"default:true"`
	TelegramEnabled bool      `gorm:"default:false"`
	PushEnabled     bool      `gorm:"default:false"`
	TelegramHandle  *string   `gorm:"type:varchar(255);null"`
	DeviceToken     *string   `gorm:"type:varchar(500);null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	User            Users     `gorm:"foreignKey:UserID"`
}

type NotificationChannel string

const (
	NotificationChannelEmail    NotificationChannel = "email"
	NotificationChannelTelegram NotificationChannel = "telegram"
	NotificationChannelPush     NotificationChannel = "push"
)

type AlertNotificationTarget struct {
	ID        string              `gorm:"type:varchar(36);primaryKey"`
	AlertID   string              `gorm:"type:varchar(36);not null;index"`
	Channel   NotificationChannel `gorm:"type:varchar(20);not null"`
	IsEnabled bool                `gorm:"default:true"`
	CreatedAt time.Time           `gorm:"autoCreateTime"`
	Alert     Alert               `gorm:"foreignKey:AlertID"`
}