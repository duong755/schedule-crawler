package crawler

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func Schedule() {
	var scheduleCollector *colly.Collector = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	selector := ""
	url := "http://112.137.129.87/qldt/"
	scheduleCollector.OnRequest(func(req *colly.Request) {
		fmt.Printf("Sending request to %s ...\n", req.URL)
	})
	scheduleCollector.OnResponse(func(res *colly.Response) {
		fmt.Printf("Received response from %s\n", res.Request.URL.String())
		fmt.Printf("%+v\n", res.Headers)
		fmt.Printf("%d\n", res.StatusCode)
	})
	scheduleCollector.OnHTML(selector, func(html *colly.HTMLElement) {})

	scheduleCollector.Visit(url)
}
