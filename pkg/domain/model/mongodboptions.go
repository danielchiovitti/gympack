package model

type MongoDbOptions struct {
	Host               string
	Port               int
	DatabaseName       string
	User               string
	Password           string
	MinPoolSize        int
	MaxPoolSize        int
	MaxIdleTimeMS      int
	ConnectTimeoutMS   int
	WaitQueueTimeoutMS int
	AuthSource         string
}

type MongoDbOptionsFunc func(opt *MongoDbOptions)
