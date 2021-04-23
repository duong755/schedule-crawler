package main

import (
	"schedule.crawler/crawler"
	"schedule.crawler/database"
)

func main() {
	dbcontext, client := database.Client()

	crawler.Schedule(dbcontext, client)
	crawler.Class(dbcontext, client)

	defer client.Disconnect(dbcontext)
}
