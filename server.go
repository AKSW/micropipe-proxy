package main

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

func initServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err != nil {
			log.Errorf("Couldn't read request body")
			return
		}
		log.Infof("Got request body: %s", data)
		route := data["route"].(string)
		dataBody := data["data"].(interface{})
		replyTo := ""
		if data["replyTo"] != nil {
			replyTo = data["replyTo"].(string)
		}
		body, errMarshal := json.Marshal(dataBody)
		if errMarshal != nil {
			log.Errorf("Couldn't marshal body")
			return
		}
		err = ch.Publish(
			exchange, // exchange
			route,    // routing key
			false,    // mandatory
			false,    // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
				ReplyTo:     replyTo,
			})
		failOnError(err, "Failed to publish a message")
		log.Infof(" [x] Sent %s", body)
	})
	err := http.ListenAndServe(serverListen, nil)
	failOnError(err, "Failed to start a server")
	log.Infof("Started server on: %s", serverListen)
}
