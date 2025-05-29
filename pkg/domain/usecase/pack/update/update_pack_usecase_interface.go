package update

import (
	"context"
	"gympack/pkg/domain/model"
)

type UpdatePackUseCaseInterface interface {
	Execute(ctx context.Context, id string, model model.PackModel) (*model.PackModel, error)
}
