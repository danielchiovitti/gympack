package shared

type ConfigInterface interface {
	GetPort() uint16
	GetTimeout() uint32

	GetMongoDbAuthSource() string
	GetMongoDbDatabaseName() string
	GetMongoDbDatabaseHost() string
	GetMongoDbMaxIdleTimeout() int
	GetMongoDbMaxPoolSize() int
	GetMongoDbMinPoolSize() int
	GetMongoDbPassword() string
	GetMongoDbPort() int
	GetMongoDbUser() string
	GetMongoDbWaitQueueTimeout() int

	GetLogStashUrl() string
}
