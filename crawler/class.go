package crawler

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func Class() {
	classCollector := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	selector := "table:nth-child(4)"
	url := "http://112.137.129.115/tkb/listbylist.php"
	classCollector.OnRequest(func(req *colly.Request) {
		fmt.Printf("Sending request to %s ...\n", req.URL)
	})
	classCollector.OnResponse(func(res *colly.Response) {
		fmt.Printf("Received response from %s\n", res.Request.URL.String())
		fmt.Printf("%+v\n", res.Headers)
		fmt.Printf("%d\n", res.StatusCode)
	})
	classCollector.OnHTML(selector, func(matchedTable *colly.HTMLElement) {
		matchedTable.ForEach("tr", func(index int, tableRow *colly.HTMLElement) {
			fmt.Println(tableRow.Text)
		})
	})

	classCollector.Visit(url)
}
