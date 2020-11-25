package app

import (
	"context"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/datastore"
	"go.mongodb.org/mongo-driver/mongo"
)

// App configurations.
type App struct {
	dbClient              *mongo.Client
	ServiceStateDatastore *datastore.ServiceStateDatastore
}

// Start initializes the application.
func (app *App) Start() {
	app.dbClient.Connect(context.Background())
}

// NewApp provides an instance of App.
func NewApp(dbClient *mongo.Client, serviceStateDatastore *datastore.ServiceStateDatastore) *App {
	return &App{
		dbClient:              dbClient,
		ServiceStateDatastore: serviceStateDatastore,
	}
}
