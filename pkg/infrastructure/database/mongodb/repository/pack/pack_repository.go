package pack

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gympack/pkg/domain/model"
	provider "gympack/pkg/infrastructure/database/mongodb"
	"gympack/pkg/infrastructure/database/mongodb/entity"
	"gympack/pkg/infrastructure/database/mongodb/mappers"
	"gympack/pkg/infrastructure/database/mongodb/repository"
	"gympack/pkg/infrastructure/database/mongodb/repository/base"
	"gympack/pkg/shared/constant"
)

func NewPackRepository(
	mongodbProvider provider.MongoDbProviderInterface,
) PackRepositoryInterface {
	client, err := mongodbProvider.GetMongoDbClient()
	if err != nil {
		return nil
	}

	return &PackRepository{
		client: client,
		BaseRepository: base.NewBaseRepository[model.PackModel, entity.PackEntity](
			client,
			constant.MGDB_CORE,
			constant.PACK,
			mappers.PackMapper{},
			repository.NewRepositoryHelper[model.PackModel, entity.PackEntity](),
		),
	}
}

type PackRepository struct {
	client *mongo.Client
	*base.BaseRepository[model.PackModel, entity.PackEntity]
}
