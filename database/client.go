package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Client(uri string) (context.Context, *mongo.Client) {
	dbcontext := context.Background()
	dbcontext = context.WithValue(dbcontext, "startCrawlingTime", time.Now().UTC().Format(time.UnixDate))

	client, err := mongo.Connect(dbcontext, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	return dbcontext, client
}
