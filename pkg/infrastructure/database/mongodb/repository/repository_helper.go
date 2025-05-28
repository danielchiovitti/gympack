package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

func NewRepositoryHelper[T, U any]() *RepositoryHelper[T, U] {
	return &RepositoryHelper[T, U]{
		uType: reflect.TypeOf((*U)(nil)).Elem(),
	}
}

type RepositoryHelper[T, U any] struct {
	uType reflect.Type
}

func (r *RepositoryHelper[T, U]) ProjectionIsValid(fields []string) bool {
	for _, field := range fields {
		exists := false

		for i := 0; i < r.uType.NumField(); i++ {
			currField := r.uType.Field(i)
			bsonTag := currField.Tag.Get("bson")

			if bsonTag == "" {
				continue
			}

			if field == currField.Name {
				exists = true
			}

		}

		if !exists {
			return false
		}
	}

	return true
}

func (r *RepositoryHelper[T, U]) GetProjection(fields []string) bson.M {
	projection := bson.M{}

	if len(fields) == 0 {
		return projection
	}

	projection["_id"] = 1

	for _, field := range fields {
		for i := 0; i < r.uType.NumField(); i++ {
			currField := r.uType.Field(i)
			bsonTag := currField.Tag.Get("bson")
			tag := strings.Split(bsonTag, ",")[0]

			if bsonTag == "" {
				continue
			}

			if field == bsonTag {
				projection[tag] = 1
			}
		}
	}

	return projection
}
