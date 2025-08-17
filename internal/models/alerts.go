package models

import (
	"time"
)

type Alert struct {
	ID            string     `gorm:"type:varchar(36);primaryKey"`
	UserID        string     `gorm:"type:varchar(36);not null;index"`
	AlertTypeID   string     `gorm:"type:varchar(50);not null"`
	Name          string     `gorm:"type:varchar(255);not null"`
	Description   string     `gorm:"type:text"`
	Conditions    string     `gorm:"type:text;not null"`
	IsActive      bool       `gorm:"default:true"`
	LastTriggered *time.Time `gorm:"null"`
	TriggerCount  int        `gorm:"default:0"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`

	User               Users                     `gorm:"foreignKey:UserID"`
	AlertType          AlertType                 `gorm:"foreignKey:AlertTypeID"`
	NotificationTargets []AlertNotificationTarget `gorm:"foreignKey:AlertID"`
}

func (Alert) TableName() string {
	return "alerts"
}