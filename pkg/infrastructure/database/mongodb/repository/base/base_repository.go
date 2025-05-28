package base

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"gympack/pkg/infrastructure/database/mongodb/mappers"
	"gympack/pkg/infrastructure/database/mongodb/repository"
	"gympack/pkg/infrastructure/database/mongodb/repository/filter"
	"gympack/pkg/shared/constant"
	"reflect"
	"strings"
)

type BaseRepository[T, U any] struct {
	Client           *mongo.Client
	Database         *mongo.Database
	Collection       constant.MongodbCollection
	Mapper           mappers.MapperInterface[T, U]
	RepositoryHelper repository.RepositoryHelperInterface[T, U]
	BaseFilter       filter.BaseFilterInterface
}

func NewBaseRepository[T, U any](
	client *mongo.Client,
	databaseName constant.MongodbInstance,
	collection constant.MongodbCollection,
	mapper mappers.MapperInterface[T, U],
	repositoryHelper repository.RepositoryHelperInterface[T, U],
) *BaseRepository[T, U] {
	return &BaseRepository[T, U]{
		Client:           client,
		Database:         client.Database(string(databaseName)),
		Collection:       collection,
		Mapper:           mapper,
		RepositoryHelper: repositoryHelper,
	}
}

func (b *BaseRepository[T, U]) InsertOne(ctx context.Context, model T) (*T, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.InsertOne.%s", uType.Name()))
	defer span.End()

	entity, err := b.Mapper.ToEntity(model)
	if err != nil {
		return nil, err
	}

	res, err := b.Database.Collection(string(b.Collection)).InsertOne(ctx, entity)
	if err != nil {
		return nil, err
	}

	insertedId := res.InsertedID.(primitive.ObjectID)
	objId := insertedId.Hex()

	respModel, err := b.FindOneById(ctx, objId, nil)
	if err != nil {
		return nil, err
	}

	return respModel, nil
}

func (b *BaseRepository[T, U]) DeleteOneById(ctx context.Context, id string) (*int64, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.DeleteOneById.%s", uType.Name()))
	defer span.End()

	innerId, _ := primitive.ObjectIDFromHex(id)
	pFilter := filter.BaseFilter{
		Equals: map[string]interface{}{
			"_id":     innerId,
			"deleted": false,
		},
	}
	builtFilter := filter.BuildBSONFilter(pFilter)

	update := bson.M{"$set": bson.M{"deleted": true}}
	res, err := b.Database.Collection(string(b.Collection)).UpdateOne(ctx, builtFilter, update)
	if err != nil {
		return nil, err
	}

	return &res.ModifiedCount, nil
}

func (b *BaseRepository[T, U]) DeleteOneByFilter(ctx context.Context, pFilter filter.BaseFilter) (*int64, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.DeleteOneByFilter.%s", uType.Name()))
	defer span.End()

	if fValid, fErr := b.IsValidMandatoryFilters(ctx, pFilter, nil); !fValid {
		return nil, fErr
	}

	update := bson.M{"$set": bson.M{"deleted": true}}

	pFilter.Equals["deleted"] = false
	builtFilter := filter.BuildBSONFilter(pFilter)

	res, err := b.Database.Collection(string(b.Collection)).UpdateOne(ctx, builtFilter, update)
	if err != nil {
		return nil, err
	}

	return &res.MatchedCount, nil
}

func (b *BaseRepository[T, U]) FindOneById(ctx context.Context, id string, project []string) (*T, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.FindOneById.%s", uType.Name()))
	defer span.End()

	var result U

	if len(project) > 0 {
		if !b.RepositoryHelper.ProjectionIsValid(project) {
			return nil, fmt.Errorf("invalid projection")
		}
	}

	projection := b.RepositoryHelper.GetProjection(project)

	pFilter := filter.BaseFilter{
		Equals: map[string]interface{}{
			"_id":     id,
			"deleted": false,
		},
	}
	builtFilter := filter.BuildBSONFilter(pFilter)

	singleResult := b.Database.Collection(string(b.Collection)).FindOne(ctx, builtFilter, options.FindOne().SetProjection(projection))
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}

	err := singleResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	mappedResult, err := b.Mapper.ToModel(result)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (b *BaseRepository[T, U]) FindOneByFilter(ctx context.Context, pFilter filter.BaseFilter, project []string) (*T, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.FindOneByFilter.%s", uType.Name()))
	defer span.End()

	if fValid, fErr := b.IsValidMandatoryFilters(ctx, pFilter, nil); !fValid {
		return nil, fErr
	}

	var result U

	if len(project) > 0 {
		if !b.RepositoryHelper.ProjectionIsValid(project) {
			return nil, fmt.Errorf("invalid projection")
		}
	}

	projection := b.RepositoryHelper.GetProjection(project)

	pFilter.Equals["deleted"] = false
	builtFilter := filter.BuildBSONFilter(pFilter)

	singleResult := b.Database.Collection(string(b.Collection)).FindOne(ctx, builtFilter, options.FindOne().SetProjection(projection))
	if singleResult.Err() != nil {
		return nil, singleResult.Err()
	}

	err := singleResult.Decode(&result)
	if err != nil {
		return nil, err
	}

	mappedResult, err := b.Mapper.ToModel(result)
	if err != nil {
		return nil, err
	}

	return mappedResult, nil
}

