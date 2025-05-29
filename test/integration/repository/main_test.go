package repository

import (
	"context"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"log"
	"os"
	"sync"
	"testing"
)

var (
	MongoContainer *mongodb.MongoDBContainer
	MongoOnce      sync.Once
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	MongoOnce.Do(func() {
		var err error
		MongoContainer, err = mongodb.Run(
			ctx,
			"mongo:6",
			mongodb.WithPassword("dc"),
			mongodb.WithUsername("dc"),
		)
		if err != nil {
			log.Fatalf("Failed to start MongoDB container: %s", err)
		}
	})

	code := m.Run()
	if MongoContainer != nil {
		if err := MongoContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate MongoDB container: %s", err)
		}
	}

	os.Exit(code)
}

func GetMongoContainer() *mongodb.MongoDBContainer {
	return MongoContainer
}
