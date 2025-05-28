package create

import (
	"context"
	"gympack/pkg/domain/model"
)

type CreatePackUseCaseInterface interface {
	Execute(ctx context.Context, model model.PackModel) (*model.PackModel, error)
}
