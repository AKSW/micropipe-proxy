package rabbit

import (
	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.com/exynize/proxy/config"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
)

// ConnectToRabbit establishes connection to RabbitMQ
func ConnectToRabbit() {
	log.Infof("Connecting to \"%s\" with exchange \"%s\"", config.Host, config.Exchange)
	var err error
	conn, err = amqp.Dial(config.Host)
	failOnError(err, "Failed to connect to RabbitMQ")

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

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		true,  // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	log.Infof("Binding queue %s to exchange %s with routing key %s", q.Name, config.Exchange, config.RoutingKey)
	err = ch.QueueBind(
		q.Name,            // queue name
		config.RoutingKey, // routing key
		config.Exchange,   // exchange
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
