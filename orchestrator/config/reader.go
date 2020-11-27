package config

import (
	"log"

	"github.com/tkanos/gonfig"
)

// Configuration contains all configured properties.
type Configuration struct {
	ApplicationIndex        int
	AmqpServerConfiguration AMQPServerConfiguration
}

// AMQPServerConfiguration contains the needed properties to connect to the AMQP server.
type AMQPServerConfiguration struct {
	Host      string
	Port      int
	Username  string
	Password  string
	QueueName string
}

// ReadConfiguration reads the configuration properties from a json file.
func ReadConfiguration() Configuration {
	config := Configuration{}
	err := gonfig.GetConf("/orchestrator_configuration/config.json", &config)
	if err != nil {
		log.Fatalln("Loading configuration file:", err)
	}
	return config
}
