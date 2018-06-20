package solvers

import (
	"net/url"
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
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
			crtPackageStatus.DateTime = value.Date + value.Time
			crtPackageStatus.Location = value.Location
			crtPackageStatus.Status = value.Status

			solver.Statuses = append(solver.Statuses, crtPackageStatus)
		}
	} else {
		fmt.Println("Error in request")
	}

	return solver.Statuses
}

func (awbsolver awbDhlSolver) GetStatuses() []IPackageStatus {
	return awbsolver.Statuses
}

func (awbsolver awbDhlSolver) GetLastStatus() IPackageStatus{
	return awbsolver.Statuses[ len(awbsolver.Statuses) - 1 ]
}

func transformDhlSolverRequest(bodyBytes []byte) AWbDhlResponse {
	var awbResponse AWbDhlResponse
	json.Unmarshal(bodyBytes, &awbResponse)

	return awbResponse
}

