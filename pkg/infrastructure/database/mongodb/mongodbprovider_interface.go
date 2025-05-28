package provider

import "go.mongodb.org/mongo-driver/mongo"

type MongoDbProviderInterface interface {
	GetMongoDbClient() (*mongo.Client, error)
}
