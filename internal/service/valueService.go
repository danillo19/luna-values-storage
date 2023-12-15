package service

import (
	"context"
	"luna-values-storage/internal/domain"
)

type ValueRepository interface {
	Get(ctx context.Context, id string) (*domain.Value, error)
	Delete(ctx context.Context, id string) error
	Set(ctx context.Context, value *domain.Value) (*domain.Value, error)
}

type ValueService struct {
	repo ValueRepository
}

func NewValueService(repo ValueRepository) *ValueService {
	return &ValueService{repo: repo}
}

func (s ValueService) GetValue(ctx context.Context, id string) (*domain.Value, error) {
	return s.repo.Get(ctx, id)
}

func (s ValueService) SetValue(ctx context.Context, value *domain.Value) (*domain.Value, error) {
	return s.repo.Set(ctx, value)
}

func (s ValueService) DeleteValue(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
