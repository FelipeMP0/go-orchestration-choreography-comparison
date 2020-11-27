package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/config"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/datastore"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/listener"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/models"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/sender"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// App configurations.
type App struct {
	dbClient              *mongo.Client
	ServiceStateDatastore *datastore.ServiceStateDatastore
	AmqpListener          *listener.AMQPListener
	AmqpSender            *sender.AMQPSender
}

var stateID primitive.ObjectID
var collectionName string

// Start initializes the application.
func (app *App) Start() {
	log.Println("Initializing database")
	app.dbClient.Connect(context.Background())

	configuration := config.ReadConfiguration()
	log.Println("Loaded configuration:", configuration)

	stateID = primitive.NewObjectID()
	collectionName = datastore.GenerateCollectionName(configuration.ApplicationIndex)

	initServiceState(app)

	listenerConfiguration := buildListenerConfiguration(configuration)
	log.Println("Listener configuration:", listenerConfiguration)

	amqpMessagesChannel := make(chan models.ServiceMessage)

	var wg sync.WaitGroup

	wg.Add(1)

	go app.AmqpListener.StartListener(listenerConfiguration, amqpMessagesChannel, &wg)
	go processServiceMessages(amqpMessagesChannel, app, configuration)

	log.Println("Initialization complete")

	wg.Wait()
}

// NewApp provides an instance of App.
func NewApp(
	dbClient *mongo.Client,
	serviceStateDatastore *datastore.ServiceStateDatastore,
	amqpListener *listener.AMQPListener,
	amqpSender *sender.AMQPSender) *App {
	return &App{
		dbClient:              dbClient,
		ServiceStateDatastore: serviceStateDatastore,
		AmqpListener:          amqpListener,
		AmqpSender:            amqpSender,
	}
}

func processServiceMessages(c <-chan models.ServiceMessage, app *App, config config.Configuration) {
	applicationIndex := config.ApplicationIndex
	senderConfiguration := buildSenderConfiguration(config)
	for value := range c {
		switch value.ApplicationIndex {
		case applicationIndex - 1:
			log.Println("Success Start time: ", time.Now().UnixNano())
			if !forceFail() {
				log.Println("Updating service state")
				update := bson.D{{"$set", bson.D{{"state", value.ServiceState}}}}
				app.ServiceStateDatastore.UpdateByID(context.TODO(), collectionName, stateID, update)
				queueName := fmt.Sprintf("/microservice_%d_queue", applicationIndex+1)
				serviceMessage := models.ServiceMessage{
					ApplicationIndex: applicationIndex,
					ServiceState:     "S2",
				}
				app.AmqpSender.Send(senderConfiguration, queueName, serviceMessage)
			} else {
				log.Println("Rolling back service state")
				sendFailureMessage(app, senderConfiguration, applicationIndex)
			}
			log.Println("Success End time: ", time.Now().UnixNano())
		case applicationIndex + 1:
			log.Println("Failure Start time: ", time.Now().UnixNano())
			if value.ServiceState == "FAILURE" {
				log.Println("Received message to roll back service state")
				update := bson.D{{"$set", bson.D{{"state", "S1"}}}}
				app.ServiceStateDatastore.UpdateByID(context.TODO(), collectionName, stateID, update)
				sendFailureMessage(app, senderConfiguration, applicationIndex)
			}
			log.Println("Failure End time: ", time.Now().UnixNano())
		}
	}
}

func sendFailureMessage(app *App, senderConfiguration sender.Configuration, applicationIndex int) {
	queueName := fmt.Sprintf("/microservice_%d_queue", applicationIndex-1)
	serviceMessage := models.ServiceMessage{
		ApplicationIndex: applicationIndex,
		ServiceState:     "FAILURE",
	}
	app.AmqpSender.Send(senderConfiguration, queueName, serviceMessage)
}

func initServiceState(app *App) {
	serviceState := models.ServiceState{
		ID:    stateID,
		State: "S1",
	}

	app.ServiceStateDatastore.Create(context.TODO(), collectionName, &serviceState)
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

func buildSenderConfiguration(config config.Configuration) sender.Configuration {
	return sender.Configuration{
		Host:     config.AmqpServerConfiguration.Host,
		Port:     config.AmqpServerConfiguration.Port,
		Username: config.AmqpServerConfiguration.Username,
		Password: config.AmqpServerConfiguration.Password,
	}
}

func forceFail() bool {
	log.Println("Force fail:", os.Getenv("FORCE_FAIL"))
	log.Println(os.Getenv("FORCE_FAIL") == "true")
	return os.Getenv("FORCE_FAIL") == "true"
}
