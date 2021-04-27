package main

import (
	"schedule.crawler/crawler"
	"schedule.crawler/database"
)

func main() {
	dbcontext, client := database.Client()
	defer client.Disconnect(dbcontext)

	go crawler.Class(dbcontext, client)
	crawler.Schedule(dbcontext, client)
}
