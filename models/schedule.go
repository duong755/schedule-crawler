package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Schedule struct {
	Id              primitive.ObjectID `bson:"_id" json:"_id"`
	StudentId       string             `bson:"studentId" json:"studentId"`
	StudentName     string             `bson:"studentName" json:"studentName"`
	StudentBirthday string             `bson:"studentBirthday" json:"studentBirthday"`
	StudentCourse   string             `bson:"studentCourse" json:"studentCourse"`
	ClassId         string             `bson:"classId" json:"classId"`
	ClassNote       string             `bson:"classNote" json:"classNote"`
	StudentNote     string             `bson:"studentNote" json:"studentNote"`
}
