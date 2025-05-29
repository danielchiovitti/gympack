package get

import (
	"context"
	"gympack/pkg/domain/model"
)

type GetPackUseCaseInterface interface {
	Execute(ctx context.Context, id string) (*[]model.PackModel, error)
}
