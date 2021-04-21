package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"schedule.crawler/database"
)

func main() {
	dbs, err := database.Client().ListDatabaseNames(context.TODO(), bson.D{})

	if err != nil {
		panic(err)
	}

	for _, dbname := range dbs {
		fmt.Println(dbname)
	}
}
