package service

import (
	"context"
	"luna-values-storage/internal/domain"
)

type VariableRepository interface {
	Get(ctx context.Context, id string) (*domain.Variable, error)
	Delete(ctx context.Context, id string) error
	Set(ctx context.Context, variable *domain.Variable) (*domain.Variable, error)
}

type VariableService struct {
	repo VariableRepository
}

func NewVariableService(repo VariableRepository) *VariableService {
	return &VariableService{repo: repo}
}

func (s VariableService) GetVariable(ctx context.Context, id string) (*domain.Variable, error) {
	return s.repo.Get(ctx, id)
}

func (s VariableService) SetVariable(ctx context.Context, variable *domain.Variable) (*domain.Variable, error) {
	return s.repo.Set(ctx, variable)
}

func (s VariableService) DeleteVariable(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
