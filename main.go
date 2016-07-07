package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	host             = "amqp://localhost:5672/"
	exchange         = "exynize"
	routingKey       = "exynize.test"
	responseEndpoint = "http://localhost:3000/"
	serverListen     = ":8080"
)

// Response is a structure for sending response from proxy to consumer
type Response struct {
	Body    interface{} `json:"body"`    // body of the message
	ReplyTo string      `json:"replyTo"` // replyto param from rabbitmq
	Route   string      `json:"route"`   // route param from rabbit
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// load config from environment
	envProduction := os.Getenv("GO_ENV") == "production"
	envHost := os.Getenv("EXYNIZE_HOST")
	if envHost != "" {
		host = envHost
	}
	envExchange := os.Getenv("EXYNIZE_EXCHANGE")
	if envExchange != "" {
		exchange = envExchange
	}
	envServerListen := os.Getenv("EXYNIZE_LISTEN")
	if envServerListen != "" {
		serverListen = envServerListen
	}

	// configure logger for production
	if envProduction == true {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// TODO: Output to ELK instead of stdout, could also be a file.
		// log.SetOutput(os.Stderr)
		// Only log the warning severity or above.
		log.SetLevel(log.WarnLevel)
	}

	// connect to rabbit
	log.Infof("Connecting to %s with exchange %s", host, exchange)
	conn, err := amqp.Dial(host)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			decoder := json.NewDecoder(r.Body)
			var data map[string]interface{}
			err = decoder.Decode(&data)
			if err != nil {
				log.Errorf("Couldn't read request body")
				return
			}
			log.Infof("Got request body: %s", data)
			route := data["route"].(string)
			dataBody := data["data"].(interface{})
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
				})
			failOnError(err, "Failed to publish a message")
			log.Infof(" [x] Sent %s", body)
		})
		err = http.ListenAndServe(serverListen, nil)
		failOnError(err, "Failed to start a server")
		log.Infof("Started server on: %s", serverListen)
	}()

	go func() {
		for d := range msgs {
			var body map[string]interface{}
			err = json.Unmarshal(d.Body, &body)
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
	}()

	log.Infof(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
