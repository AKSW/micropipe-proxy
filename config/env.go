package config

import (
	"os"
	"strconv"
)

var (
	// Host points to RabbitMQ host
	Host = "amqp://localhost:5672/"
	// Exchange to be used with RabbitMQ
	Exchange = "microproxy"
	// RoutingKey to be used with RabbitMQ
	RoutingKey = "microproxy.default"
	// RoutingUniqueKey to be used with RabbitMQ
	RoutingUniqueKey = "microproxy.default.unique"
	// HeartbeatRoute to be used with RabbitMQ
	HeartbeatRoute = "microproxy.heartbeats"
	// ResponseEndpoint to be used for incoming messages
	ResponseEndpoint = "http://localhost:3000/"
	// ServerListen bind for server for replies
	ServerListen = ":8080"
	// IsProduction determines if running in production env
	IsProduction = false
	// SendHeartbeats determines whether proxy should send service info to heartbeat rabbit topic
	SendHeartbeats = false
)

// ReadEnvConfig reads config from environment
func ReadEnvConfig() {
	// load config from environment
	IsProduction = os.Getenv("GO_ENV") == "production"
	envHost := os.Getenv("MICROPROXY_RABBIT_HOST")
	if envHost != "" {
		Host = envHost
	}
	envExchange := os.Getenv("MICROPROXY_EXCHANGE")
	if envExchange != "" {
		Exchange = envExchange
	}
	envServerListen := os.Getenv("MICROPROXY_SERVER_LISTEN")
	if envServerListen != "" {
		ServerListen = envServerListen
	}
	envSendHeartbeats := os.Getenv("MICROPROXY_HEARTBEATS")
	if envSendHeartbeats != "" {
		parsedHeartbeats, err := strconv.ParseBool(envSendHeartbeats)
		if err != nil {
			SendHeartbeats = parsedHeartbeats
		}
	}
	envHeartbeat := os.Getenv("MICROPROXY_HEARTBEAT_ROUTE")
	if envHeartbeat != "" {
		HeartbeatRoute = envHeartbeat
	}
}
