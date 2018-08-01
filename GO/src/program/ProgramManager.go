package program

import (
	"state"
	"subscription"
	"chat"
	"sync"
	"wit"
	"encoding/json"
	"solvers"
	"time"
	"log"
	"net/url"
	"net/http"
	"io/ioutil"
	"fmt"
)

type ProgramManager struct {
	stateManager state.StateManager
	subscriptionManager subscription.SubscriptionManager
	chatManager chat.IChat
}

// To do: change this
var resolverMap = map[string]func(string) solvers.ISolver {
	"dhl" : solvers.AwbDhlSolverBuilder,
	"fanCourier": solvers.AwbFanCourierSolverBuilder,
	"sameDay": solvers.SameDaySolverBuilder,
	"urgentCargus": solvers.UrgentCargusBuilder,
	"dpd": solvers.DpdSolverBuile,
	"unknown": solvers.UnknownFanCourierSolverBuilder,
}

func ProgramManagerBuilder() ProgramManager{
	programManager := ProgramManager{}

	// Create the state object
	programManager.stateManager = state.StateManagerBuilder()

	// Create the subscription manager object
	programManager.subscriptionManager = subscription.SubscriptionManagerBuilder()

	// Create the chat manager
	programManager.chatManager = chat.FacebookMessengerBuilder()

	return programManager
}


func (programManager *ProgramManager) StartProd() {
	// Start the go routines
	var wg sync.WaitGroup
	wg.Add(2)

	// Start Subscription listening
	go programManager.startSubscriptionListening(&programManager.subscriptionManager, programManager.chatManager)

	// Start Facebook Messenger Server
	go programManager.startFacebookMessengerServer(&programManager.stateManager, programManager.chatManager)

	wg.Wait()
}


func (programManager *ProgramManager) startFacebookMessengerServer(stateManager *state.StateManager, chatManager chat.IChat) {
	chatManager.HandleMessages(programManager.MessageHandle)
}

func (programManager *ProgramManager) startSubscriptionListening(subscriptionManager *subscription.SubscriptionManager, chat chat.IChat) {
	for range time.Tick(time.Duration(time.Second)) {
		changes, err := subscriptionManager.CheckForChanges()

		if err != nil {
			log.Fatal("Error in startSubscriptionListening")
		}
		// Send the messages
		for key, value := range changes {
			if len(value) != 0 {
				chat.SendMessage(key, value)
			}
		}
	}
}


// To change here the parameters
func (programManager *ProgramManager) MessageHandle( senderId string, message string) []string {
	// Get the message text & form WIT request
	var urlToSend string
	urlToSend = "https://api.wit.ai/message?v=20180617&q=" + url.QueryEscape(message)

	clientWit := &http.Client{}
	reqWit, _ := http.NewRequest("GET", urlToSend, nil)
	reqWit.Header.Add("Authorization", "Bearer XSNNOAK5JCAEYUULJ6V6YJ6G45VSJ6TV")
	respWit, _ := clientWit.Do(reqWit)

	// Get wit response, if it is ok send it further to be parsed & get a result
	if respWit.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(respWit.Body)

		// Transform byte array into an response
		var sentToUSer []string
		sentToUSer = programManager.processMessageByState(fmt.Sprintf("%v", senderId), bodyBytes)

		// Return the result ( a list of strings )
		return sentToUSer
	} else {
		return []string{ "Something fishy has happened" }
	}

	return []string{}
}

