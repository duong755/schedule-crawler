package crawler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schedule.crawler/models"
)

type OnlyObjectId struct {
	Id primitive.ObjectID `bson:"_id"`
}

func buildQuery(semesterId string, page, offset int) string {
	return fmt.Sprintf("?SinhvienLmh[term_id]=%s&SinhvienLmh_page=%d&ajax=sinhvien-lmh-grid&pageSize=%d&r=sinhvienLmh/admin", semesterId, page, offset)
}

func Schedule(dbcontext context.Context, client *mongo.Client) {
	const semesterSelector string = "#SinhvienLmh_term_id"
	// const tableSelector string = "#sinhvien-lmh-grid > table.items > tbody"
	const lastPageSelector string = "#sinhvien-lmh-grid #yw0 > li.last > a[href]"
	const rootUrl string = "http://112.137.129.87/qldt/index.php"

	scheduleCollection := client.Database("uet").Collection("schedule")
	scheduleCollection.Drop(dbcontext)

	var lastPage int64 = 0
	var semesterId string = ""

	var semesterCollector *colly.Collector = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	var lastPageCollector *colly.Collector = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	var scheduleCollector *colly.Collector = colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	scheduleCollector.OnHTML("tbody", func(htmlTableBodyElement *colly.HTMLElement) {
		fmt.Println("Collecting schedules...")
		numberOfRows := len(htmlTableBodyElement.ChildTexts("tr"))
		documents := make([]interface{}, 0, numberOfRows)

		htmlTableBodyElement.ForEach("tr", func(rowIndex int, htmlTRowElement *colly.HTMLElement) {
			schedule := models.Schedule{}

			htmlTRowElement.ForEach("td", func(cellIndex int, htmlTCellElement *colly.HTMLElement) {
				// ignore case 0
				switch cellIndex {
				case 1:
					schedule.StudentId = htmlTCellElement.Text
				case 2:
					schedule.StudentName = htmlTCellElement.Text
				case 3:
					schedule.StudentBirthday = htmlTCellElement.Text
				case 4:
					schedule.StudentCourse = htmlTCellElement.Text
				case 5:
					schedule.ClassId = htmlTCellElement.Text
				case 7:
					schedule.ClassNote = htmlTCellElement.Text
				case 9:
					schedule.StudentNote = htmlTCellElement.Text
				}
			})

			bsonDocument := bson.D{
				{Key: "studentId", Value: schedule.StudentId},
				{Key: "studentName", Value: schedule.StudentName},
				{Key: "studentBirthday", Value: schedule.StudentBirthday},
				{Key: "studentCourse", Value: schedule.StudentCourse},
				{Key: "classId", Value: schedule.ClassId},
				{Key: "classNote", Value: schedule.ClassNote},
				{Key: "studentNote", Value: schedule.StudentNote},
			}
			documents = append(documents, bsonDocument)
			fmt.Printf("\r Scanned %d/%d", rowIndex+1, numberOfRows)
		})
		fmt.Println()

		_, err := scheduleCollection.InsertMany(context.TODO(), documents)
		if err != nil {
			panic(err)
		}
	})

	lastPageCollector.OnHTML(lastPageSelector, func(htmlLinkElement *colly.HTMLElement) {
		pageRegexp, _ := regexp.Compile(`(?:SinhvienLmh_page\=)(\d+)`)
		href := htmlLinkElement.Attr("href")
		matchedResults := pageRegexp.FindStringSubmatch(href)
		lastPageString := matchedResults[1]
		lastPage, _ = strconv.ParseInt(lastPageString, 10, 64)
		requestQueue, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10})
		for page := 1; page <= int(lastPage); page++ {
			requestQueue.AddURL(rootUrl + buildQuery(semesterId, page, 5000))
		}
		requestQueue.Run(scheduleCollector)
	})

	semesterCollector.OnHTML(semesterSelector, func(htmlSelectElement *colly.HTMLElement) {
		fmt.Println("Searching for semester id")
		htmlSelectElement.ForEach("option", func(htmlOptionElementIndex int, htmlOptionElement *colly.HTMLElement) {
			if htmlOptionElementIndex == 1 {
				semesterId = htmlOptionElement.Attr("value")
				if semesterId == "" {
					fmt.Println("No semester id found")
				} else {
					fmt.Printf("Found semester id: %s\n", semesterId)
					lastPageCollector.Visit(rootUrl + buildQuery(semesterId, 1, 5000))
				}
			}
		})
	})

	semesterCollector.Visit(rootUrl)
	semesterCollector.Wait()
}
