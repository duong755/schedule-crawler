package crawler

import (
	"context"
	"fmt"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"schedule.crawler/models"
)

type OnlyObjectId struct {
	Id primitive.ObjectID `bson:"_id"`
}

func Schedule(dbcontext context.Context, client *mongo.Client) {
	var scheduleCollector *colly.Collector = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	client.Database("uet").Collection("schedule").Drop(dbcontext)
	semesterSelector := "#SinhvienLmh_term_id"
	var semesterId string = ""
	tableSelector := "#sinhvien-lmh-grid > table.items > tbody"
	rootUrl := "http://112.137.129.87/qldt/index.php"
	currentSchedulePage := "0"
	isLastPage := false

	scheduleCollector.OnRequest(func(req *colly.Request) {
		fmt.Printf("Sending request to %s ...\n", req.URL)
	})

	scheduleCollector.OnResponse(func(res *colly.Response) {
		fmt.Printf("Received response from %s\n", res.Request.URL.String())
		fmt.Printf("%+v\n", res.Headers)
		fmt.Printf("%d\n", res.StatusCode)
	})

	scheduleCollector.OnHTML(semesterSelector, func(htmlSelectElement *colly.HTMLElement) {
		if semesterId != "" {
			return
		}
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

	scheduleCollector.OnHTML("#yw0 li.page.selected > a", func(htmlLinkElement *colly.HTMLElement) {
		if currentSchedulePage == htmlLinkElement.Text {
			// stop crawling
			isLastPage = true
		} else {
			currentSchedulePage = htmlLinkElement.Text
		}
	})

	scheduleCollector.OnHTML(tableSelector, func(htmlTBodyElement *colly.HTMLElement) {
		if htmlTBodyElement.ChildText("span.empty") != "" {
			return
		}
		scheduleCollection := client.Database("uet").Collection("schedule")
		scheduleCollection.Drop(dbcontext)

		fmt.Println("Collecting schedules...")
		numberOfRows := len(htmlTBodyElement.ChildTexts("tr"))
		documents := make([]interface{}, 0, numberOfRows)

		if isLastPage {
			return
		}

		htmlTBodyElement.ForEach("tr", func(rowIndex int, htmlTRowElement *colly.HTMLElement) {
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
			fmt.Printf("\r Scanned %d/%d rows of page %s", rowIndex+1, numberOfRows, currentSchedulePage)
		})

		results, err := scheduleCollection.InsertMany(dbcontext, documents)
		fmt.Printf("%v\n", results)
		if err != nil {
			panic(err)
		}
	})

	scheduleCollector.Visit(rootUrl)
	scheduleCollector.Wait()
	for page := 1; !isLastPage; page++ {
		queryString := fmt.Sprintf("?SinhvienLmh[term_id]=%s&SinhvienLmh_page=%d&ajax=sinhvien-lmh-grid&pageSize=%d&r=sinhvienLmh/admin", semesterId, page, 5000)
		fmt.Printf("Scanning page %d\n", page)
		cloneScheduleCollector := scheduleCollector.Clone()
		cloneScheduleCollector.Visit(rootUrl + queryString)
		cloneScheduleCollector.Wait()
	}
}
