package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func consumeMessages() {
	for d := range msgs {
		var body map[string]interface{}
		err := json.Unmarshal(d.Body, &body)
		if err != nil {
			log.Errorf("Couldn't decode message body")
			return
		}
		replyTo := d.ReplyTo
		route := d.RoutingKey
		log.Infof(" [x] Got:\n  - body: %s\n  - replyTo: %s\n  - route: %s", body, replyTo, route)
		r := Response{Body: body, ReplyTo: replyTo, Route: route}
		rbody, errMarshal := json.Marshal(r)
		if errMarshal != nil {
			log.Errorf("Couldn't marshal response body")
			return
		}
		log.Infof(" [x] prepared response: %s", rbody)
		resp, err := http.Post(responseEndpoint, "application/json", bytes.NewBuffer(rbody))
		if err != nil {
			log.Errorf("Couldn't send POST request to consumer")
		} else {
			log.Infof(" [x] Sent via REST: %s", resp.Status)
		}
	}
}
