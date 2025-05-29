package create

import (
	"context"
	"gympack/pkg/domain/model"
)

type CreatePackUseCaseInterface interface {
	Execute(ctx context.Context, pModel model.PackModel) (*model.PackModel, error)
}
