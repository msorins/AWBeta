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

// V3
type WitResponseStructMap struct {
	Text string `json:"_text"`
	MsgId string `json:"msg_id"`
	Entities map[string][] WitEntity `json:"entities"`
}

// V2
type withResponseMap map[string] interface{}

// V1
type WitResponseStruct struct {
	Text string `json:"_text"`
	MsgId string `json:"msg_id"`
	Entities WitEntities `json:"entities"`
}

type WitEntities struct {
	Dhl []WitEntity `json:"dhl,omitempty"`
	FanCourier []WitEntity `json:"fancourier,omitempty"`
	Cargus []WitEntity `json:"cargus,omitempty"`
}

type WitEntity struct {
	Confidence float64 `json:"confidence"`
	Value string `json:"value"`
	Type string `json:"type"`
}


func main() {
	/*
	{"_text":"jjjjjjkjk","entities":{"dhl":[{"suggested":true,"confidence":0.57024255304067,"value":"jjjjjjkjk","type":"value"}]},"msg_id":"0bCijamf5xGrLEfdH"}
	 */

	 bytes := []byte(`{"_text":"jjjjjjkjk","entities":{"dhl":[{"suggested":true,"confidence":0.57024255304067,"value":"jjjjjjkjk","type":"value"}]},"msg_id":"0bCijamf5xGrLEfdH"}`)
	 transformWitResponse(bytes)
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
	//witResponse := map[string]interface{} {
	//	"_text": "",
	//	"msg_id": "",
	//	"entities":
	//	map[string]interface{} {
	//		"confidence": "",
	//		"value": "",
	//		"type": "",
	//	},
	//
	//}

	var witResponse WitResponseStructMap
	json.Unmarshal(bodyBytes, &witResponse)

	fmt.Println(witResponse)

	processMessageType(witResponse)
}


func processMessageType(data WitResponseStructMap) {
	//len(data.Entities)
	fmt.Println(data)


	//v := reflect.ValueOf(data)
	//
	//for i := 0; i < v.NumField(); i++ {
	//	fmt.Printf("name: %+v, value: %+v (%T)\n",
	//		v.Type().Field(i).Name, // Name attribute gives us the struct's value
	//		v.Field(i).Elem(), 	// Elem() dereferences the pointer value
	//		v.Field(i).Interface()) // Interface() provides memory address of the value
	//}
	//fmt.Println(a)
	//
	//data.Entitie
}

