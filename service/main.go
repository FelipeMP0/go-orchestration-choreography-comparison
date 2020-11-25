package main

import (
	"context"
	"log"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/app"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	log.Println("Application started")
	app := app.Initialize()

	id, _ := primitive.ObjectIDFromHex("5fb9c4a3027965ff3cdf07a1")
	log.Println(id)
	serviceState := models.ServiceState{
		ID:    id,
		State: "S2",
	}
	//app.ServiceStateDatastore.Create(context.TODO(), &serviceState)´
	update := bson.D{{"$set", bson.D{{"state", serviceState.State}}}}
	app.ServiceStateDatastore.UpdateByID(context.TODO(), serviceState.ID, update)
}
