package models

import "time"

type SubscriptionPlan struct {
	ID              string    `gorm:"type:varchar(36);primaryKey"`
	Name            string    `gorm:"type:varchar(100);not null"`
	Price           int       `gorm:"not null"`
	Currency        string    `gorm:"type:varchar(3);default:'USD'"`
	BillingInterval string    `gorm:"type:varchar(20);not null"`
	TrialDays       int       `gorm:"default:0"`
	Features        string    `gorm:"type:text"`
	Limits          string    `gorm:"type:text"`
	IsActive        bool      `gorm:"default:true"`
	ExternalPlanIDs string    `gorm:"type:text"`
	Description     string    `gorm:"type:text"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