// Here must implement the flow
func (programManager *ProgramManager) processMessageByState(userId string, bodyBytes []byte) []string {
	// Transform byte array into an response
	rw := programManager.transformWitResponse(bodyBytes)

	// If the user has already a state associated with him / her => deal with it accordingly
	if programManager.stateManager.IdExists(userId) {
		stateOfUser, _ := programManager.stateManager.GetState(userId)

		switch stateOfUser.State {
		case state.USER_STATE_AWB_CONFUSING:
			handler := programManager.getHandlerFromName(stateOfUser, rw)
			res, responseCode := handler.GetLastStatus()

			switch responseCode {
			// Operation completed successfully -> delete the state
			case solvers.SOLVER_OK:
				programManager.stateManager.SetState(userId, handler, state.USER_STATE_AWB_OK)

				// Provided awb was incorect -> ask him to specify the name of the awb
			case solvers.SOLVER_AWB_INCORRECT:
				res = append(res, "Please try again with other awb :)")
				programManager.stateManager.RemoveState(userId)
			}

			return res;

		case state.USER_STATE_AWB_OK:
			messageIntent := programManager.getMessageIntent(rw)

			switch messageIntent {
			case wit.MESSAGE_NO_INTENT:
				res, _ := stateOfUser.Solver.GetLastStatus()
				return res
			case wit.MESSAGE_REQUEST_ALL_HISTORY:
				res, _ := stateOfUser.Solver.GetStatuses()
				return res
			case wit.MESSAGE_REQUEST_SUBSCRIPTION:
				statuses, _ := stateOfUser.Solver.GetStatuses()
				programManager.subscriptionManager.AddSubscription(stateOfUser.Solver.GetAwb(), subscription.SubscriptionManagerEntity{stateOfUser.Solver, len(statuses), userId})
			case wit.MESSAGE_REQUEST_NEW_AWB:
				// Remove the state of the old awb && recall the function
				programManager.stateManager.RemoveState(userId)
				return programManager.processMessageByState(userId, bodyBytes)
			}
		}
	} else { // User has no state associated -> check the message for an awb
		// Get the handler needed to process
		handler := programManager.getHandlerFromAwb(rw)
		res, responseCode := handler.GetLastStatus()

		// Update the stateManager
		switch responseCode {
		// Operation completed successfully -> delete the state
		case solvers.SOLVER_OK:
			programManager.stateManager.SetState(userId, handler, state.USER_STATE_AWB_OK)

			// Provided awb was incorect -> ask him to specify the name of the awb
		case solvers.SOLVER_AWB_INCORRECT:
			res = append(res, "Could you please specify a courier name?")
			programManager.stateManager.SetState(userId, handler, state.USER_STATE_AWB_CONFUSING)
		}

		return res
	}

	return []string{}
}

func (programManager *ProgramManager) transformWitResponse(bodyBytes []byte) wit.WitResponseStructMap {
	// Transform byte array ti WitResponseStructMap
	var witResponse wit.WitResponseStructMap
	json.Unmarshal(bodyBytes, &witResponse)

	return witResponse
}

func (programManager *ProgramManager) getHandlerFromAwb(data wit.WitResponseStructMap) solvers.ISolver {
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
	return resolverMap[bestEntityCourierName](bestEntity.Value)
}

func (programManager *ProgramManager) getHandlerFromName(stateOfRequester state.StateManagerStruct, data wit.WitResponseStructMap) solvers.ISolver {
	var bestEntityCourierName string
	bestEntityCourierName = "unknown"

	companyName, exists := data.Entities["companyName"]
	if exists == true {
		bestEntityCourierName = companyName[0].Value
	}

	// Call the resolver for the given awb & courier firm
	return resolverMap[bestEntityCourierName](stateOfRequester.Solver.GetAwb())
}

func (programManager *ProgramManager) getMessageIntent(data wit.WitResponseStructMap) wit.MessageIntent {
	intent := wit.MESSAGE_NO_INTENT

	intentEntity, exists := data.Entities["Intent"]
	if exists == true {
		switch intentEntity[0].Value {
		case "REQUEST_ALL_HISTORY":
			intent = wit.MESSAGE_REQUEST_ALL_HISTORY

		case "REQUEST_SUBSCRIPTION":
			intent = wit.MESSAGE_REQUEST_SUBSCRIPTION

		case "REQUEST_NEW_AWB":
			intent = wit.MESSAGE_REQUEST_NEW_AWB
		}
	}

	return intent
}