package app

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/orchestrator/v2/config"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/orchestrator/v2/listener"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/orchestrator/v2/models"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/orchestrator/v2/sender"
)

var channel0 chan models.ServiceMessage
var channel1 chan models.ServiceMessage
var channel2 chan models.ServiceMessage
var channel3 chan models.ServiceMessage
var channel4 chan models.ServiceMessage
var channel5 chan models.ServiceMessage
var channel6 chan models.ServiceMessage
var channel7 chan models.ServiceMessage
var channel8 chan models.ServiceMessage

// Start initializes the application.
func Start() {
	log.Println("Initializing database")
	configuration := config.ReadConfiguration()
	log.Println("Loaded configuration:", configuration)

	var wg sync.WaitGroup

	wg.Add(8)

	log.Println("Initialization complete")

	channel0 = make(chan models.ServiceMessage)
	channel1 = make(chan models.ServiceMessage)
	channel2 = make(chan models.ServiceMessage)
	channel3 = make(chan models.ServiceMessage)
	channel4 = make(chan models.ServiceMessage)
	channel5 = make(chan models.ServiceMessage)
	channel6 = make(chan models.ServiceMessage)
	channel7 = make(chan models.ServiceMessage)
	channel8 = make(chan models.ServiceMessage)

	go listener.StartListener(buildListenerConfiguration(configuration, 0), channel0, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 1), channel1, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 2), channel2, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 3), channel3, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 4), channel4, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 5), channel5, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 6), channel6, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 7), channel7, &wg)
	go listener.StartListener(buildListenerConfiguration(configuration, 8), channel8, &wg)
	go processServiceMessages(configuration)

	wg.Wait()
}

func processServiceMessages(config config.Configuration) {
	senderConfiguration := buildSenderConfiguration(config)
	for {
		select {
		case v := <-channel0:
			log.Println("Channel 0 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 1, -1)
			log.Println("Channel 0 message sent:", time.Now().UnixNano())
		case v := <-channel1:
			log.Println("Channel 1 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 2, 0)
			log.Println("Channel 1 message sent:", time.Now().UnixNano())
		case v := <-channel2:
			log.Println("Channel 2 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 3, 1)
			log.Println("Channel 2 message sent:", time.Now().UnixNano())
		case v := <-channel3:
			log.Println("Channel 3 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 4, 2)
			log.Println("Channel 3 message sent:", time.Now().UnixNano())
		case v := <-channel4:
			log.Println("Channel 4 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 5, 3)
			log.Println("Channel 4 message sent:", time.Now().UnixNano())
		case v := <-channel5:
			log.Println("Channel 5 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 6, 4)
			log.Println("Channel 5 message sent:", time.Now().UnixNano())
		case v := <-channel6:
			log.Println("Channel 6 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 7, 5)
			log.Println("Channel 6 message sent:", time.Now().UnixNano())
		case v := <-channel7:
			log.Println("Channel 7 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, 8, 6)
			log.Println("Channel 7 message sent:", time.Now().UnixNano())
		case v := <-channel8:
			log.Println("Channel 8 message received: ", time.Now().UnixNano())
			executeOrchestratorCommandByServiceIndex(senderConfiguration, v, -1, 7)
			log.Println("Channel 8 message sent:", time.Now().UnixNano())
		}
	}
}

func executeOrchestratorCommandByServiceIndex(senderConfiguration sender.Configuration, v models.ServiceMessage, nextServiceIndex int, previousServiceIndex int) {
	if v.ServiceState == "SUCCESS" {
		sendCommandByServiceIndex(senderConfiguration, nextServiceIndex)
	} else if v.ServiceState == "FAILURE" {
		sendFailureCommandByServiceIndex(senderConfiguration, previousServiceIndex)
	}
}

func sendCommandByServiceIndex(config sender.Configuration, queueIndex int) {
	queueName := fmt.Sprintf("/microservice_%d_queue", queueIndex)
	serviceMessage := models.ServiceMessage{
		ApplicationIndex: 0,
		ServiceState:     "S2",
	}
	sender.Send(config, queueName, serviceMessage)
}

func sendFailureCommandByServiceIndex(config sender.Configuration, queueIndex int) {
	queueName := fmt.Sprintf("/microservice_%d_queue", queueIndex)
	serviceMessage := models.ServiceMessage{
		ApplicationIndex: 0,
		ServiceState:     "FAILURE",
	}
	sender.Send(config, queueName, serviceMessage)
}

func buildListenerConfiguration(config config.Configuration, index int) listener.Configuration {
	queueName := fmt.Sprintf("/microservice_%d_queue_orchestrator", index)
	return listener.Configuration{
		Host:      config.AmqpServerConfiguration.Host,
		Port:      config.AmqpServerConfiguration.Port,
		Username:  config.AmqpServerConfiguration.Username,
		Password:  config.AmqpServerConfiguration.Password,
		QueueName: queueName,
	}
}

func buildSenderConfiguration(config config.Configuration) sender.Configuration {
	return sender.Configuration{
		Host:     config.AmqpServerConfiguration.Host,
		Port:     config.AmqpServerConfiguration.Port,
		Username: config.AmqpServerConfiguration.Username,
		Password: config.AmqpServerConfiguration.Password,
	}
}
