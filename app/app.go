package app

import (
	"github.com/AKSW/micropipe-proxy/config"
	log "github.com/Sirupsen/logrus"
)

// InitApp inits app
func InitApp() {
	// read config from env
	config.ReadEnvConfig()
	// configure logger for production
	if config.IsProduction == true {
		// Log as Logstash JSON instead of the default ASCII formatter
		log.SetFormatter(&log.JSONFormatter{})
	}
	// read config from yaml
	config.ReadYamlConfig()
	// start app in new goroutine
	go startApp()
}
