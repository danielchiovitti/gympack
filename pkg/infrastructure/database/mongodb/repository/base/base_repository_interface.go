package base

import (
	"context"
	"gympack/pkg/infrastructure/database/mongodb/repository/filter"
)

type BaseRepositoryInterface[T, U any] interface {
	InsertOne(ctx context.Context, pModel T) (*T, error)
	DeleteOneById(ctx context.Context, id string) (*int64, error)
	DeleteOneByFilter(ctx context.Context, pFilter filter.BaseFilter) (*int64, error)
	FindOneById(ctx context.Context, id string, project []string) (*T, error)
	FindOneByFilter(ctx context.Context, pFilter filter.BaseFilter, project []string) (*T, error)
	FindByFilter(ctx context.Context, pFilter filter.BaseFilter, project []string) (*[]T, error)
	UpdateOneById(ctx context.Context, id string, model T) (*int64, error)
	UpdateOneByFilter(ctx context.Context, pFilter filter.BaseFilter, model T) (*int64, error)
	IsValidMandatoryFilters(ctx context.Context, pFilter filter.BaseFilter, project []string) (bool, error)
}
