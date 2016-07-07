package main

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
)

// HeartbeatInfo to be sent over heartbeat topic
type HeartbeatInfo struct {
	Name        string
	Description string
	Version     string
	Route       string
}

func connectToRabbit() {
	log.Infof("Connecting to \"%s\" with exchange \"%s\"", host, exchange)
	var err error
	conn, err = amqp.Dial(host)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		true,  // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Infof("Binding queue %s to exchange %s with routing key %s", q.Name, exchange, routingKey)
	err = ch.QueueBind(
		q.Name,     // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// send heartbeats
	sendHeartBeats()
}

func sendHeartBeats() {
	// generate heartbeat info object
	dataBody := HeartbeatInfo{
		Name:        cfg.Name,
		Description: cfg.Description,
		Version:     cfg.Version,
		Route:       cfg.Route,
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
					exchange,            // exchange
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
