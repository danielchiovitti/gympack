package calc

import (
	"context"
	"fmt"
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/repository/filter"
	"gympack/pkg/infrastructure/database/mongodb/repository/pack"
	"gympack/pkg/shared"
	"gympack/pkg/shared/constant"
)

func NewCalcPackUseCase(
	log shared.LoggerInterface,
	packRepository pack.PackRepositoryInterface,
) CalcPackUseCaseInterface {
	return &CalcPackUseCase{
		log:            log,
		packRepository: packRepository,
	}
}

type CalcPackUseCase struct {
	log            shared.LoggerInterface
	packRepository pack.PackRepositoryInterface
}

func (c *CalcPackUseCase) Execute(ctx context.Context, qty uint) (*[]model.PackModel, error) {
	current := int(qty)
	var respPacks []model.PackModel

	pFilter := filter.BaseFilter{
		Sort: map[string]constant.RepositoryOrder{"maxSize": constant.RepositoryOrderDESC},
	}

	packs, err := c.packRepository.FindByFilter(ctx, pFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("FindByFilter: %w", err)
	}

	for current > 0 {
		var currentPack *model.PackModel
		for i := range *packs {
			v := (*packs)[i]

			if current >= v.MaxSize {
				currentPack = &v
				break
			}
		}

		if currentPack == nil {
			currentPack = &(*packs)[len(*packs)-1]
		}

		current -= currentPack.MaxSize
		respPacks = append(respPacks, *currentPack)
	}

	return &respPacks, nil
}
