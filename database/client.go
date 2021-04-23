package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Client() (context.Context, *mongo.Client) {
	dbcontext := context.Background()

	client, err := mongo.Connect(dbcontext, options.Client().ApplyURI(DATABASE_URL))
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to database")
	}

	return dbcontext, client
}
