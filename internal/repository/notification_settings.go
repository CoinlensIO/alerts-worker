package repository

import (
	"alerts-worker/internal/models"
	"context"

	"gorm.io/gorm"
)

type NotificationSettingsRepository interface {
	GetUserNotificationSettings(ctx context.Context, userID string) (*models.UserNotificationSettings, error)
	CreateOrUpdateNotificationSettings(ctx context.Context, settings *models.UserNotificationSettings) error
	DeleteNotificationSettings(ctx context.Context, userID string) error
}

type notificationSettingsRepository struct {
	db *gorm.DB
}

func NewNotificationSettingsRepository(db *gorm.DB) NotificationSettingsRepository {
	return &notificationSettingsRepository{db: db}
}

func (r *notificationSettingsRepository) GetUserNotificationSettings(ctx context.Context, userID string) (*models.UserNotificationSettings, error) {
	var settings models.UserNotificationSettings
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *notificationSettingsRepository) CreateOrUpdateNotificationSettings(ctx context.Context, settings *models.UserNotificationSettings) error {
	return r.db.WithContext(ctx).Save(settings).Error
}

func (r *notificationSettingsRepository) DeleteNotificationSettings(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.UserNotificationSettings{}).Error
}
