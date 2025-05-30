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

func (c *CreatePackUseCase) Execute(ctx context.Context, pModel model.PackModel) (*model.PackModel, error) {
	ctx, span := otel.Tracer("CreatePackUseCase").Start(ctx, "CreatePackUseCase.Execute")
	defer span.End()

	if pModel.Id == "" {
		pModel.Id = primitive.NewObjectID().Hex()
	}

	pFilter := filter.BaseFilter{
		Equals: map[string]interface{}{"maxSize": pModel.MaxSize},
	}

	fRes, err := c.packRepository.FindByFilter(ctx, pFilter, []string{"Id"})
	if err != nil {
		return nil, fmt.Errorf("FindByFilter: %w", err)
	}

	if *fRes != nil && len(*fRes) > 0 {
		return nil, fmt.Errorf("max already registered")
	}

	c.log.Info(fmt.Sprintf("CreatePackUseCase execute pModel: %v", pModel))
	resPack, err := c.packRepository.InsertOne(ctx, pModel)
	if err != nil {
		c.log.Error(fmt.Sprintf("CreatePackUseCase execute packRepository.InsertOne error: %v", err))
		span.RecordError(err)
		return nil, err
	}

	return resPack, nil
}
