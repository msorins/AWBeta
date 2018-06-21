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
)

type awbFanCourierSolver struct {
	awb string
	url string
	Statuses []IPackageStatus
}

type AWbFanCourierResponse struct {
	Entities map[string] AWbFanCourierCheckpoint `json:"1"`
}

//type Aux struct {
//	FirstEntity  `json:"0"`
//}
//type AWbFanCourierCheckpoints struct {
//	Entities map[string][] AWbFanCourierCheckpoint
//}

type AWbFanCourierCheckpoint struct {
		Index int `json:"nstex"`
		Status string `json:"mstex"`
		Date string `json:"dstex"`
}
//type AWbFanCourierResponse struct {
//	Results []AwbFanCourierCheckpoint `json:"1"`
//}
//
//type AwbFanCourierCheckpoint struct {
//	Status string `json:"mstex"`
//	Date string `json:"dstex"`
//}

func AwbFanCourierSolverBuilder(awb string) ISolver{
	awbSolver := awbFanCourierSolver{}
	awbSolver.url = "https://www.fancourier.ro/wp-content/themes/fancourier/webservice.php"
	awbSolver.awb = awb
	return awbSolver
}

func (solver awbFanCourierSolver) GetStatusesForAwb() []IPackageStatus {
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
		fmt.Println(rs)

		for key, value := range rs.Entities {
			fmt.Printf("%s -> %s\n", key, value)
		}

		// Assign it to class member
		solver.Statuses = []IPackageStatus{}
		lst := []IPackageStatus{}

		for _, value := range  rs.Entities {
			if value.Index != 0 {
				var crtPackage IPackageStatus
				crtPackage.Index = value.Index
				crtPackage.Status = value.Status[: len(value.Status) - 9]
				crtPackage.DateTime = value.Date
				crtPackage.Location = ""

				lst = append(lst, crtPackage)
			}

		}

		sort.Slice(lst, func(i, j int) bool {
			return lst[i].Index > lst[j].Index
		})

		solver.Statuses = lst
		fmt.Println(lst)
	} else {
		fmt.Println("Error in request")
	}

	return solver.Statuses
}

func (awbsolver awbFanCourierSolver) GetStatuses() []IPackageStatus {
	return awbsolver.Statuses
}

func (awbsolver awbFanCourierSolver) GetLastStatus() IPackageStatus{
	return awbsolver.Statuses[ len(awbsolver.Statuses) - 1 ]
}

func transformFanCourierSolverRequest(bodyBytes []byte) AWbFanCourierResponse {
	var awbResponse AWbFanCourierResponse
	json.Unmarshal(bodyBytes, &awbResponse)

	return awbResponse
}

