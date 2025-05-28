package shared

import (
	"github.com/kelseyhightower/envconfig"
	_ "github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

var lockConfig sync.Mutex
var configInstance *Config

func NewConfig() ConfigInterface {
	if configInstance == nil {
		lockConfig.Lock()
		defer lockConfig.Unlock()
		if configInstance == nil {
			configInstance = &Config{}
			err := envconfig.Process("", configInstance)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return configInstance
}

type Config struct {
	Port    uint16 `envconfig:"PORT" required:"true"`
	Timeout uint32 `envconfig:"TIMEOUT" required:"true"`

	MongoDbAuthSource       string `envconfig:"MONGODB_AUTH_SOURCE" required:"true"`
	MongoDbDatabaseName     string `envconfig:"MONGODB_DATABASE_NAME" required:"true"`
	MongoDbDatabaseHost     string `envconfig:"MONGODB_HOST" required:"true"`
	MongoDbMaxIdleTimeout   int    `envconfig:"MONGODB_MAX_IDLE_TIMEOUT" required:"true"`
	MongoDbMaxPoolSize      int    `envconfig:"MONGODB_MAX_POOL_SIZE" required:"true"`
	MongoDbMinPoolSize      int    `envconfig:"MONGODB_MIN_POOL_SIZE" required:"true"`
	MongoDbPassword         string `envconfig:"MONGODB_PASSWORD" required:"true"`
	MongoDbPort             int    `envconfig:"MONGODB_PORT" required:"true"`
	MongoDbUser             string `envconfig:"MONGODB_USER" required:"true"`
	MongoDbWaitQueueTimeout int    `envconfig:"MONGODB_WAIT_QUEUE_TIMEOUT" required:"true"`

	LogStashUrl string `envconfig:"LOG_STASH_URL" required:"true"`
}

func (c *Config) GetPort() uint16 {
	return c.Port
}

func (c *Config) GetTimeout() uint32 {
	return c.Timeout
}

func (c *Config) GetMongoDbAuthSource() string {
	return c.MongoDbAuthSource
}

func (c *Config) GetMongoDbDatabaseName() string {
	return c.MongoDbDatabaseName
}

func (c *Config) GetMongoDbDatabaseHost() string {
	return c.MongoDbDatabaseHost
}

func (c *Config) GetMongoDbMaxIdleTimeout() int {
	return c.MongoDbMaxIdleTimeout
}

func (c *Config) GetMongoDbMaxPoolSize() int {
	return c.MongoDbMaxPoolSize
}

func (c *Config) GetMongoDbMinPoolSize() int {
	return c.MongoDbMinPoolSize
}

func (c *Config) GetMongoDbPassword() string {
	return c.MongoDbPassword
}

func (c *Config) GetMongoDbPort() int {
	return c.MongoDbPort
}

func (c *Config) GetMongoDbUser() string {
	return c.MongoDbUser
}

func (c *Config) GetMongoDbWaitQueueTimeout() int {
	return c.MongoDbWaitQueueTimeout
}

func (c *Config) GetLogStashUrl() string {
	return c.LogStashUrl
}
