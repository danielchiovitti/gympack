package delete

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/repository/pack"
	"gympack/pkg/shared"
)

func NewDeletePackUseCase(
	log shared.LoggerInterface,
	packRepository pack.PackRepositoryInterface,
) DeletePackUseCaseInterface {
	return &DeletePackUseCase{
		log:            log,
		packRepository: packRepository,
	}
}

type DeletePackUseCase struct {
	log            shared.LoggerInterface
	packRepository pack.PackRepositoryInterface
}

func (d *DeletePackUseCase) Execute(ctx context.Context, id string) (*model.PackModel, error) {
	ctx, span := otel.Tracer("DeletePackUseCase").Start(ctx, "DeletePackUseCase.Execute")
	defer span.End()

	pModel, err := d.packRepository.FindOneById(ctx, id, nil)
	if err != nil {
		d.log.Error(fmt.Sprintf("DeletePackUseCase.Execute error trying to find register %s", id))
		span.RecordError(err)
		return nil, err
	}

	dQty, err := d.packRepository.DeleteOneById(ctx, id)
	if dQty != nil && *dQty != 1 {
		d.log.Error(fmt.Sprintf("DeletePackUseCase.Execute error trying to delete workspace %s", id))
	}

	if err != nil {
		d.log.Error(fmt.Sprintf("DeletePackUseCase.Execute error trying to delete pack %s", id))
		span.RecordError(err)
		return nil, err
	}

	return pModel, nil
}
