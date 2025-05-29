package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	constant2 "gympack/pkg/shared/constant"
	"reflect"
)

type BaseFilter struct {
	AndFilters  []BaseFilter
	OrFilters   []BaseFilter
	Exists      map[string]bool
	Equals      map[string]interface{}
	NotEquals   map[string]interface{}
	In          map[string][]interface{}
	NotIn       map[string][]interface{}
	Range       map[string]RangeFilter
	GreaterThan map[string]interface{}
	LessThan    map[string]interface{}
	Skip        *int64
	Limit       *int64
	Sort        map[string]constant2.RepositoryOrder
}

func TransformKeyValue(key string, value interface{}) interface{} {
	if key == "_id" {
		if reflect.TypeOf(value).Kind() == reflect.String {
			v, _ := primitive.ObjectIDFromHex(value.(string))
			return v
		}
	}
	return value
}

func BuildBSONFilter(filter BaseFilter) bson.M {
	bsonFilter := bson.M{}

	for k, v := range filter.Equals {
		bsonFilter[k] = TransformKeyValue(k, v)
	}

	for k, v := range filter.NotEquals {
		bsonFilter[k] = bson.M{"$ne": TransformKeyValue(k, v)}
	}

	for k, v := range filter.In {
		bsonFilter[k] = bson.M{"$in": TransformKeyValue(k, v)}
	}

	for k, v := range filter.NotIn {
		bsonFilter[k] = bson.M{"$nin": TransformKeyValue(k, v)}
	}

	for k, rangeFilter := range filter.Range {
		condition := bson.M{}

		if rangeFilter.Min != nil {
			condition["$gte"] = rangeFilter.Min
		}

		if rangeFilter.Max != nil {
			condition["$lte"] = rangeFilter.Max
		}

		bsonFilter[k] = condition
	}

	for k, v := range filter.GreaterThan {
		bsonFilter[k] = bson.M{"$gt": TransformKeyValue(k, v)}
	}

	for k, v := range filter.LessThan {
		bsonFilter[k] = bson.M{"$lt": TransformKeyValue(k, v)}
	}

	for k, v := range filter.Exists {
		bsonFilter[k] = bson.M{"$exists": TransformKeyValue(k, v)}
	}

	if len(filter.AndFilters) > 0 {
		var andConditions []bson.M
		for _, subFilter := range filter.AndFilters {
			//andConditions = append(andConditions, bson.M{"$and": subFilter})
			andConditions = append(andConditions, BuildBSONFilter(subFilter))
		}
		bsonFilter["$and"] = andConditions
	}

	if len(filter.OrFilters) > 0 {
		var orConditions []bson.M
		for _, subFilter := range filter.OrFilters {
			//orConditions = append(orConditions, bson.M{"$or": subFilter})
			orConditions = append(orConditions, BuildBSONFilter(subFilter))
		}
		bsonFilter["$or"] = orConditions
	}

	if len(filter.Sort) > 0 {

	}

	return bsonFilter
}

func BuildSort(sort map[string]int) bson.D {
	var sortQuery bson.D
	for k, v := range sort {
		if v != 1 && v != -1 {
			continue
		}
		sortQuery = append(sortQuery, bson.E{Key: k, Value: v})
	}
	return sortQuery
}
