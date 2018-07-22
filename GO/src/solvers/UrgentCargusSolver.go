package solvers

import (
	"fmt"
	"time"
	"github.com/gocolly/colly"
	"strings"
	"sort"
)

type awbUrgentCargusSolver struct {
	awb string
	url string
	requestData map[string]string
	Statuses []AwbUrgentCargusCheckpoint
	LastSolverResponse SolverResponse
	lastUpdateCheck	time.Time
}


type AwbUrgentCargusCheckpoint struct {
	Index int
	Status string
	Date string
	Time string
	Location string
}

func UrgentCargusBuilder(awb string) ISolver {
	sds := awbUrgentCargusSolver{}
	sds.awb = awb
	sds.url = "https://app.urgentcargus.ro/Private/Tracking.aspx?CodBara=" + awb
	sds.requestData = map[string]string{}
	sds.requestData["__VIEWSTATE"] = "VT9Kd6nuu3XDtlEZG55hZt4rmcJYuEHKMi/T9upSMq3OXf8pA2YCmnzkNqLjQJsLwBzw3eJLGkP6S6JJ0NH3r3LJ+UYqhc5soePreJUIqQJcg6gwzGmbZVPMgvMSEBseYoTQAq2V5os1eQXakZPdCZq1fpDsiKTDcJ5Q3FlK42Pw8yiQsterrQ3mfPInCfqhaVHyc4+aDh51fcZ8ifvqbAJAekmqJobxz3xdXt+wo/QV0lkq5HqO8b1F0CzPVYY1dZv7suJ42RPhf0eRJuJ6VHRWlZRYrVSsY8gsbdhBdoovnJUJj8oOwEdOF5iVZ9M3gpxPbGrlnGBzdZ/5lVKTU0VTTTUrjc9iZRwWfBzU5G67MdMfLCNSSX7iiw22IxTjyKQ2oN7eiZ5p9R2G3kqIrHJMT9Cw169gcNrUfAzQMHIv01CxjuawHjWGIrXtjFVS5Hw4fhgGuNVZw7v7zXbYGEZghJNzYLxrwrTmGUMNQeYZdIwBgL0LNfi/ZXXm/10bd8OtR+G4L4H24XsyyfhfhmfFtM7AD1VJlhnR58GqWq+MGyMpgz9LTMzqxR9YkozuO6FW/jVxMdZoCqg94w5jzGv7BQJdo8vbR5df36xi+Huxi+gwuDiMOKL7PyZui+2YSdSF6o1Z5HatmABc4v2Zz+DZezce4BFS7QbWwpfDY1APbM79bg58zxUOVePm0YNjNBsVLl0rpxeiDRFkgh3pQ48cDVLKdgwu/ChnHBkz8G0Jn85Ikcc3PstU8Npl5Z/gJfdSJ58NUXUNZsHTGsZ0xgc/zZ3p1vRgYkgfSE6clG1Kb2p4KHcrX5EuMnq8TWkDNmzk7Juj3InrgeWKXzi7HEeQc0+I2Q8MM9EGx9VcFumYhgQaSj4MKcU1qkdz7mJbirm2ciGN2JKtN7so8KUXENAHQCmIfpqoBBZLiRWDH9RKlNY1zC/u+xeIrxOPqbDaKC7bV76RMd4DQlzSG5ph6B+8+4BtmYvUMqj2jAYMBSPj6lBLM2Vx4uBcQUCRBZSOAsqOS05PNGAEava6Fuggwh2Z8c9h5l2MtbP57f0AvxZ0GrXvpt69VYG3GacZNgOVGdCJO8Kc1ov7vWUId/cATsZN2RN59AUJ2MzRUdfDUU+4mS765/Bl7Oe60ouXX0v165pL8H0RnCLHoKBuutrXaX1rl09/ACDcCe5l80kH6xqKdCBwvgfDBsMFCrNszYNcEF8/68kHiBpcqT6xWyQwOfqJNFXXa9oaX002sv/0h4hCVju++uFMl0bCiRgEufE1fb5ew0n2ypfd1QZsxoiUUy25m8s4YgqEw3aX5HRqQIIi6QEWZXV/sV2C03Afq1cesOgjR/yO8s7LtOXDLY1XC0C4o697GeMdyhpAuwmcxRoEN79xoBndjqbZGJrFohhEIBh4JqAKNR0Nw7aXeoVeVw8yWV7ECret11SzmPzs30UeFYrw3qK+TdFcrgB4LfbTBNxA15Zg9zcTB/8P2HMb2/I6FTJ9Dd1jTX9/YPw3kGRQjmdsd0jRN22FNms3CstGopRSGh3JqWkhVzPUxPKoEv+nw0/1ZKYYUgirnyTMsRbxnZZgRnQuIP+BWgTUzvK7Z7JIqUCJ2qE1yBQH3+OkuDOqLyxCytKNZb7A5qed7nYEbY3bqfeqiJPz8+4LAnuZy4cLazf8lUR29WDCYp8Tg4IWSfDQ3G1cxikGzTpd0j0lGU+Fq4YEEZNV2+mviM0Divk8lDRfsdt7Se6yGGptUDIqBY/10vf74OSIPzTZebJ7XJEyEYowrJCQ81uZY8DN3SLK/TEeYTUvNsgdA9oS4YhO0rUGgwjXgJ/GiNXdiVMLmbmxWJZ0r77zE4pwzuK0w6XVcEb/z3vuQpCsWWqBY3RWQBbJfbKYHFhh5oXhLQL5WkGx/ebzXC1q3hhrrfZjgME9YKm63+OaZtBllS/eIrEKxX/D+L2kQTI6Wyl7ytxNMJ+Aq3VD7VjGHmZolRvewISfGTNJ0k7pvZEK4s99dO4iWuFzRESAn50p9GNJcJ59/gOQeO443AULI/2Ag/To9fG7PO5keaQCkEzNzg4roTVxM0bRNdZMLjHY1ldA4hj2K7vYE3DxOPPTra8ICohAM5Sg7Jq8mNjNoftalLKdMC4KwQ8NWrzEJ9hTTM8bMYRWayQ2rUF20fk9uv/VTBtLjnYgyY2g3jewpJAr6jZoYlvQ1qB+hdxjvNssA8igrKh+qW8ZP/7wTbNWkGIsqx4RnYkvDtDwbfYqd9knp0Kp0fHPA2Ol6d/DT/GVaUyYMQExfiJT8AUjLzUY4Nty6xS585he0nfPeMh1cza4A9HgkOpYdXg/opO9nDVEeK7FCgaVCP4gFqCmyx9hC848bEEFIo+WW1B8ZlZN7n21lpH6b74XlOhOUq2rIDh+5ydhGeRKUKPb7m6pvodld2T/f/nnpV3qDX8XsLv+UYYBQWAcSOSilEj3cGLRc7qzuM5NyzK4hMCU4YUupRpUGNPALYsgtDKk3+gukXgJwbtnSMi1CFdTT6R6DJHXQs+kuZ0ga39FiqEH+Ry+GiGsmvaE+ulpn+KmpBNTl+3Fp7iBSRBkwg607SRm+QwVKIn2dqvHNOK2nXe5P+Izrb7IT4+Nf2qElVPx9BLo2sYNg9K7e++TzaKjC7hsCPUeYkQV2sOVRvhbB3xzQNKPrIc9PjgD8rR4vp57dmI6m0yqYTgs9Ce8WSH+lYzhmD8pUAeEhwezqtd5vVcrE118sU9M9c+E25dSfzz8+sJTqyunKXXVqc54ZdcK+ESahF4ywcblocC/23XvHwIU/7n6FORSYIgvtZ6TqWquYfCIxdNa4OwBfrkA2X5wEMnTE+3NuqfuN9CdrfI870C1gMX7LVTAMwjF2vVMNZUZA1+5f3nSrDptXnWoMPSNU4M8jGDh6J5zGSYaSzh8DeClceBjXW5TOvtyqIyTcLoUKMJ9fMw4i7Y4Gi0FDcFZyPI+w1liwHTkcqXV7uxPONzNqalYCN1vOoMU5Zie2k9XjWhoJwlEMB6WC3r1c08dMrF125scTX4Z2ynFomMsGzr+R4GHG1GLilWemiaSAViBpHyePj3lnt6opE7a3TrKGBFAOWp1nlo6Q5FP2fikwb4/5Ir53mhxRhXu2aRbx8uIAIEse770tw78QRVNzpXlUQllMhHnL/mRUlRpFQgLDXddUHE+2GrfuJc4A48mVH/JlXDBpJmffBUMAj0juXfZGZeKI4Q7rFvUwoxvdFowTir9"
	return &sds
}

