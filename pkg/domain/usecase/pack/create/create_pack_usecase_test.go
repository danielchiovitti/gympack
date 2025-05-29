package create

import (
	"context"
	"github.com/go-faker/faker/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
	"gympack/pkg/domain/model"
	"gympack/test/mocks"
	"testing"
	"time"
)

func TestCreatePackUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	loggerMock := mocks.NewMockLoggerInterface(ctrl)
	loggerMock.EXPECT().Info(gomock.Any()).AnyTimes()

	var sliceModelReturn []model.PackModel

	mockModelReturn := model.PackModel{
		Id:          primitive.NewObjectID().Hex(),
		Name:        faker.Word(),
		Description: faker.Word(),
		MinSize:     1,
		MaxSize:     10,
		BaseModel: model.BaseModel{
			CreatedAt:   time.Now().UTC().Truncate(time.Millisecond),
			CreatedById: primitive.NewObjectID().Hex(),
			UpdatedAt:   time.Now().UTC().Truncate(time.Millisecond),
			UpdatedById: primitive.NewObjectID().Hex(),
		},
	}

	packRepositoryMock := mocks.NewMockPackRepositoryInterface(ctrl)
	packRepositoryMock.
		EXPECT().
		FindByFilter(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		Return(&sliceModelReturn, nil)
	packRepositoryMock.
		EXPECT().
		InsertOne(gomock.Any(), gomock.Any()).AnyTimes().
		Return(&mockModelReturn, nil)

	useCase := NewCreatePackUseCase(loggerMock, packRepositoryMock)

	_, err := useCase.Execute(ctx, mockModelReturn)
	if err != nil {
		t.Fatal(err)
	}

}
