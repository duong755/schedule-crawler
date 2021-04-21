package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	Id primitive.ObjectID `bson:"_id" json:"_id"`
}
