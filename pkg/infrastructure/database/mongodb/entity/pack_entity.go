package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type PackEntity struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	MaxSize     int                `bson:"maxSize,omitempty"`
	BaseEntity  `bson:",inline"`
}