func (solver *awbUrgentCargusSolver) updateStatuses() SolverResponse {
	// Check to see if the request is already cached
	if time.Since(solver.lastUpdateCheck).Minutes() < TIME_BETWEEN_REQUEST_MIN {
		return SOLVER_CACHED
	}

	// Crawler object
	c := colly.NewCollector()

	// Find and visit all links
	var errMsg error
	statuses := []AwbUrgentCargusCheckpoint{}
	indexTraceData := 0
	indexTraceLocatie := 0
	indexTraceEvent := 0
	index := 0

	c.OnHTML("div.trace_data", func(e *colly.HTMLElement) {
		statuses = append(statuses, AwbUrgentCargusCheckpoint{})
	})


	c.OnHTML("div.trace_data", func(e *colly.HTMLElement) {
		statuses[indexTraceData].Index = index
		statuses[indexTraceData].Date = strings.Split(e.Text, " ")[0]
		statuses[indexTraceData].Time = strings.Split(e.Text, " ")[1]

		index += 1
		indexTraceData += 1
	})

	c.OnHTML("div.trace_locatie", func(e *colly.HTMLElement) {
		statuses[indexTraceLocatie].Location = e.Text

		indexTraceLocatie += 1
	})

	c.OnHTML("div.trace_event", func(e *colly.HTMLElement) {
		statuses[indexTraceEvent].Status = e.Text

		indexTraceEvent += 1
	})


	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(response *colly.Response, e error) {
		errMsg = e
	})

	// Start crawling
	c.Post(solver.url, solver.requestData)


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

func (solver *awbUrgentCargusSolver) GetStatuses() ([]string, SolverResponse) {
	responseCode := solver.updateStatuses()
	updatedStatuses := solver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "UrgentCargus: These are all the steps taken by your SameDay package")
		for _, status := range updatedStatuses {
			results = append(results, fmt.Sprintf("%s in %s at %s %s", status.Status, status.Location, status.Date, status.Time))
		}
	} else {
		results = append(results, "UrgentCargus: Could not found any records for your AWB")
	}


	return results, responseCode
}

func (solver *awbUrgentCargusSolver) GetLastStatus() ([]string, SolverResponse) {
	responseCode := solver.updateStatuses()
	updatedStatuses := solver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "UrgentCargus: Successfully found the latest status of your UrgentCargus package")
		results = append(results, fmt.Sprintf("%s in %s at %s %s", updatedStatuses[0].Status, updatedStatuses[0].Location, updatedStatuses[0].Date, updatedStatuses[0].Time))
	} else {
		results = append(results, "UrgentCargus: Could not found any records for your AWB")
	}

	return results, responseCode
}

func (solver *awbUrgentCargusSolver) GetAwb() string {
	return solver.awb
}
