package repository

import (
	"context"
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gympack/pkg/domain/model"
	"gympack/pkg/infrastructure/database/mongodb/repository/pack"
	"gympack/test/mocks"
	"testing"
	"time"
)

func TestPackRepositoryInsertOne(t *testing.T) {
	sWant := []model.PackModel{
		{
			Id:          primitive.NewObjectID().Hex(),
			Name:        faker.Word(),
			Description: faker.Word(),
			MaxSize:     5,
			BaseModel: model.BaseModel{
				CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
				CreatedById: primitive.NewObjectID().Hex(),
				UpdatedAt:   time.Now().UTC().Truncate(time.Millisecond),
				UpdatedById: primitive.NewObjectID().Hex(),
			},
		},
		{
			Id:          primitive.NewObjectID().Hex(),
			Name:        faker.Word(),
			Description: faker.Word(),
			MaxSize:     10,
			BaseModel: model.BaseModel{
				CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
				CreatedById: primitive.NewObjectID().Hex(),
				UpdatedAt:   time.Now().UTC().Truncate(time.Millisecond),
				UpdatedById: primitive.NewObjectID().Hex(),
			},
		},
		{
			Id:          primitive.NewObjectID().Hex(),
			Name:        faker.Word(),
			Description: faker.Word(),
			MaxSize:     20,
			BaseModel: model.BaseModel{
				CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
				CreatedById: primitive.NewObjectID().Hex(),
				UpdatedAt:   time.Now().UTC().Truncate(time.Millisecond),
				UpdatedById: primitive.NewObjectID().Hex(),
			},
		},
	}

	ctx := context.Background()
	uri, err := GetMongoContainer().ConnectionString(ctx)
	if err != nil {
		t.Fatal("Failed to setup mongo container")
	}

	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	mongodbProviderMock := mocks.NewMongodbProviderMock(mongoClient)

	repository := pack.NewPackRepository(mongodbProviderMock)

	for _, v := range sWant {
		res, err := repository.InsertOne(ctx, v)
		if err != nil {
			t.Errorf("Failed to insert pack: %v", err)
		}

		if res.Id != v.Id {
			t.Errorf("Failed to insert pack Id. Expected %v, got %v", v.Id, res.Id)
		}

		if res.Name != v.Name {
			t.Errorf("Failed to insert pack Name. Expected %v, got %v", v.Name, res.Name)
		}

		if res.Description != v.Description {
			t.Errorf("Failed to insert pack Description. Expected %v, got %v", v.Description, res.Description)
		}

		if res.MaxSize != v.MaxSize {
			t.Errorf("Failed to insert pack MaxSize. Expected %v, got %v", v.MaxSize, res.MaxSize)
		}

		if res.BaseModel.CreatedAt != v.BaseModel.CreatedAt {
			t.Errorf("Failed to insert pack CreatedAt. Expected %v, got %v", v.MaxSize, res.MaxSize)
		}

		if res.BaseModel.CreatedById != v.BaseModel.CreatedById {
			t.Errorf("Failed to insert pack CreatedById. Expected %v, got %v", v.MaxSize, res.MaxSize)
		}

		if res.BaseModel.UpdatedAt != v.BaseModel.UpdatedAt {
			t.Errorf("Failed to insert pack UpdatedAt. Expected %v, got %v", v.MaxSize, res.MaxSize)
		}

		if res.BaseModel.UpdatedById != v.BaseModel.UpdatedById {
			t.Errorf("Failed to insert pack UpdatedById. Expected %v, got %v", v.MaxSize, res.MaxSize)
		}
	}

}
