package calc

import (
	"context"
	"gympack/pkg/domain/model"
)

type CalcPackUseCaseInterface interface {
	Execute(ctx context.Context, qty uint) (*[]model.PackModel, error)
}
