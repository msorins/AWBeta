package solvers

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"sort"
)

type awbDhlSolver struct {
	awb string
	url string
	Statuses []IPackageStatus
}

type AWbDhlResponse struct {
	Results []AwbDhlCheckpointHolder `json:"results"`
}
type AwbDhlCheckpointHolder struct {
	Checkpoints []AwbDhlCheckpoint `json:"checkpoints"`
}

type AwbDhlCheckpoint struct {
	Index int `json:"counter"`
	Status string `json:"description"`
	Date string `json:"date"`
	Time string `json:"time"`
	Location string `json:"location"`
}

func AwbDhlSolverBuilder(awb string) ISolver{
	awbSolver := awbDhlSolver{}
	awbSolver.url = "https://www.dhl.ro/shipmentTracking?AWB="
	awbSolver.awb = awb
	return awbSolver
}

func (solver awbDhlSolver) GetStatusesForAwb() []IPackageStatus {
	var urlToSend string
	urlToSend = solver.url + url.QueryEscape(solver.awb)

	respAwb, _ := http.Get(urlToSend)

	// Get the response
	if respAwb.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(respAwb.Body)

		// Transform it to a struct
		rs := transformDhlSolverRequest(bodyBytes)

		// Assign it to class member
		solver.Statuses = []IPackageStatus{}
		for _, value := range  rs.Results[0].Checkpoints {
			var crtPackageStatus IPackageStatus
			crtPackageStatus.Index = value.Index
			crtPackageStatus.DateTime = value.Date + value.Time
			crtPackageStatus.Location = value.Location
			crtPackageStatus.Status = value.Status

			solver.Statuses = append(solver.Statuses, crtPackageStatus)
		}

		sort.Slice(solver.Statuses, func(i, j int) bool {
			return solver.Statuses[i].Index > solver.Statuses[j].Index
		})
	} else {
		fmt.Println("Error in request")
	}

	return solver.Statuses
}

func (awbsolver awbDhlSolver) GetStatuses() []string {
	updatedStatuses := awbsolver.GetStatusesForAwb()
	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "These are all the steps taken by your DHL package")
		for _, status := range updatedStatuses {
			results = append(results, fmt.Sprintf("%s %s", status.Status, status.DateTime))
		}
	} else {
		results = append(results, "Could not found any records for your AWB")
	}


	return results
}

func (awbsolver awbDhlSolver) GetLastStatus() []string{
	updatedStatuses := awbsolver.GetStatusesForAwb()
	results := []string{}

	if len(updatedStatuses) >= 1 {
		results = append(results, "Successfully found the latest status of your DHL package")
		results = append(results, fmt.Sprintf("%s, %s", updatedStatuses[0].Status, updatedStatuses[0].DateTime))
	} else {
		results = append(results, "Could not found any records for your AWB")
	}


	return results
}

func transformDhlSolverRequest(bodyBytes []byte) AWbDhlResponse {
	var awbResponse AWbDhlResponse
	json.Unmarshal(bodyBytes, &awbResponse)

	return awbResponse
}

