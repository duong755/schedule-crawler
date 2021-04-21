package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	Id primitive.ObjectID `bson:"_id" json:"id"`
}
