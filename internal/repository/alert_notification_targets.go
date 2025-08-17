package repository

import (
	"alerts-worker/internal/models"
	"context"
	"gorm.io/gorm"
)

type AlertNotificationTargetRepository interface {
	GetAlertNotificationTargets(ctx context.Context, alertID string) ([]models.AlertNotificationTarget, error)
}

type alertNotificationTargetRepository struct {
	db *gorm.DB
}

func NewAlertNotificationTargetRepository(db *gorm.DB) AlertNotificationTargetRepository {
	return &alertNotificationTargetRepository{db: db}
}

func (r *alertNotificationTargetRepository) GetAlertNotificationTargets(ctx context.Context, alertID string) ([]models.AlertNotificationTarget, error) {
	var targets []models.AlertNotificationTarget
	err := r.db.WithContext(ctx).Where("alert_id = ?", alertID).Find(&targets).Error
	return targets, err
}
