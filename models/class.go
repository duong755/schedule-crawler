package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	Id               primitive.ObjectID `bson:"_id," json:"_id"`
	SubjectId        string             `bson:"subjectId" json:"subjectId"`
	SubjectName      string             `bson:"subjectName" json:"subjectName"`
	Credit           int8               `bson:"credit" json:"credit"`
	ClassId          string             `bson:"classId" json:"classId"`
	Teacher          string             `bson:"teacher" json:"teacher"`
	NumberOfStudents int8               `bson:"numberOfStudents" json:"numberOfStudents"`
	Session          string             `bson:"session" json:"session"` // "sáng" hoặc "chiều" hoặc "tối"
	WeekDay          int8               `bson:"weekDay" json:"weekDay"`
	Periods          []int8             `bson:"periods" json:"periods"`
	Place            string             `bson:"place" json:"place"`
	Note             string             `bson:"note" json:"note"` // "CL" (chung lớp), số (nhóm)
	Students         []string           `bson:"students" json:"students"`
}
