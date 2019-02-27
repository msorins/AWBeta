package chat

import (
	"github.com/paked/messenger"
	"flag"
	"fmt"
	"time"
	"strconv"
	"log"
	"net/http"
)

type facebookMessenger struct {
	client *messenger.Messenger
}

var (
	verifyToken = flag.String("verify-token", "soarecelmaifainb@T", "The token used to verify facebook (required)")
	verify      = flag.Bool("should-verify", false, "Whether or not the app should verify itself")
	pageToken   = flag.String("page-token", "EAAcJw2oDsswBAARloVZAVfmrIJWIuZAXHsxfaQhWJHxdMbgHJ18sjrcGvgWQONDWcPnoCsNT5dJVkXOYw75LHMHk8OwjGsuthwauk27PTEOc9kFaZBD5VvlyTQZCRJZBYCDhYj6AltdQaYzJL4bDKsCEDz3ZBsBG653HFZCwL1w4AZDZD", "The token that is used to verify the page on facebook")
	appSecret   = flag.String("app-secret", "596b7437a204b6aaff57b4e72938afec", "The app secret from the facebook developer portal (required)")
	host        = flag.String("host", "localhost", "The host used to serve the messenger bot")
	port        = flag.Int("port", 3000, "The port used to serve the messenger bot")
	witToken        = flag.String("witToken", "XSNNOAK5JCAEYUULJ6V6YJ6G45VSJ6TV", "Token for wit.ai")
	couriers = []string{"Dhl", "FanCourier", "Cargus"}
)


func FacebookMessengerBuilder() IChat {
	client := messenger.New(messenger.Options{
		Verify:      *verify,
		AppSecret:   *appSecret,
		VerifyToken: *verifyToken,
		Token:       *pageToken,
	})

	fbm := facebookMessenger{client}

	return &fbm
}

func (fb *facebookMessenger) HandleMessages(messageReceivedCallBack func(string, string) []string) {

	fb.client.HandleMessage(func(m messenger.Message, r *messenger.Response) {
		fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

		// Get the results for the message received
		results := messageReceivedCallBack(fmt.Sprintf("%v", m.Sender.ID), m.Text)

		// Send them to the user
		for _, str := range results {
			if str == "<sendExtraAwbOptions>" {
				p, _ := fb.client.ProfileByID(m.Sender.ID)
				sendExtraAwbOptions(p, r)
				continue
			}
			r.Text(str, messenger.ResponseType)
		}
	})
	
	fb.client.HandleAccountLinking(func(linking messenger.AccountLinking, response *messenger.Response) {
		response.Text("Hello, ce mai faaaci", messenger.ResponseType)

	})

	fb.client.HandlePostBack(func(back messenger.PostBack, response *messenger.Response) {
		response.Text("Hello you idiot", messenger.ResponseType)
	})

	addr := fmt.Sprintf("%s:%d", *host, *port)
	log.Println("Serving messenger bot on", addr)
	log.Fatal(http.ListenAndServe(addr, fb.client.Handler()))
}

func sendExtraAwbOptions(p messenger.Profile, r *messenger.Response) error {
	text := fmt.Sprintf(
		"%s, want to know even more about your package?",
		p.FirstName,
	)

	replies := []messenger.QuickReply{
		{
			ContentType: "text",
			Title:       "Past statuses",
		},
		{
			ContentType: "text",
			Title:       "Subscribe to changes",
		},
	}

	return r.TextWithReplies(text, replies, messenger.ResponseType)
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