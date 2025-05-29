package get

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/repository/filter"
	"gympack/pkg/infrastructure/database/mongodb/repository/pack"
	"gympack/pkg/shared"
	"gympack/pkg/shared/constant"
)

func NewGetPackUseCase(
	log shared.LoggerInterface,
	packRepository pack.PackRepositoryInterface,
) GetPackUseCaseInterface {
	return &GetPackUseCase{
		log:            log,
		packRepository: packRepository,
	}
}

type GetPackUseCase struct {
	log            shared.LoggerInterface
	packRepository pack.PackRepositoryInterface
}

func (g *GetPackUseCase) Execute(ctx context.Context, id string) (*[]model.PackModel, error) {
	ctx, span := otel.Tracer("CreatePackUseCase").Start(ctx, "CreatePackUseCase.Execute")
	defer span.End()

	var pFilter filter.BaseFilter

	pFilter.Sort = map[string]constant.RepositoryOrder{"minSize": constant.RepositoryOrderASC}

	if id != "" {
		pFilter = filter.BaseFilter{
			Equals: map[string]interface{}{"_id": id},
			Sort:   map[string]constant.RepositoryOrder{"minSize": -1},
		}
	}

	fRes, err := g.packRepository.FindByFilter(ctx, pFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("FindByFilter: %w", err)
	}

	return fRes, nil
}
