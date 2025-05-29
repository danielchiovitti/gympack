package mocks

import (
	"go.mongodb.org/mongo-driver/mongo"
)

func NewMongodbProviderMock(
	client *mongo.Client,
) *MongodbProviderMock {
	return &MongodbProviderMock{
		client: client,
	}
}

type MongodbProviderMock struct {
	client *mongo.Client
}

func (m *MongodbProviderMock) GetMongoDbClient() (*mongo.Client, error) {
	return m.client, nil
}
