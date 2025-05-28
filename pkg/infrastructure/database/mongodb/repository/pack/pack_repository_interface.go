package pack

import (
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/entity"
	"gympack/pkg/infrastructure/database/mongodb/repository/base"
)

type PackRepositoryInterface interface {
	base.BaseRepositoryInterface[model.PackModel, entity.PackEntity]
}
