package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"github.com/paked/messenger"
	"time"
	"net/url"
	"net/http"
	"io/ioutil"
	"log"
	"os"
	"wit"
	"solvers"
)


var (
	verifyToken = flag.String("verify-token", "soarecelmaifainb@T", "The token used to verify facebook (required)")
	verify      = flag.Bool("should-verify", false, "Whether or not the app should verify itself")
	pageToken   = flag.String("page-token", "EAAcJw2oDsswBAPfJfJEXYC96SRHOAV37UmoPWVQ8ssaidzLdUPmYSOy1eGp7wEmJZC6MdiU10SuU5ptVE784YrsF092PmuUzPEmolR5pxYZAUaEH6PNL8hwRJKWBHjhRBDl9L6D2WyE6eJkBcY0buocNjuZAGD9n9fcopREFjiSR4qWeXFU", "The token that is used to verify the page on facebook")
	appSecret   = flag.String("app-secret", "596b7437a204b6aaff57b4e72938afec", "The app secret from the facebook developer portal (required)")
	host        = flag.String("host", "localhost", "The host used to serve the messenger bot")
	port        = flag.Int("port", 3000, "The port used to serve the messenger bot")
	witToken        = flag.String("witToken", "XSNNOAK5JCAEYUULJ6V6YJ6G45VSJ6TV", "Token for wit.ai")
	couriers = []string{"Dhl", "FanCourier", "Cargus"}
)


func main() {
	//bytes := []byte(`{"_text":"jjjjjjkjk","entities":{"dhl":[{"suggested":true,"confidence":0.57024255304067,"value":"jjjjjjkjk","type":"value"}]},"msg_id":"0bCijamf5xGrLEfdH"}`)
	//transformWitResponse(bytes)

	solver := AwbDhlSolverBuilder("1627190725")
	fmt.Println( solver.GetStatusesForAwb()[0] )
}

func messengerServer() {

	flag.Parse()

	if *verifyToken == "" || *appSecret == "" || *pageToken == "" {
		fmt.Println("missing arguments")
		fmt.Println()
		flag.Usage()

		os.Exit(-1)
	}

	// Create a new messenger client
	client := messenger.New(messenger.Options{
		Verify:      *verify,
		AppSecret:   *appSecret,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})


	// Setup a handler to be triggered when a message is received
	client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		var urlToSend string
		urlToSend = "https://api.wit.ai/message?v=20180617&q=" + url.QueryEscape(m.Text)

		clientWit := &http.Client{}
		reqWit, _ := http.NewRequest("GET", urlToSend, nil)
		reqWit.Header.Add("Authorization", "Bearer XSNNOAK5JCAEYUULJ6V6YJ6G45VSJ6TV")
		respWit, err := clientWit.Do(reqWit)
		if respWit.StatusCode == http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(respWit.Body)
			transformWitResponse(bodyBytes)
		}
		fmt.Println(respWit)

		p, err := client.ProfileByID(m.Sender.ID)
		if err != nil {
			fmt.Println("Something went wrong!", err)
		}

		r.Text(fmt.Sprintf("Hello, %v!", p.FirstName), messenger.ResponseType)
	})


	// Setup a handler to be triggered when a message is delivered
	client.HandleDelivery(func(d messenger.Delivery, r *messenger.Response) {
		fmt.Println("Delivered at:", d.Watermark().Format(time.UnixDate))
	})

	// Setup a handler to be triggered when a message is read
	client.HandleRead(func(m messenger.Read, r *messenger.Response) {
		fmt.Println("Read at:", m.Watermark().Format(time.UnixDate))
	})

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Serving messenger bot on", addr)
	log.Fatal(http.ListenAndServe(addr, client.Handler()))
}

func transformWitResponse(bodyBytes []byte) {
	// Transform byte array ti WitResponseStructMap
	var witResponse wit.WitResponseStructMap
	json.Unmarshal(bodyBytes, &witResponse)

	// Process the message based on the entities provided by the AI
	processMessageType(witResponse)
}

func processMessageType(data wit.WitResponseStructMap) {
	// Get the courier intent with the biggest probability
	var bestEntityCourierName string
	bestEntity := wit.WitEntity{}
	bestEntity.Confidence = -1


	for key, value := range data.Entities {
		if value[0].Confidence > bestEntity.Confidence{
			bestEntity = value[0]
			bestEntityCourierName = key
		}
		fmt.Printf("%s   ->  v%s \n", key, value)
	}

	// Call the resolver for the given awb & courier firm
	fmt.Println(bestEntity)
	fmt.Println(bestEntityCourierName)
}



type awbDhlSolver struct {
	awb string
	url string
	Statuses []solvers.IPackageStatus
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

func AwbDhlSolverBuilder(awb string) awbDhlSolver {
	awbSolver := awbDhlSolver{}
	awbSolver.url = "https://www.dhl.ro/shipmentTracking?AWB="
	awbSolver.awb = awb
	return awbSolver
}

func (solver *awbDhlSolver) GetStatusesForAwb() []solvers.IPackageStatus {
	var urlToSend string
	urlToSend = solver.url + url.QueryEscape(solver.awb)

	respAwb, _ := http.Get(urlToSend)

	// Get the response
	if respAwb.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(respAwb.Body)

		// Transform it to a struct
		rs := transformAwbSolverRequest(bodyBytes)

		// Assign it to class member
		solver.Statuses = []solvers.IPackageStatus{}
		for _, value := range  rs.Results[0].Checkpoints {
			var crtPackageStatus solvers.IPackageStatus
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

func transformAwbSolverRequest(bodyBytes []byte) AWbDhlResponse {
	var awbResponse AWbDhlResponse
	json.Unmarshal(bodyBytes, &awbResponse)

	return awbResponse
}

func (awbsolver *awbDhlSolver) GetStatuses() []solvers.IPackageStatus {
	return awbsolver.Statuses
}

func (awbsolver *awbDhlSolver) GetLastStatu() solvers.IPackageStatus{
	return awbsolver.Statuses[ len(awbsolver.Statuses) - 1 ]
}