//go:build wireinject
// +build wireinject

package gympack

import (
	"github.com/google/wire"
	"gympack/pkg/domain/usecase/pack/create"
	provider "gympack/pkg/infrastructure/database/mongodb"
	"gympack/pkg/infrastructure/database/mongodb/repository/pack"
	"gympack/pkg/presentation"
	middlewares "gympack/pkg/presentation/middleware"
	"gympack/pkg/presentation/route"
	"gympack/pkg/shared"
)

var superSet = wire.NewSet(
	presentation.NewLoader,
	shared.NewConfig,
	middlewares.NewDtoValidationMiddleware,
	create.NewCreatePackUseCase,
	route.NewPackRoute,
	shared.NewLogger,
	provider.NewMongoDbProvider,
	pack.NewPackRepository,
)

func InitializeLoader() *presentation.Loader {
	wire.Build(superSet)
	return &presentation.Loader{}
}
