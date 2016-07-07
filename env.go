package main

import "os"

func readEnv() {
	// load config from environment
	isProduction = os.Getenv("GO_ENV") == "production"
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
}
