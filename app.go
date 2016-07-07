package main

import (
	"bufio"
	"io/ioutil"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ApplicationConfig describes application config.yml file
type ApplicationConfig struct {
	Name             string
	Route            string
	Command          string
	ResponseEndpoint string
}

var cfg ApplicationConfig

func initApp() {
	cfg = ApplicationConfig{}

	// read the whole file at once
	cfgYaml, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(cfgYaml), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("Got application config:")
	log.Info(cfg)

	// update configs
	if cfg.ResponseEndpoint != "" {
		responseEndpoint = cfg.ResponseEndpoint
	}
	if cfg.Route != "" {
		routingKey = cfg.Route
	}

	// start app in new goroutine
	go startApp()
}

func startApp() {
	// parse command
	parts := strings.Fields(cfg.Command)
	executable := parts[0]
	args := parts[1:]
	log.Infof("Starting app with command: %s - and args: %s", executable, args)

	// start app
	cmd := exec.Command(executable, args...)

	// listen for stdout
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error creating StdoutPipe for Cmd: %s", err)
	}
	scanner := bufio.NewScanner(stdOut)
	go func() {
		for scanner.Scan() {
			log.Infof("  >> Child log:  %s\n", scanner.Text())
		}
	}()

	// listen for stderr
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("Error creating StderrPipe for Cmd: %s", err)
	}

	stdErrScanner := bufio.NewScanner(stdErr)
	go func() {
		for stdErrScanner.Scan() {
			log.Fatalf("  >> Child error: %s\n", stdErrScanner.Text())
		}
	}()

	// start command
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Error starting app: %s", err)
	}

	// wait for command
	err = cmd.Wait()
	// go generate command will fail when no generate command find.
	if err != nil {
		if err.Error() != "exit status 1" {
			log.Fatalf("Error waiting for app: %s", err)
		}
	}
}
