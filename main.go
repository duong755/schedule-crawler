package main

import (
	"sync"

	"schedule.crawler/crawler"
)

func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		crawler.Schedule()
		defer waitGroup.Done()
	}()
	go func() {
		crawler.Class()
		defer waitGroup.Done()
	}()
	waitGroup.Wait()
}
