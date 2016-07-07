package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
)

var (
	host             = "amqp://localhost:5672/"
	exchange         = "exynize"
	routingKey       = "exynize.test"
	responseEndpoint = "http://localhost:3000/"
	serverListen     = ":8080"
	isProduction     = false
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// read config from env
	readEnv()
	// configure logger for production
	if isProduction == true {
		// Log as Logstash JSON instead of the default ASCII formatter
		log.SetFormatter(&logstash.LogstashFormatter{})
	}

	// read app config and start the app
	initApp()

	// connect to rabbit
	connectToRabbit()

	// create endless channel
	forever := make(chan bool)

	// init REST server
	go initServer()

	// start consuming messages from queue
	go consumeMessages()

	log.Infof(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
