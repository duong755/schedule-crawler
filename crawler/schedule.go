package crawler

import (
	"context"
	"fmt"
	"net/url"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Schedule(dbcontext context.Context, client *mongo.Client) {
	var scheduleCollector *colly.Collector = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	semesterSelector := "#SinhvienLmh_term_id"
	var semesterId string = ""
	tableSelector := "#sinhvien-lmh-grid > table.items > tbody"
	rootUrl := "http://112.137.129.87/qldt/index.php"

	scheduleCollector.OnRequest(func(req *colly.Request) {
		fmt.Printf("Sending request to %s ...\n", req.URL)
	})

	scheduleCollector.OnResponse(func(res *colly.Response) {
		fmt.Printf("Received response from %s\n", res.Request.URL.String())
		fmt.Printf("%+v\n", res.Headers)
		fmt.Printf("%d\n", res.StatusCode)
	})

	scheduleCollector.OnHTML(semesterSelector, func(htmlSelectElement *colly.HTMLElement) {
		fmt.Println("Searching for semester id")
		htmlSelectElement.ForEach("option", func(htmlOptionElementIndex int, htmlOptionElement *colly.HTMLElement) {
			if htmlOptionElementIndex == 1 {
				semesterId = htmlOptionElement.Attr("value")
				if semesterId == "" {
					fmt.Println("No semester id found")
				} else {
					fmt.Printf("Found semster id: %s\n", semesterId)
				}
			}
		})
	})
	scheduleCollector.OnHTML(tableSelector, func(htmlTBodyElement *colly.HTMLElement) {
		if htmlTBodyElement.DOM.ChildrenFiltered("span.empty").First().Text() == "" {
			return
		}
		fmt.Println("Collecting schedules...")
		htmlTBodyElement.ForEach("tr", func(rowIndex int, htmlTRowElement *colly.HTMLElement) {
			//
		})
	})

	scheduleCollector.Visit(rootUrl)
	url.QueryEscape(fmt.Sprintf("?SinhvienLmh[term_id]=%s&SinhvienLmh_page=%d&ajax=sinhvien-lmh-grid&pageSize=%d&r=sinhvienLmh/admin", semesterId, 1, 5000))
}
