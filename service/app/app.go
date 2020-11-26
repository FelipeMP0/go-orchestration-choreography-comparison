package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/config"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/datastore"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/listener"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// App configurations.
type App struct {
	dbClient              *mongo.Client
	ServiceStateDatastore *datastore.ServiceStateDatastore
	AmqpListener          *listener.AMQPListener
}

// Start initializes the application.
func (app *App) Start() {
	log.Println("Initializing database")
	app.dbClient.Connect(context.Background())

	configuration := config.ReadConfiguration()
	log.Println("Loaded configuration:", configuration)

	listenerConfiguration := buildListenerConfiguration(configuration)
	log.Println("Listener configuration:", listenerConfiguration)

	amqpMessagesChannel := make(chan models.ServiceMessage)

	var wg sync.WaitGroup

	wg.Add(1)
	go app.AmqpListener.StartListener(listenerConfiguration, amqpMessagesChannel, &wg)
	go processServiceMessages(amqpMessagesChannel)
	log.Println("Initialization complete")
	wg.Wait()
}

// NewApp provides an instance of App.
func NewApp(
	dbClient *mongo.Client,
	serviceStateDatastore *datastore.ServiceStateDatastore,
	amqpListner *listener.AMQPListener) *App {
	return &App{
		dbClient:              dbClient,
		ServiceStateDatastore: serviceStateDatastore,
		AmqpListener:          amqpListner,
	}
}

func processServiceMessages(c <-chan models.ServiceMessage) {
	for value := range c {
		log.Println(value)
	}
}

func buildListenerConfiguration(config config.Configuration) listener.Configuration {
	queueName := fmt.Sprintf("/microservice_%d_queue", config.ApplicationIndex)
	return listener.Configuration{
		Host:      config.AmqpServerConfiguration.Host,
		Port:      config.AmqpServerConfiguration.Port,
		Username:  config.AmqpServerConfiguration.Username,
		Password:  config.AmqpServerConfiguration.Password,
		QueueName: queueName,
	}
}
