package models

import "time"

type UserSubscription struct {
	ID                     string    `gorm:"type:varchar(36);primaryKey"`
	UserID                 string    `gorm:"type:varchar(36);not null;index"`
	PlanID                 string    `gorm:"type:varchar(36);not null"`
	Platform               string    `gorm:"type:varchar(20);not null"`
	Status                 string    `gorm:"type:varchar(20);not null"`
	ExternalSubscriptionID string    `gorm:"type:varchar(255)"`
	PaymentProvider        string    `gorm:"type:varchar(50)"`
	StartDate              time.Time `gorm:"not null"`
	EndDate                *time.Time
	NextBillingDate        *time.Time
	TrialEndDate           *time.Time
	CancelledAt            *time.Time
	CurrentPeriodStart     *time.Time
	CurrentPeriodEnd       *time.Time
	CancelAtPeriodEnd      bool      `gorm:"default:false"`
	Metadata               string    `gorm:"type:text"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`

	User Users            `gorm:"foreignKey:UserID"`
	Plan SubscriptionPlan `gorm:"foreignKey:PlanID"`
}
