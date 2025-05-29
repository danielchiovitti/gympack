package create

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/repository/filter"
	"gympack/pkg/infrastructure/database/mongodb/repository/pack"
	"gympack/pkg/shared"
	"gympack/pkg/shared/helpers"
)

func NewCreatePackUseCase(
	log shared.LoggerInterface,
	packRepository pack.PackRepositoryInterface,
) CreatePackUseCaseInterface {
	return &CreatePackUseCase{
		log:            log,
		packRepository: packRepository,
	}
}

type CreatePackUseCase struct {
	log            shared.LoggerInterface
	packRepository pack.PackRepositoryInterface
}

func (c *CreatePackUseCase) Execute(ctx context.Context, model model.PackModel) (*model.PackModel, error) {
	ctx, span := otel.Tracer("CreatePackUseCase").Start(ctx, "CreatePackUseCase.Execute")
	defer span.End()

	if model.Id == "" {
		model.Id = primitive.NewObjectID().Hex()
	}

	if model.MaxSize <= model.MinSize {
		return nil, fmt.Errorf("max size must be greater than min size")
	}

	pFilter := filter.BaseFilter{
		OrFilters: []filter.BaseFilter{
			{Range: map[string]filter.RangeFilter{
				"minSize": {Min: helpers.ToInterfacePtr(model.MinSize), Max: helpers.ToInterfacePtr(model.MaxSize)}},
			},
			{Range: map[string]filter.RangeFilter{
				"maxSize": {Min: helpers.ToInterfacePtr(model.MinSize), Max: helpers.ToInterfacePtr(model.MaxSize)}},
			},
		},
	}

	fRes, err := c.packRepository.FindByFilter(ctx, pFilter, []string{"Id"})
	if err != nil {
		return nil, fmt.Errorf("FindByFilter: %w", err)
	}

	if *fRes != nil && len(*fRes) > 0 {
		return nil, fmt.Errorf("Min or Max already registered")
	}

	c.log.Info(fmt.Sprintf("CreatePackUseCase execute model: %v", model))
	resPack, err := c.packRepository.InsertOne(ctx, model)
	if err != nil {
		c.log.Error(fmt.Sprintf("CreatePackUseCase execute packRepository.InsertOne error: %v", err))
		span.RecordError(err)
		return nil, err
	}

	return resPack, nil
}
