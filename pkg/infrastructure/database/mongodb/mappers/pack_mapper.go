package mappers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/entity"
	"strings"
)

type PackMapper struct{}

func (p PackMapper) ToEntity(model model.PackModel) (*entity.PackEntity, error) {
	id, err := primitive.ObjectIDFromHex(model.Id)
	if err != nil {
		return nil, err
	}

	var createdById primitive.ObjectID
	if strings.TrimSpace(model.CreatedById) != "" {
		createdById, err = primitive.ObjectIDFromHex(model.CreatedById)
		if err != nil {
			return nil, err
		}
	}

	var updatedById primitive.ObjectID
	if strings.TrimSpace(model.UpdatedById) != "" {
		updatedById, err = primitive.ObjectIDFromHex(model.UpdatedById)
		if err != nil {
			return nil, err
		}
	}

	var deletedById primitive.ObjectID
	if strings.TrimSpace(model.DeletedById) != "" {
		deletedById, err = primitive.ObjectIDFromHex(model.DeletedById)
		if err != nil {
			return nil, err
		}
	}

	return &entity.PackEntity{
		Id:          id,
		Name:        model.Name,
		Description: model.Description,
		MaxSize:     model.MaxSize,
		BaseEntity: entity.BaseEntity{
			CreatedAt:   model.CreatedAt,
			CreatedById: createdById,
			UpdatedAt:   model.UpdatedAt,
			UpdatedById: updatedById,
			Deleted:     model.Deleted,
			DeletedAt:   model.DeletedAt,
			DeletedById: deletedById,
		},
	}, nil
}

func (p PackMapper) ToModel(entity entity.PackEntity) (*model.PackModel, error) {
	return &model.PackModel{
		Id:          entity.Id.Hex(),
		Name:        entity.Name,
		Description: entity.Description,
		MaxSize:     entity.MaxSize,
		BaseModel: model.BaseModel{
			CreatedAt:   entity.CreatedAt,
			CreatedById: entity.CreatedById.Hex(),
			UpdatedAt:   entity.UpdatedAt,
			UpdatedById: entity.UpdatedById.Hex(),
			Deleted:     entity.Deleted,
			DeletedAt:   entity.DeletedAt,
			DeletedById: entity.DeletedById.Hex(),
		},
	}, nil
}
