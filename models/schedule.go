package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Student struct {
	Id       string `bson:"id" json:"id"`
	Name     string `bson:"name" json:"name"`
	Birthday string `bson:"birthday" json:"birthday"`
	Course   string `bson:"course" json:"course"`
	Note     string `bson:"note" json:"note"`
}

type Schedule struct {
	Id      primitive.ObjectID `bson:"id" json:"id"`
	Student Student            `bson:"student" json:"student"`
	Classes []primitive.ObjectID
}
