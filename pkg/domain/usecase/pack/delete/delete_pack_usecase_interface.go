package delete

import (
	"context"
	"gympack/pkg/domain/model"
)

type DeletePackUseCaseInterface interface {
	Execute(ctx context.Context, id string) (*model.PackModel, error)
}
