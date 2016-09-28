package main

import (
	"github.com/AKSW/micropipe-proxy/app"
	"github.com/AKSW/micropipe-proxy/rabbit"
	log "github.com/Sirupsen/logrus"
)

func main() {
	// read app config and start the app
	app.InitApp()
	// connect to rabbit, 3 retries
	rabbit.ConnectToRabbit(3)
	// create endless channel
	forever := make(chan bool)
	// init REST server
	go rabbit.InitResponseServer()
	// start consuming messages from queue
	go rabbit.ConsumeMessages()
	// log successful init
	log.Infof(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
