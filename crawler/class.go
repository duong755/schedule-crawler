package crawler

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"schedule.crawler/models"
)

func Class(dbcontext context.Context, client *mongo.Client) {
	classCollector := colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	SELECTOR := "table:nth-child(4) > tbody"
	ROOT_URL := "http://112.137.129.115/tkb/listbylist.php"

	classCollection := client.Database("uet").Collection("class")
	classCollection.Indexes().CreateOne(dbcontext, mongo.IndexModel{
		Keys:    bson.D{{Key: "classId", Value: "hashed"}},
		Options: options.Index().SetName("classIdHashedIndex"),
	})

	var anyClass models.Class
	classCollection.FindOne(dbcontext, bson.D{}).Decode(&anyClass)
	defer func() {
		fmt.Println("Deleting outdated classes data...")
		classCollection.DeleteMany(dbcontext, bson.D{{Key: "crawledAt", Value: anyClass.CrawledAt}})
	}()

	classCollector.OnRequest(func(req *colly.Request) {
		fmt.Printf("Sending request to %s ...\n", req.URL)
	})

	classCollector.OnResponse(func(res *colly.Response) {
		fmt.Printf("Received response from %s\n", res.Request.URL.String())
		fmt.Printf("%d\n", res.StatusCode)
	})

	classCollector.OnHTML(SELECTOR, func(matchedTable *colly.HTMLElement) {
		firstPeriodRegexp, _ := regexp.Compile(`^\d+`)
		lastPeriodRegexp, _ := regexp.Compile(`\d+$`)

		numberOfRows := len(matchedTable.ChildTexts("tr"))
		documents := []interface{}{}

		fmt.Println("Collecting classes...")
		matchedTable.ForEach("tr", func(rowIndex int, tableRow *colly.HTMLElement) {
			if rowIndex == 0 {
				// ignore table header
				return
			}
			class := models.Class{}
			tableRow.ForEach("td", func(cellIndex int, tableCell *colly.HTMLElement) {
				class.CrawledAt = dbcontext.Value("startCrawlingTime").(string)
				// 0. ordinary number
				// 1. subject id
				// 2. subject name
				// 3. credit
				// 4. class id
				// 5. teacher
				// 6. number of students
				// 7. session
				// 8. week day
				// 9. periods
				// 10. place
				// 11. group
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

						periods := []int8{}
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
			for index := range class.Periods {
				periods = append(periods, class.Periods[index])
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
				{Key: "crawledAt", Value: class.CrawledAt},
			}
			documents = append(documents, bsonDocument)
			fmt.Printf("\rScanned %d/%d rows", rowIndex+1, numberOfRows)
		})

		fmt.Println()
		_, err := classCollection.InsertMany(dbcontext, documents)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	})

	classCollector.Visit(ROOT_URL)
}
