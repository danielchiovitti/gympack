package provider

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gympack/pkg/domain/model"
	"gympack/pkg/shared"
	"net/url"
	"strconv"
	"sync"
)

var mongodbProviderLock sync.Mutex
var mongodbClientLock sync.Mutex
var mongodbProviderInstance *MongoDbProvider
var mongodbClientInstance *mongo.Client

type MongoDbProvider struct {
	databaseOptions *model.MongoDbOptions
	config          shared.ConfigInterface
}

func NewMongoDbProvider(config shared.ConfigInterface) MongoDbProviderInterface {
	if mongodbProviderInstance == nil {
		mongodbProviderLock.Lock()
		defer mongodbProviderLock.Unlock()
		if mongodbProviderInstance == nil {
			opts := []model.MongoDbOptionsFunc{
				WithHost(config.GetMongoDbDatabaseHost()),
				WithPort(config.GetMongoDbPort()),
				WithDatabaseName(config.GetMongoDbDatabaseName()),
				WithUser(config.GetMongoDbUser()),
				WithPassword(config.GetMongoDbPassword()),
				WithMinPoolSize(config.GetMongoDbMinPoolSize()),
				WithMaxPoolSize(config.GetMongoDbMaxPoolSize()),
				WithMaxIdleTimeMS(config.GetMongoDbMaxIdleTimeout()),
				WithWaitQueueTimeoutMS(config.GetMongoDbWaitQueueTimeout()),
				WithAuthSource(config.GetMongoDbAuthSource()),
			}
			o := DatabaseDefaultOpts()
			for _, fn := range opts {
				fn(o)
			}
			mongodbProviderInstance = &MongoDbProvider{
				databaseOptions: o,
				config:          config,
			}
		}
	}
	return mongodbProviderInstance
}

func WithHost(host string) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.Host = host
	}
}

func WithPort(port int) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.Port = port
	}
}

func WithDatabaseName(databaseName string) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.DatabaseName = databaseName
	}
}

func WithUser(user string) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.User = user
	}
}

func WithPassword(password string) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.Password = password
	}
}

func WithMinPoolSize(minPoolSize int) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.MinPoolSize = minPoolSize
	}
}

func WithMaxPoolSize(maxPoolSize int) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.MaxPoolSize = maxPoolSize
	}
}

func WithMaxIdleTimeMS(maxIdleTimeMS int) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.MaxIdleTimeMS = maxIdleTimeMS
	}
}

func WithConnectTimeoutMS(connectTimeoutMS int) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.ConnectTimeoutMS = connectTimeoutMS
	}
}

func WithWaitQueueTimeoutMS(waitQueueTimeoutMS int) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.WaitQueueTimeoutMS = waitQueueTimeoutMS
	}
}

func WithAuthSource(authSource string) model.MongoDbOptionsFunc {
	return func(opt *model.MongoDbOptions) {
		opt.AuthSource = authSource
	}
}

func DatabaseDefaultOpts() *model.MongoDbOptions {
	return &model.MongoDbOptions{
		Host: "127.0.0.1",
		Port: 3306,
	}
}

func (d *MongoDbProvider) GetMongoDbClient() (*mongo.Client, error) {
	if mongodbClientInstance != nil {
		return mongodbClientInstance, nil
	}

	query := url.Values{}

	if d.databaseOptions.MinPoolSize != 0 {
		query.Add("minPoolSize", strconv.Itoa(d.databaseOptions.MinPoolSize))
	}

	if d.databaseOptions.MaxPoolSize > 0 {
		query.Add("maxPoolSize", strconv.Itoa(d.databaseOptions.MaxPoolSize))
	}

	if d.databaseOptions.MaxIdleTimeMS > 0 {
		query.Add("maxIdleTimeMS", strconv.Itoa(d.databaseOptions.MaxIdleTimeMS))
	}

	if d.databaseOptions.ConnectTimeoutMS > 0 {
		query.Add("connectTimeoutMS", strconv.Itoa(d.databaseOptions.ConnectTimeoutMS))
	}

	if d.databaseOptions.WaitQueueTimeoutMS > 0 {
		query.Add("waitQueueTimeoutMS", strconv.Itoa(d.databaseOptions.WaitQueueTimeoutMS))
	}

	if d.databaseOptions.AuthSource != "" {
		query.Add("authSource", d.databaseOptions.AuthSource)
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?%s", d.databaseOptions.User, d.databaseOptions.Password, d.databaseOptions.Host, d.databaseOptions.Port, d.databaseOptions.DatabaseName, query.Encode())

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	mongoOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverApi)

	client, err := mongo.Connect(context.Background(), mongoOptions)
	if err != nil {
		return nil, err
	}

	if mongodbClientInstance == nil {
		mongodbClientLock.Lock()
		defer mongodbClientLock.Unlock()
		if mongodbClientInstance == nil {
			mongodbClientInstance = client
		}
	}

	return mongodbClientInstance, nil
}
