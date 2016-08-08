package config

import (
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ApplicationConfig describes application config.yml file
type ApplicationConfig struct {
	// app info
	ID          string `yaml:"id"`
	Name        string
	Description string
	Version     string

	// command to execute
	Command string

	// optional override for responseEndpoint
	ResponseEndpoint string `yaml:"responseEndpoint"`

	// input schema
	InputSchema interface{} `yaml:"inputSchema"`
	// output schema
	OutputSchema interface{} `yaml:"outputSchema"`
	// config schema
	ConfigSchema interface{} `yaml:"configSchema"`
}

// Cfg holds application current config
var Cfg ApplicationConfig

// Fixes parsed YAML to conform to JSON string->interface{} format
func fixJSON(input map[interface{}]interface{}) map[string]interface{} {
	fixedInput := make(map[string]interface{})
	for key, value := range input {
		switch key := key.(type) {
		case string:
			switch value := value.(type) {
			case string:
				fixedInput[key] = value
			case map[interface{}]interface{}:
				fixedInput[key] = fixJSON(value)
			}
		default:
			log.Debugf("other key: %s", key)
		}
	}
	return fixedInput
}

// ReadYamlConfig reads config from file
func ReadYamlConfig() {
	Cfg = ApplicationConfig{}

	// read the whole file at once
	cfgYaml, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal([]byte(cfgYaml), &Cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("Got application config:")
	Cfg.InputSchema = fixJSON(Cfg.InputSchema.(map[interface{}]interface{}))
	Cfg.OutputSchema = fixJSON(Cfg.OutputSchema.(map[interface{}]interface{}))
	Cfg.ConfigSchema = fixJSON(Cfg.ConfigSchema.(map[interface{}]interface{}))
	log.Info(Cfg)

	// update configs
	if Cfg.ResponseEndpoint != "" {
		ResponseEndpoint = Cfg.ResponseEndpoint
	}
	RoutingKey = Cfg.ID + "-" + Cfg.Version + ".#"
}
