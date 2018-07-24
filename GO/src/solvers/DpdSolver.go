package solvers

import (
	"time"
	"fmt"
	"github.com/gocolly/colly"
	"sort"
)

type dpdSolver struct {
	awb string
	url string
	Statuses []AwbDpdCheckpoint
	LastSolverResponse SolverResponse
	lastUpdateCheck	time.Time
}

type AwbDpdCheckpoint struct {
	Index int
	Status string
	Date string
	Time string
	Location string
}


func DpdSolverBuile(awb string) ISolver {
	dpd := dpdSolver{}
	dpd.awb = awb
	dpd.url = "https://tracking.dpd.ro/?shipmentNumber=" + dpd.awb
	return &dpd
}

func (solver *dpdSolver) updateStatuses() SolverResponse {
	// Check to see if the request is already cached
	if time.Since(solver.lastUpdateCheck).Minutes() < TIME_BETWEEN_REQUEST_MIN {
		return SOLVER_CACHED
	}

	// Crawler object
	c := colly.NewCollector()

	// Find and visit all links
	count := 0
	var errMsg error
	statuses := []AwbDpdCheckpoint{}
	crtStatus := AwbDpdCheckpoint{}
	index := 0

	c.OnHTML("tr>td", func(e *colly.HTMLElement) {
		count += 1
		elem := e.Text

		switch count{
		case 1:
			// Date
			crtStatus.Date = elem
			break
		case 2:
			// Time
			crtStatus.Time = elem
			break
		case 3:
			// Blank column
			crtStatus.Status = elem
			break
		case 4:
			// Location
			index += 1

			crtStatus.Index = index
			if len(elem) > 3 {
				crtStatus.Location = elem[3:]
			} else {
				crtStatus.Location = elem
			}


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

	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Index > statuses[j].Index
	})

	solver.Statuses = statuses
	solver.lastUpdateCheck = time.Now()
	solver.LastSolverResponse = SOLVER_OK
	return SOLVER_OK
}

func (solver *dpdSolver) GetStatuses() ([]string, SolverResponse) {
	responseCode := solver.updateStatuses()
	updatedStatuses := solver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "Dpd: These are all the steps taken by your Dpd package")
		for _, status := range updatedStatuses {
			results = append(results, fmt.Sprintf("%s in %s at %s %s", status.Status, status.Location, status.Date, status.Time))
		}
	} else {
		results = append(results, "Dpd: Could not found any records for your AWB")
	}


	return results, responseCode
}

func (solver *dpdSolver) GetLastStatus() ([]string, SolverResponse) {
	responseCode := solver.updateStatuses()
	updatedStatuses := solver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "Dpd: Successfully found the latest status of your Dpd package")
		results = append(results, fmt.Sprintf("%s in %s at %s %s", updatedStatuses[0].Status, updatedStatuses[0].Location, updatedStatuses[0].Date, updatedStatuses[0].Time))
	} else {
		results = append(results, "Dpd: Could not found any records for your AWB")
	}

	return results, responseCode
}

func (solver *dpdSolver) GetAwb() string {
	return solver.awb
}
