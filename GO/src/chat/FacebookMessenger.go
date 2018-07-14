package chat

import (
	"github.com/paked/messenger"
	"flag"
	"fmt"
	"time"
	"state"
	"strconv"
	"log"
	"net/http"
)

type facebookMessenger struct {
	client *messenger.Messenger
	stateManager *state.StateManager
}

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


func FacebookMessengerBuilder(stateManager *state.StateManager) IChat {
	client := messenger.New(messenger.Options{
		Verify:      *verify,
		AppSecret:   *appSecret,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})

	fbm := facebookMessenger{client, stateManager}

	return &fbm
}

func (fb *facebookMessenger) HandleMessages(messageReceivedCallBack func(*state.StateManager, string, string) []string) {

	fb.client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		// Get the results for the message received
		results := messageReceivedCallBack(fb.stateManager, fmt.Sprintf("%v", m.Sender.ID), m.Text)

		// Send them to the user
		for _, str := range results {
			r.Text(str, messenger.ResponseType)
		}
	})

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Serving messenger bot on", addr)
	log.Fatal(http.ListenAndServe(addr, fb.client.Handler()))
}

func (fb *facebookMessenger) SendMessage(userId string, msgs []string) {
	// Send the messages to the recipient
	var id int64
	id, _ = strconv.ParseInt(userId, 10, 64)

	recipient := messenger.Recipient{id}
	for _, str := range msgs {
		fb.client.Send(recipient, str, messenger.UpdateType)
	}
}