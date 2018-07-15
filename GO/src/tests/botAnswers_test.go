package tests

import (
	"testing"
	"program"
	"strconv"
	"github.com/paked/messenger"
	"strings"
)

func TestFanCourierLastStatus(t *testing.T) {
	// Instantiate the program
	programManager := program.ProgramManagerBuilder()

	// Send message for Fan Courier
	messageMock := messenger.Message{}
	messageMock.Sender.ID = 123456

	texts := []string {"Hi, what's the status for 2032810250356"}

	for _, text := range texts {
		messageMock.Text = text
		res := programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

		// Check to see if response is correct
		if len(res) != 2  || res[0] != "FAN: Successfully found the latest status of your FanCourier package" || res[1] != "Ultimul status al expeditiei: livrat, 03.02.2018 10:36" {
			t.Error("TestFanCourierLastStatus wrong")
		}
	}

}

func TestFanCourierPastStatuses(t *testing.T) {
	// Instantiate the program
	programManager := program.ProgramManagerBuilder()

	// Send message for Fan Courier
	messageMock := messenger.Message{}
	messageMock.Sender.ID = 123456

	messageMock.Text = "Hi, what's the status for 2032810250356"
	programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

	texts := []string {"Could you please tell all my past statuses?"}
	for _, text := range texts {
		messageMock.Text = text
		res := programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

		// Check to see if response is correct
		if len(res) != 5 {
			t.Error("TestFanCourierLastStatus wrong")
		}
	}
}

func TestDHLLastStatus(t *testing.T) {
	// Instantiate the program
	programManager := program.ProgramManagerBuilder()

	// Send message for Fan Courier
	messageMock := messenger.Message{}
	messageMock.Sender.ID = 123456
	texts := []string {"Hi, what's the status for 1627190725", "Could you please tell my where my collet is: 1627190725 ?", "Where is 1627190725 ?? Does it take much longer"}

	for _, text := range texts {
		messageMock.Text = text
		res := programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

		// Check to see if response is correct
		if len(res) != 2  || res[0] != "DHL: Successfully found the latest status of your DHL package" || res[1] != "Delivered - Signed for by: mircea sorin sebastianb, Friday, May 25, 2018  11:54" {
			t.Error("TestFanCourierLastStatus wrong")
		}
	}
}

func TestDHLPastStatuses(t *testing.T) {
	// Instantiate the program
	programManager := program.ProgramManagerBuilder()

	// Send message for Fan Courier
	messageMock := messenger.Message{}
	messageMock.Sender.ID = 123456

	messageMock.Text = "Hi, what's the status for 1627190725"
	programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

	texts := []string {"Could you please tell all my past statuses?", "I want the full history, please"}
	for _, text := range texts {
		messageMock.Text = text
		res := programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

		// Check to see if response is correct
		if len(res) != 19 {
			t.Error("TestFanCourierLastStatus wrong")
		}
	}
}

func TestWrongAwb(t *testing.T) {
	// Instantiate the program
	programManager := program.ProgramManagerBuilder()

	// Send message for Fan Courier
	messageMock := messenger.Message{}
	messageMock.Sender.ID = 123456
	texts := []string {"Hi, what's the status for 4327190725"}

	for _, text := range texts {
		messageMock.Text = text
		res := programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

		// Check to see if response is correct
		if len(res) != 2 || !strings.Contains(res[0], "Could not found any records") {
			t.Error("TestFanCourierLastStatus wrong")
		}
	}
}

func TestNoAwb(t *testing.T) {
	// Instantiate the program
	programManager := program.ProgramManagerBuilder()

	// Send message for Fan Courier
	messageMock := messenger.Message{}
	messageMock.Sender.ID = 123456
	texts := []string {"Hi, are you thereeee?", "Hello", "Thank you"}

	for _, text := range texts {
		messageMock.Text = text
		res := programManager.MessageHandle(strconv.FormatInt(messageMock.Sender.ID, 10), messageMock.Text)

		// Check to see if response is correct
		if len(res) != 2 || !strings.Contains(res[0], "Could not found") {
			t.Error("TestFanCourierLastStatus wrong")
		}
	}
}
