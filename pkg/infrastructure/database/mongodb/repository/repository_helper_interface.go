package repository

import "go.mongodb.org/mongo-driver/bson"

type RepositoryHelperInterface[T, U any] interface {
	ProjectionIsValid(fields []string) bool
	GetProjection(fields []string) bson.M
}
