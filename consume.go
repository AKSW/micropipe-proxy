package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	js "github.com/xeipuuv/gojsonschema"
)

// Response is a structure for sending response from proxy to consumer
type Response struct {
	Body     interface{} `json:"body"`     // body of the message
	ReplyTo  string      `json:"replyTo"`  // replyto param from rabbitmq
	Route    string      `json:"route"`    // route param from rabbit
	NewRoute string      `json:"newRoute"` // new route for next message
}

func validateMessage(body map[string]interface{}) error {
	schema := js.NewGoLoader(cfg.InputSchema)
	doc := js.NewGoLoader(body)

	result, err := js.Validate(schema, doc)
	if err != nil {
		log.Errorf("Error validating message: %s", err)
		return err
	}

	if result.Valid() {
		log.Infof("The document is valid")
		return nil
	}

	log.Errorf("The document is not valid. see errors:")
	for _, desc := range result.Errors() {
		log.Errorf("- %s", desc)
	}
	return errors.New("Docment not valid")
}

func consumeMessages() {
	// prepare route replacement regex
	reg, err := regexp.Compile(cfg.Route + "-" + cfg.Version + "(.?)")
	if err != nil {
		log.Fatalf("Error compiling route regex: %s", err)
	}

	// consume messages
	for d := range msgs {
		// try to unmarshal incoming data
		var body map[string]interface{}
		err := json.Unmarshal(d.Body, &body)
		if err != nil {
			log.Errorf("Couldn't decode message body")
			continue
		}
		replyTo := d.ReplyTo
		route := d.RoutingKey
		newRoute := ""
		if replyTo != "" {
			newRoute = replyTo
		} else {
			newRoute = reg.ReplaceAllString(route, "")
		}
		log.Infof(" [x] Got:\n  - body: %s\n  - replyTo: %s\n  - route: %s\n  - newRoute: %s", body, replyTo, route, newRoute)

		// validate document using input schema
		err = validateMessage(body)
		if err != nil {
			log.Errorf("Error validating input: %s", err)
			continue
		}

		// create response body
		r := Response{Body: body, ReplyTo: replyTo, Route: route, NewRoute: newRoute}
		// try to marshal it to json
		rbody, errMarshal := json.Marshal(r)
		if errMarshal != nil {
			log.Errorf("Couldn't marshal response body")
			continue
		}
		log.Infof(" [x] prepared response: %s", rbody)

		// try sending it to the response endpoint
		resp, err := http.Post(responseEndpoint, "application/json", bytes.NewBuffer(rbody))
		if err != nil {
			log.Errorf("Couldn't send POST request to consumer")
		} else {
			log.Infof(" [x] Sent via REST: %s", resp.Status)
		}
	}
}
