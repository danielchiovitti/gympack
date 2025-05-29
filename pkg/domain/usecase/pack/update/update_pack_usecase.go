package update

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/repository/pack"
	"gympack/pkg/shared"
)

func NewUpdatePackUseCase(
	log shared.LoggerInterface,
	packRepository pack.PackRepositoryInterface,
) UpdatePackUseCaseInterface {
	return &UpdatePackUseCase{
		log:            log,
		packRepository: packRepository,
	}
}

type UpdatePackUseCase struct {
	log            shared.LoggerInterface
	packRepository pack.PackRepositoryInterface
}

func (u *UpdatePackUseCase) Execute(ctx context.Context, id string, model model.PackModel) (*model.PackModel, error) {
	ctx, span := otel.Tracer("UpdatePackUseCase").Start(ctx, "UpdatePackUseCase.Execute")
	defer span.End()

	u.log.Info(fmt.Sprintf("UpdatePackUseCase execute model: %v", model))
	model.Id = id
	_, err := u.packRepository.UpdateOneById(ctx, id, model)
	if err != nil {
		u.log.Error(fmt.Sprintf("UpdatePackUseCase execute error: %v", err))
		span.RecordError(err)
		return nil, err
	}

	resp, err := u.packRepository.FindOneById(ctx, id, nil)
	if err != nil {
		u.log.Error(fmt.Sprintf("UpdatePackUseCase execute packRepository.FindOneById error: %v", err))
		span.RecordError(err)
		return nil, err
	}

	return resp, nil
}
