package datastore

import (
	"context"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientInstance *mongo.Client

var mongoOnce sync.Once

//NewMongoClient Returns mongodb connection.
func NewMongoClient() *mongo.Client {
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(getConnectionString())

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatalln(err)
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatalln(err)
		}
		clientInstance = client
	})
	return clientInstance
}

// NewDatabase - provides mongo database instance
func NewDatabase(client *mongo.Client) *mongo.Database {
	return client.Database("service_state")
}

func getConnectionString() string {
	return os.Getenv("MONGO_CONNECTION_STRING")
}
