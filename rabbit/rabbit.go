package rabbit

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/exynize/proxy/config"
)

var (
	conn       *amqp.Connection
	ch         *amqp.Channel
	msgs       <-chan amqp.Delivery
	uniqueMsgs <-chan amqp.Delivery
	maxRetries = 3
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// ConnectToRabbit establishes connection to RabbitMQ
func ConnectToRabbit(retries int) {
	if retries > maxRetries {
		retries = maxRetries
	}

	log.Infof("Connecting to \"%s\" with exchange \"%s\"", config.Host, config.Exchange)
	var err error
	conn, err = amqp.Dial(config.Host)
	if err != nil {
		if retries > 0 {
			newTries := retries - 1
			log.Infof("Failed to connect, retrying %s more times...", newTries)
			// sleep
			t := 2000 * time.Duration(maxRetries+1-retries) * time.Millisecond
			time.Sleep(t)
			ConnectToRabbit(newTries)
			return
		}

		failOnError(err, "Failed to connect to RabbitMQ")
	}

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		config.Exchange, // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// subscribe to normal messages
	subscribeTo(config.RoutingKey, false)

	// subscribe to messages from unique queue
	subscribeTo(config.RoutingUniqueKey, true)

	// start sending heartbeats
	sendHeartBeats()
}

func subscribeTo(routingKey string, exclusive bool) {
	// Create queue for normal message exchange
	q, err := ch.QueueDeclare(
		"",        // name
		true,      // durable
		true,      // delete when usused
		exclusive, // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	// Bind queue for normal message exchange
	log.Infof("Binding queue %s to exchange %s with routing key %s", q.Name, config.Exchange, routingKey)
	err = ch.QueueBind(
		q.Name,          // queue name
		routingKey,      // routing key
		config.Exchange, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
	// consume from normal queue
	msgs, err = ch.Consume(
		q.Name,    // queue
		"",        // consumer
		true,      // auto ack
		exclusive, // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")
}
