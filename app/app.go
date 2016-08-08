package app

import (
	log "github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"gitlab.com/exynize/proxy/config"
)

// InitApp inits app
func InitApp() {
	// read config from env
	config.ReadEnvConfig()
	// configure logger for production
	if config.IsProduction == true {
		// Log as Logstash JSON instead of the default ASCII formatter
		log.SetFormatter(&logstash.LogstashFormatter{})
	}
	// read config from yaml
	config.ReadYamlConfig()
	// start app in new goroutine
	go startApp()
}
