package models

import (
	"time"
)

type AlertType struct {
	ID               string     `gorm:"type:varchar(50);primaryKey"`
	Name             string     `gorm:"type:varchar(255);not null"`
	Description      string     `gorm:"type:text;not null"`
	ConfigSchema     string     `gorm:"type:text;not null"`
	IsCustom         bool       `gorm:"default:false"`
	RequiredPlan     string     `gorm:"type:varchar(50);not null;default:'Free'"`
	CreatedBy        *string    `gorm:"type:varchar(36);null"`
	CreatedAt        time.Time  `gorm:"autoCreateTime"`

	Alerts []Alert `gorm:"foreignKey:AlertTypeID"`
}

func (AlertType) TableName() string {
	return "alert_types"
}