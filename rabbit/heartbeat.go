package rabbit

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/AKSW/micropipe-proxy/config"
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

// HeartbeatInfo to be sent over heartbeat topic
type HeartbeatInfo struct {
	// app info
	UID         string
	ID          string
	Name        string
	Description string
	Version     string

	// input schema
	InputSchema interface{}
	// output schema
	OutputSchema interface{}
	// config schema
	ConfigSchema interface{}
}

func sendHeartBeats() {
	if config.SendHeartbeats == false {
		return
	}

	// generate heartbeat info object
	dataBody := HeartbeatInfo{
		UID:          config.Cfg.UID,
		ID:           config.Cfg.ID,
		Name:         config.Cfg.Name,
		Description:  config.Cfg.Description,
		Version:      config.Cfg.Version,
		InputSchema:  config.Cfg.InputSchema,
		OutputSchema: config.Cfg.OutputSchema,
		ConfigSchema: config.Cfg.ConfigSchema,
	}
	// marshal it into json
	body, errMarshal := json.Marshal(dataBody)
	failOnError(errMarshal, "Couldn't marshal heartbeat info")

	// send heartbeats every 30 seconds
	ticker := time.NewTicker(30 * time.Second)
	// heartbeats expire in 30s
	expirationTime := strconv.Itoa(30 * 1000)

	// attach to ticker in goroutine
	go func() {
		for {
			select {
			case <-ticker.C:
				err := ch.Publish(
					config.Exchange,       // exchange
					config.HeartbeatRoute, // routing key
					false, // mandatory
					false, // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(body),
						Timestamp:   time.Now(),
						Expiration:  expirationTime,
					})
				failOnError(err, "Failed to publish a heartbeat message")
				log.Infof(" [x] Sent heartbeat: %s", body)
			}
		}
	}()

}
