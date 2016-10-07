package rabbit

import (
	"encoding/json"
	"net/http"

	"github.com/AKSW/micropipe-proxy/config"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Payload to be encoded for message bus
type Payload struct {
	Data   interface{} `json:"data"`
	Config interface{} `json:"config"`
}

// InitResponseServer inits response server
func InitResponseServer() {
	// Message handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err != nil {
			log.Errorf("Couldn't read request body: %s", err)
			return
		}
		log.Infof("Got request body: %s", data)
		route := ""
		if data["route"] != nil {
			route = data["route"].(string)
		}
		dataBody := interface{}(nil)
		if data["data"] != nil {
			dataBody = data["data"].(interface{})
		}
		// try to get config
		configBody := interface{}(nil)
		if data["config"] != nil {
			configBody = data["config"].(interface{})
		}
		replyTo := ""
		if data["replyTo"] != nil {
			replyTo = data["replyTo"].(string)
		}
		payload := Payload{Data: dataBody, Config: configBody}
		body, errMarshal := json.Marshal(payload)
		if errMarshal != nil {
			log.Errorf("Couldn't marshal body")
			return
		}
		err = ch.Publish(
			config.Exchange, // exchange
			route,           // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
				ReplyTo:     replyTo,
			})
		failOnError(err, "Failed to publish a message")
		log.Infof(" [x] Sent to %s - %s", route, body)
	})
	// Health checks handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Infof(" [x] Healthcheck..")
	})
	err := http.ListenAndServe(config.ServerListen, nil)
	failOnError(err, "Failed to start a server")
	log.Infof("Started server on: %s", config.ServerListen)
}
