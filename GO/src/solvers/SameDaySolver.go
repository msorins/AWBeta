package solvers

import (
	"time"
	"github.com/gocolly/colly"
	"fmt"
)

type SameDaySolver struct {
	awb string
	url string
	LastSolverResponse SolverResponse
	lastUpdateCheck	time.Time
}

func SameDaySolverBuilder(awb string) ISolver {
	sds := SameDaySolver{}
	sds.awb = awb
	sds.url = "https://www.sameday.ro/awb-tracking?awb="
	return &sds
}

func (solver *SameDaySolver) updateStatuses() SolverResponse {

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")

	return SOLVER_NOT_IMPLEMENTED
}

func (solver *SameDaySolver) GetStatuses() ([]string, SolverResponse) {
	return nil, SOLVER_NOT_IMPLEMENTED
}

func (solver *SameDaySolver) GetLastStatus() ([]string, SolverResponse) {
	return nil, SOLVER_NOT_IMPLEMENTED
}

func (solver *SameDaySolver) GetAwb() string {
	return ""
}
