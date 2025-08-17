package repository

import "gorm.io/gorm"

type Repository struct {
	db                       *gorm.DB
	NotificationSettings     NotificationSettingsRepository
	AlertNotificationTargets AlertNotificationTargetRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db:                       db,
		NotificationSettings:     NewNotificationSettingsRepository(db),
		AlertNotificationTargets: NewAlertNotificationTargetRepository(db),
	}
}
