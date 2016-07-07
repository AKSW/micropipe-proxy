package main

import log "github.com/Sirupsen/logrus"

var (
	host             = "amqp://localhost:5672/"
	exchange         = "exynize"
	routingKey       = "exynize.test"
	responseEndpoint = "http://localhost:3000/"
	serverListen     = ":8080"
	isProduction     = false
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
	// read config from env
	readEnv()
	// configure logger for production
	if isProduction == true {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})
		// TODO: Output to ELK instead of stdout, could also be a file.
		// log.SetOutput(os.Stderr)
		// Only log the warning severity or above.
		log.SetLevel(log.WarnLevel)
	}

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
