package solvers

import (
"net/http"
"io/ioutil"
"fmt"
"encoding/json"
	"net/url"
	"strings"
	"strconv"
	"sort"
	"time"
)

type awbFanCourierSolver struct {
	awb                string
	url                string
	Statuses           []AwbFanCourierCheckpoint
	LastSolverResponse SolverResponse
	lastUpdateCheck	time.Time
}

type AWbFanCourierResponse struct {
	Entities map[string]AwbFanCourierCheckpoint `json:"1"`
}

type AwbFanCourierCheckpoint struct {
		Index int `json:"nstex"`
		Status string `json:"mstex"`
		Date string `json:"dstex"`
}

func AwbFanCourierSolverBuilder(awb string) ISolver{
	awbSolver := awbFanCourierSolver{}
	awbSolver.url = "https://www.fancourier.ro/wp-content/themes/fancourier/webservice.php"
	awbSolver.awb = awb
	return &awbSolver
}

func (solver *awbFanCourierSolver) updateStatuses() SolverResponse {
	// Check to see if the request is already cached
	if time.Since(solver.lastUpdateCheck).Minutes() < TIME_BETWEEN_REQUEST_MIN {
		return SOLVER_CACHED
	}

	var urlToSend string
	urlToSend = solver.url

	form := url.Values{}
	form.Add("awb", solver.awb)
	form.Add("metoda", "tracking")

	hc := http.Client{}
	req, _ := http.NewRequest("POST", urlToSend, strings.NewReader(form.Encode()) )
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	respAwb, _ := hc.Do(req)

	// Get the response
	if respAwb.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(respAwb.Body)

		// Transform it to a struct
		rs := transformFanCourierSolverRequest(bodyBytes)

		if len(rs.Entities) <= 3 {
			solver.LastSolverResponse = SOLVER_AWB_INCORRECT
			return SOLVER_AWB_INCORRECT
		}

		// Assign it to class member
		lst := []AwbFanCourierCheckpoint{}

		for _, value := range  rs.Entities {
			if value.Index != 0 {
				value.Status = value.Status[: len(value.Status) - 9]
				lst = append(lst, value)
			}

		}

		sort.Slice(lst, func(i, j int) bool {
			return lst[i].Index > lst[j].Index
		})

		solver.Statuses = lst
	} else {
		solver.LastSolverResponse = SOLVER_BAD_REQUEST
		return SOLVER_BAD_REQUEST
	}

	solver.lastUpdateCheck = time.Now()
	solver.LastSolverResponse = SOLVER_OK
	return SOLVER_OK
}

func (awbsolver *awbFanCourierSolver) GetStatuses() ([]string, SolverResponse) {
	responseCode := awbsolver.updateStatuses()
	updatedStatuses := awbsolver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "FAN: These are all the steps taken by your FanCourier package")
		for _, status := range updatedStatuses {
			results = append(results, fmt.Sprintf("%s, %s", status.Status, status.Date))
		}
	} else {
		results = append(results, "FAN: Could not found any records for your AWB")
	}

	return results, responseCode
}

func (awbsolver *awbFanCourierSolver) GetLastStatus() ([]string, SolverResponse)  {
	responseCode := awbsolver.updateStatuses()
	updatedStatuses := awbsolver.Statuses

	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "FAN: Successfully found the latest status of your FanCourier package")
		results = append(results, fmt.Sprintf("%s, %s", updatedStatuses[0].Status, updatedStatuses[0].Date))
	} else {
		results = append(results, "FAN: Could not found any records for your AWB")
	}

	return results, responseCode
}

func (awbsolver *awbFanCourierSolver) GetLastSolverResponse() SolverResponse {
	return awbsolver.LastSolverResponse
}

func transformFanCourierSolverRequest(bodyBytes []byte) AWbFanCourierResponse {
	var awbResponse AWbFanCourierResponse
	json.Unmarshal(bodyBytes, &awbResponse)

	return awbResponse
}

func (awbsolver *awbFanCourierSolver) GetAwb() string {
	return awbsolver.awb
}
