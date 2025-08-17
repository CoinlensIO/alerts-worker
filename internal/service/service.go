package service

import "alerts-worker/internal/repository"

type Service struct {
	userRepo *repository.Repository
}

func New(userRepo *repository.Repository) *Service {
	return &Service{userRepo: userRepo}
}
