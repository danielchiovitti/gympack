package filter

import "go.mongodb.org/mongo-driver/bson"

type BaseFilterInterface interface {
	BuildBSONFilter(filter BaseFilterInterface) bson.M
}