func (b *BaseRepository[T, U]) FindByFilter(ctx context.Context, pFilter filter.BaseFilter, project []string) (*[]T, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.FindByFilter.%s", uType.Name()))
	defer span.End()

	if fValid, fErr := b.IsValidMandatoryFilters(ctx, pFilter, nil); !fValid {
		return nil, fErr
	}

	var mappedResult []T

	if len(project) > 0 {
		if !b.RepositoryHelper.ProjectionIsValid(project) {
			return nil, fmt.Errorf("invalid projection")
		}
	}

	projection := b.RepositoryHelper.GetProjection(project)

	if pFilter.Equals == nil {
		pFilter.Equals = map[string]interface{}{}
	}

	pFilter.Equals["deleted"] = false
	builtFilter := filter.BuildBSONFilter(pFilter)

	opts := &options.FindOptions{}
	opts.SetProjection(projection)

	if pFilter.Skip != nil {
		opts.SetSkip(*pFilter.Skip)
	}

	if pFilter.Limit != nil {
		opts.SetLimit(*pFilter.Limit)
	}

	if pFilter.Sort != nil {
		var sortQuery bson.D
		for k, v := range pFilter.Sort {
			sortQuery = append(sortQuery, bson.E{Key: k, Value: int(v)})
		}
		opts.SetSort(sortQuery)
	}

	cursor, err := b.Database.Collection(string(b.Collection)).Find(ctx, builtFilter, opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var item U

		err = cursor.Decode(&item)
		if err == nil {

			innerMappedResult, err := b.Mapper.ToModel(item)
			if err != nil {
				return nil, err
			}

			mappedResult = append(mappedResult, *innerMappedResult)
		}
	}

	return &mappedResult, nil
}

func (b *BaseRepository[T, U]) UpdateOneById(ctx context.Context, id string, model T) (*int64, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.UpdateOneByFilter.%s", uType.Name()))
	defer span.End()

	update, err := b.Mapper.ToEntity(model)
	if err != nil {
		return nil, err
	}

	updateBson, _ := bson.Marshal(update)
	updateMap := make(bson.M)
	_ = bson.Unmarshal(updateBson, &updateMap)
	delete(updateMap, "_id")

	finalUpdate := bson.M{"$set": updateMap}

	pId, _ := primitive.ObjectIDFromHex(id)
	pFilter := filter.BaseFilter{
		Equals: map[string]interface{}{
			"_id":     pId,
			"deleted": false,
		},
	}
	builtFilter := filter.BuildBSONFilter(pFilter)
	result, err := b.Database.Collection(string(b.Collection)).UpdateOne(ctx, builtFilter, finalUpdate)
	if err != nil {
		return nil, err
	}

	return &result.MatchedCount, nil
}

func (b *BaseRepository[T, U]) UpdateOneByFilter(ctx context.Context, pFilter filter.BaseFilter, model T) (*int64, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.UpdateOneByFilter.%s", uType.Name()))
	defer span.End()

	if fValid, fErr := b.IsValidMandatoryFilters(ctx, pFilter, nil); !fValid {
		return nil, fErr
	}

	update, err := b.Mapper.ToEntity(model)
	if err != nil {
		return nil, err
	}

	updateBson, _ := bson.Marshal(update)
	updateMap := make(bson.M)
	_ = bson.Unmarshal(updateBson, &updateMap)
	delete(updateMap, "_id")

	finalUpdate := bson.M{"$set": updateMap}

	pFilter.Equals["deleted"] = false
	builtFilter := filter.BuildBSONFilter(pFilter)

	result, err := b.Database.Collection(string(b.Collection)).UpdateOne(ctx, builtFilter, finalUpdate)
	if err != nil {
		return nil, err
	}

	return &result.MatchedCount, nil
}

func (b *BaseRepository[T, U]) IsValidMandatoryFilters(ctx context.Context, pFilter filter.BaseFilter, project []string) (bool, error) {
	uType := reflect.TypeOf((*U)(nil)).Elem()
	_, span := otel.Tracer("BaseRepository").Start(ctx, fmt.Sprintf("BaseRepository.IsValidMandatoryFilters.%s", uType.Name()))
	defer span.End()

	var alwaysFields []string

	for i := 0; i < uType.NumField(); i++ {
		field := uType.Field(i)
		tag := field.Tag.Get("filter")

		if tag == "always" {
			bsonTag := field.Tag.Get("bson")
			bsonName := strings.Split(bsonTag, ",")[0]
			alwaysFields = append(alwaysFields, bsonName)
		}
	}

	if len(alwaysFields) == 0 {
		return true, nil
	}

	verified := false
	for k, _ := range pFilter.Equals {
		for _, v := range alwaysFields {
			if k == v {
				verified = true
				break
			}
		}
	}

	if !verified {
		return false, fmt.Errorf("invalid filter in %s", uType.Name())
	}

	return true, nil
}
