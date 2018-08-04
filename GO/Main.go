package main

import (
	"program"
)

func main() {
	// Instantiate the program
	programManager := program.ProgramManagerBuilder()

	// Send some mock message
	//messageMock := messenger.Message{}
	//messageMock.Sender.ID = 123456

	// CASE 1 -> SPECIFING THE RIGHT COURIER FIRM
	//messageMock.Text = "Hi, what's the status for 2032810250356"
	//fmt.Println( programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text) )

	//messageMock.Text = "Please tell me all its history"
	//fmt.Println( programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text) )

	// Start the production code listeners
	programManager.StartProd()
}

/*
//
	//messageMock.Text = "DHL"
	//fmt.Println( messageHandleToRes(&stateManager, messageMock) )
	//
	//
	// CASE 2 -> ALRIGHT AWB
	//messageMock.Text = "Hi, what's the status for 2032810250356"
	//fmt.Println( messageHandleToRes(&stateManager, messageMock) )
	//
	// CASE 3 -> ALRIGHT AWB -> Request all history for that awb
	//messageMock.Text = "Hi, what's the status for 2032810250356"
	//fmt.Println( messageHandleToRes(&stateManager, messageMock) )
	//
	//messageMock.Text = "Please show me all the statistics"
	//fmt.Println( messageHandleToRes(&stateManager, messageMock) )

 */

