package rabbit

import (
	"encoding/json"
	"strconv"
	"time"

	"gitlab.com/exynize/proxy/config"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

// HeartbeatInfo to be sent over heartbeat topic
type HeartbeatInfo struct {
	ID          string
	Name        string
	Description string
	Version     string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendHeartBeats() {
	// generate heartbeat info object
	dataBody := HeartbeatInfo{
		ID:          config.Cfg.ID,
		Name:        config.Cfg.Name,
		Description: config.Cfg.Description,
		Version:     config.Cfg.Version,
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
					config.Exchange,     // exchange
					"exynize.heartbeat", // routing key
					false,               // mandatory
					false,               // immediate
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
