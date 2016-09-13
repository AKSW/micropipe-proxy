package config

import "os"

var (
	// Host points to RabbitMQ host
	Host = "amqp://localhost:5672/"
	// Exchange to be used with RabbitMQ
	Exchange = "exynize"
	// RoutingKey to be used with RabbitMQ
	RoutingKey = "exynize.test"
	// RoutingUniqueKey to be used with RabbitMQ
	RoutingUniqueKey = "exynize.test.unique"
	// ResponseEndpoint to be used for incoming messages
	ResponseEndpoint = "http://localhost:3000/"
	// ServerListen bind for server for replies
	ServerListen = ":8080"
	// IsProduction determines if running in production env
	IsProduction = false
)

// ReadEnvConfig reads config from environment
func ReadEnvConfig() {
	// load config from environment
	IsProduction = os.Getenv("GO_ENV") == "production"
	envHost := os.Getenv("EXYNIZE_HOST")
	if envHost != "" {
		Host = envHost
	}
	envExchange := os.Getenv("EXYNIZE_EXCHANGE")
	if envExchange != "" {
		Exchange = envExchange
	}
	envServerListen := os.Getenv("EXYNIZE_LISTEN")
	if envServerListen != "" {
		ServerListen = envServerListen
	}
}
