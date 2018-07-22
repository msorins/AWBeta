package solvers

import (
	"time"
	"github.com/gocolly/colly"
	"fmt"
	"strings"
)

type awbSameDaySolver struct {
	awb string
	url string
	Statuses []AwbSameDayCheckpoint
	LastSolverResponse SolverResponse
	lastUpdateCheck	time.Time
}

type AwbSameDayCheckpoint struct {
	Index int
	Status string
	Date string
	Time string
	Location string
}

func SameDaySolverBuilder(awb string) ISolver {
	sds := awbSameDaySolver{}
	sds.awb = awb
	sds.url = "https://www.sameday.ro/awb-tracking?awb=" + awb
	return &sds
}

func (solver *awbSameDaySolver) updateStatuses() SolverResponse {
	// Check to see if the request is already cached
	if time.Since(solver.lastUpdateCheck).Minutes() < TIME_BETWEEN_REQUEST_MIN {
		return SOLVER_CACHED
	}

	// Crawler object
	c := colly.NewCollector()

	// Find and visit all links
	count := 0
	var errMsg error
	statuses := []AwbSameDayCheckpoint{}
	crtStatus := AwbSameDayCheckpoint{}
	index := 101

	c.OnHTML("tr>td", func(e *colly.HTMLElement) {
		count += 1
		elem := e.Text

		switch count{
		case 1:
			// Date and time
			crtStatus.Date = strings.Split(elem, " ")[1]
			crtStatus.Time = strings.Split(elem, " ")[2]
			break
		case 2:
			// Location
			crtStatus.Location = elem
			break
		case 3:
			// Blank column
			break
		case 4:
			// Awb
			break
		case 5:
			// Status
			index -= 1

			crtStatus.Status = elem
			crtStatus.Index = index
			statuses = append(statuses, crtStatus)
			count = 0
			break
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(response *colly.Response, e error) {
		errMsg = e
	})

	// Start crawling
	c.Visit(solver.url)

	// Error handling
	if errMsg != nil {
		fmt.Println(errMsg)
		return SOLVER_BAD_REQUEST
	}

	if len(statuses) == 0 {
		return SOLVER_AWB_INCORRECT
	}

	solver.Statuses = statuses
	solver.lastUpdateCheck = time.Now()
	solver.LastSolverResponse = SOLVER_OK
	return SOLVER_OK
}

func (solver *awbSameDaySolver) GetStatuses() ([]string, SolverResponse) {
	responseCode := solver.updateStatuses()
	updatedStatuses := solver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "SameDay: These are all the steps taken by your SameDay package")
		for _, status := range updatedStatuses {
			results = append(results, fmt.Sprintf("%s in %s at %s %s", status.Status, status.Location, status.Date, status.Time))
		}
	} else {
		results = append(results, "SameDay: Could not found any records for your AWB")
	}


	return results, responseCode
}

func (solver *awbSameDaySolver) GetLastStatus() ([]string, SolverResponse) {
	responseCode := solver.updateStatuses()
	updatedStatuses := solver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "SameDay: Successfully found the latest status of your SameDay package")
		results = append(results, fmt.Sprintf("%s in %s at %s %s", updatedStatuses[0].Status, updatedStatuses[0].Location, updatedStatuses[0].Date, updatedStatuses[0].Time))
	} else {
		results = append(results, "SameDay: Could not found any records for your AWB")
	}

	return results, responseCode
}

func (solver *awbSameDaySolver) GetAwb() string {
	return solver.awb
}
