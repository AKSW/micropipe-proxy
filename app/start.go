package app

import (
	"bufio"
	"os/exec"
	"strings"

	"github.com/AKSW/micropipe-proxy/config"
	log "github.com/Sirupsen/logrus"
)

func startApp() {
	// parse command
	parts := strings.Fields(config.Cfg.Command)
	log.Infof("cmd: %s, parts: %s", config.Cfg.Command, parts)
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
			log.Infof("  >> Child log:  %s", scanner.Text())
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
			log.Fatalf("  >> Child error: %s", stdErrScanner.Text())
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
		log.Fatalf("Error waiting for app: %s", err)
	}

	log.Fatal("Error waiting for app - app exited!")
}
