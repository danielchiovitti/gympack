package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BaseEntity struct {
	CreatedAt   time.Time          `bson:"created_at,omitempty"`
	CreatedById primitive.ObjectID `bson:"created_by_id,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
	UpdatedById primitive.ObjectID `bson:"updated_by_id,omitempty"`
	Deleted     bool               `bson:"deleted"`
	DeletedAt   time.Time          `bson:"deleted_at,omitempty"`
	DeletedById primitive.ObjectID `bson:"deleted_by_id,omitempty"`
}
