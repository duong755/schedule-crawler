package main

import (
	"schedule.crawler/crawler"
	"schedule.crawler/database"
)

func main() {
	dbcontext, client := database.Client()

        crawler.Class(dbcontext, client)
	crawler.Schedule(dbcontext, client)

	defer client.Disconnect(dbcontext)
}
