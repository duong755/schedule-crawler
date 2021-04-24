package crawler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"schedule.crawler/models"
)

func Class(dbcontext context.Context, client *mongo.Client) {
	classCollector := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	selector := "table:nth-child(4) > tbody"
	url := "http://112.137.129.115/tkb/listbylist.php"
	classCollector.OnRequest(func(req *colly.Request) {
		fmt.Printf("Sending request to %s ...\n", req.URL)
	})
	classCollector.OnResponse(func(res *colly.Response) {
		fmt.Printf("Received response from %s\n", res.Request.URL.String())
		fmt.Printf("%d\n", res.StatusCode)
	})
	classCollector.OnHTML(selector, func(matchedTable *colly.HTMLElement) {
		firstPeriodRegexp, _ := regexp.Compile(`^\d+`)
		lastPeriodRegexp, _ := regexp.Compile(`\d+$`)
		classCollection := client.Database("uet").Collection("class")
		classCollection.Drop(dbcontext)

		numberOfRows := len(matchedTable.ChildTexts("tr"))
		documents := make([]interface{}, 0, numberOfRows)

		fmt.Println("Collecting classes...")
		matchedTable.ForEach("tr", func(rowIndex int, tableRow *colly.HTMLElement) {
			if rowIndex == 0 {
				// ignore table header
				return
			}
			class := models.Class{}
			tableRow.ForEach("td", func(cellIndex int, tableCell *colly.HTMLElement) {
				switch cellIndex {
				case 1:
					class.SubjectId = tableCell.Text
				case 2:
					class.SubjectName = tableCell.Text
				case 3:
					credit, convertCreditErr := strconv.ParseInt(tableCell.Text, 10, 8)
					if convertCreditErr == nil {
						class.Credit = int8(credit)
					}
				case 4:
					class.ClassId = tableCell.Text
				case 5:
					class.Teacher = tableCell.Text
				case 6:
					numberOfStudents, convertNumberOfStudentsErr := strconv.ParseInt(tableCell.Text, 10, 8)
					if convertNumberOfStudentsErr == nil {
						class.NumberOfStudents = int8(numberOfStudents)
					}
				case 7:
					class.Session = tableCell.Text
				case 8:
					weekDay, convertWeekDayErr := strconv.ParseInt(tableCell.Text, 10, 8)
					if convertWeekDayErr == nil {
						class.WeekDay = int8(weekDay)
					} else {
						class.WeekDay = 8
					}
				case 9:
					firstPeriodMatch := firstPeriodRegexp.FindString(tableCell.Text)
					lastPeriodMatch := lastPeriodRegexp.FindString(tableCell.Text)
					if firstPeriodMatch == "" || lastPeriodMatch == "" {
						break
					} else {
						firstPeriod, _ := strconv.ParseInt(firstPeriodMatch, 10, 8)
						lastPeriod, _ := strconv.ParseInt(lastPeriodMatch, 10, 8)

						periods := make([]int8, 0)
						for period := firstPeriod; period <= lastPeriod; period++ {
							periods = append(periods, int8(period))
						}
						class.Periods = periods
					}
				case 10:
					class.Place = tableCell.Text
				case 11:
					class.Note = tableCell.Text
				}
			})

			periods := bson.A{}
			for _, period := range class.Periods {
				periods = append(periods, period)
			}
			bsonDocument := bson.D{
				{Key: "subjectId", Value: class.SubjectId},
				{Key: "subjectName", Value: class.SubjectName},
				{Key: "credit", Value: class.Credit},
				{Key: "classId", Value: class.ClassId},
				{Key: "teacher", Value: class.Teacher},
				{Key: "numberOfStudents", Value: class.NumberOfStudents},
				{Key: "session", Value: class.Session},
				{Key: "weekDay", Value: class.WeekDay},
				{Key: "periods", Value: periods},
				{Key: "place", Value: class.Place},
				{Key: "note", Value: class.Note},
			}
			documents = append(documents, bsonDocument)
			fmt.Printf("\r Scanned %d/%d rows", rowIndex+1, numberOfRows)
		})

		fmt.Println()
		_, err := classCollection.InsertMany(dbcontext, documents)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	})

	classCollector.Visit(url)
}
