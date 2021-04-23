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
		classCollection := client.Database("uet").Collection("class")
		scheduleCollection := client.Database("uet").Collection("schedule")

		fmt.Println("Collecting schedules...")
		numberOfRows := len(htmlTBodyElement.ChildTexts("tr"))
		documents := make([]interface{}, 0, numberOfRows)

		htmlTBodyElement.ForEach("tr", func(rowIndex int, htmlTRowElement *colly.HTMLElement) {
			student := models.Student{}
			schedule := models.Schedule{Classes: make([]primitive.ObjectID, 0, 4)}

			htmlTRowElement.ForEach("td", func(cellIndex int, htmlTCellElement *colly.HTMLElement) {
				classId := ""
				classNote := "" // "CL", "1", "2", ...
				// ignore case 0
				switch cellIndex {
				case 1:
					student.Id = htmlTCellElement.Text
				case 2:
					student.Name = htmlTCellElement.Text
				case 3:
					student.Birthday = htmlTCellElement.Text
				case 4:
					student.Course = htmlTCellElement.Text
				case 5:
					classId = htmlTCellElement.Text
				case 7:
					classNote = htmlTCellElement.Text
				case 9:
					student.Note = htmlTCellElement.Text
				}

				schedule.Student = student
				filter := bson.D{
					{Key: "classId", Value: classId},
					{Key: "note", Value: bson.D{{Key: "$in", Value: bson.A{"CL", classNote}}}},
				}
				cursorResult, queryErr := classCollection.Find(dbcontext, filter)
				if queryErr != nil {
					panic(queryErr)
				}
				for cursorResult.Next(dbcontext) {
					classItem := OnlyObjectId{}
					cursorResult.Decode(&classItem)
					schedule.Classes = append(schedule.Classes, classItem.Id)
				}
			})

			classes := bson.A{}
			for _, class := range schedule.Classes {
				classes = append(classes, class)
			}

			bsonDocument := bson.D{
				{Key: "student", Value: bson.D{
					{Key: "id", Value: schedule.Student.Id},
					{Key: "name", Value: schedule.Student.Name},
					{Key: "birthday", Value: schedule.Student.Birthday},
					{Key: "course", Value: schedule.Student.Course},
					{Key: "note", Value: schedule.Student.Note},
				}},
				{Key: "classes", Value: classes},
			}
			documents = append(documents, bsonDocument)
		})

		results, err := scheduleCollection.InsertMany(dbcontext, documents)
		fmt.Printf("%v\n", results)
		if err != nil {
			panic(err)
		}
	})

	scheduleCollector.Visit(rootUrl)
	for page := 1; !isLastPage; page++ {
		queryString := fmt.Sprintf("?SinhvienLmh[term_id]=%s&SinhvienLmh_page=%d&ajax=sinhvien-lmh-grid&pageSize=%d&r=sinhvienLmh/admin", semesterId, page, 5000)
		scheduleCollector.Visit(rootUrl + queryString)
	}
}
