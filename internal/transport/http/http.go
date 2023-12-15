package http

import (
	"context"
	"go.uber.org/zap"
	"luna-values-storage/internal/common"
	"luna-values-storage/internal/domain"
)

type VariableService interface {
	GetVariable(ctx context.Context, id string) (*domain.Variable, error)
	SetVariable(ctx context.Context, variable *domain.Variable) (*domain.Variable, error)
	DeleteVariable(ctx context.Context, id string) error
}

type ValueService interface {
	GetValue(ctx context.Context, id string) (*domain.Value, error)
	SetValue(ctx context.Context, value *domain.Value) (*domain.Value, error)
	DeleteValue(ctx context.Context, id string) error
}

type Resolver struct {
	logger          *zap.SugaredLogger
	valueService    ValueService
	variableService VariableService
	s3Client        *common.S3Client
}

func NewResolver(vars VariableService, values ValueService, s3Client *common.S3Client, logger *zap.SugaredLogger) Resolver {
	return Resolver{
		valueService:    values,
		variableService: vars,
		s3Client:        s3Client,
		logger:          logger,
	}
}
