package main

import (
	log "github.com/Sirupsen/logrus"
	"gitlab.com/exynize/proxy/app"
	"gitlab.com/exynize/proxy/rabbit"
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
