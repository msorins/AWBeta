package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"github.com/paked/messenger"
	"time"
	"net/url"
	"net/http"
	"log"
	"os"
	"wit"
	"solvers"
	"io/ioutil"
	"state"
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

var resolverMap = map[string]func(string,  map[string][]wit.WitEntity) solvers.ISolver {
	"dhl" : solvers.AwbDhlSolverBuilder,
	"fanCourier": solvers.AwbFanCourierSolverBuilder,
	"unknown": solvers.UnknownFanCourierSolverBuilder,
}

func main() {
	// Real server
	//messengerServer()

	// Mock messages
	st := state.StateManagerBuilder()
	messageMock := messenger.Message{}
	messageMock.Sender.ID = 123456
	messageMock.Text = "Hi, what's the status for 2032810250356"

	fmt.Println( messageHandleToRes(st, messageMock) )
}

func messengerServer(stateManager *state.StateManager) {

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

		// Get the results for the message received
		results := messageHandleToRes(stateManager, m)

		// Send them to the user
		for _, str := range results {
			r.Text(str, messenger.ResponseType)
		}
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

func messageHandleToRes(stateManager *state.StateManager, message messenger.Message) []string {
	// Get the message text & form WIT request
	var urlToSend string
	urlToSend = "https://api.wit.ai/message?v=20180617&q=" + url.QueryEscape(message.Text)

	clientWit := &http.Client{}
	reqWit, _ := http.NewRequest("GET", urlToSend, nil)
	reqWit.Header.Add("Authorization", "Bearer XSNNOAK5JCAEYUULJ6V6YJ6G45VSJ6TV")
	respWit, _ := clientWit.Do(reqWit)

	// Get wit response, if it is ok send it further to be parsed & get a result
	if respWit.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(respWit.Body)

		// Transform byte array into an response
		var sentToUSer []string
		sentToUSer = witToRes(stateManager, fmt.Sprintf("%v", message.Sender.ID), bodyBytes)

		// Return the result ( a list of strings )
		return sentToUSer
	}

	return []string{}
}

// Here must implement the flow
// TO DO: Pass by reference
func witToRes(stateManager *state.StateManager, userId string, bodyBytes []byte) []string {
	// Transform byte array into an response
	rw := transformWitResponse(bodyBytes)

	if stateManager.IdExists(userId) {

	} else {
		// Get the handler needed to process
		handler := processMessageType(rw)
		res, responseCode := handler.GetLastStatus()

		// Update the stateManager
		switch responseCode {
			case solvers.SOLVER_OK:
				stateManager.SetState(userId, handler, "REQUESTED_AWB_STATUS")
			case solvers.SOLVER_AWB_INCORRECT:
				res = append(res, "Could you please specify a courier name?")
				stateManager.SetState(userId, handler, "REQUESTED_PROVIDER_NAME")

		}

		return res
	}

	return []string{}
}

func transformWitResponse(bodyBytes []byte) wit.WitResponseStructMap {
	// Transform byte array ti WitResponseStructMap
	var witResponse wit.WitResponseStructMap
	json.Unmarshal(bodyBytes, &witResponse)

	return witResponse
}

func processMessageType(data wit.WitResponseStructMap) solvers.ISolver {
	// Get the courier intent with the biggest probability
	var bestEntityCourierName string
	bestEntityCourierName = "unknown"
	bestEntity := wit.WitEntity{}
	bestEntity.Confidence = -1


	for key, value := range data.Entities {
		if value[0].Confidence > bestEntity.Confidence{
			_, ok := resolverMap[ key ]
			if ok == true {
				bestEntity = value[0]
				bestEntityCourierName = key
			}

		}
	}

	// Call the resolver for the given awb & courier firm
	return resolverMap[bestEntityCourierName](bestEntity.Value, data.Entities)
}
