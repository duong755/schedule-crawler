package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"schedule.crawler/crawler"
	"schedule.crawler/database"
)

func main() {
	// load from .env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal(envErr)
	}
	uri := os.Getenv("DATABASE_URL")
	dbcontext, client := database.Client(uri)
	defer client.Disconnect(dbcontext)

	go crawler.Class(dbcontext, client)
	crawler.Schedule(dbcontext, client)
}
