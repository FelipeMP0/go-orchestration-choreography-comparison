package datastore

import (
	"context"
	"fmt"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ServiceStateDatastore - interface to interact with mongo db
type ServiceStateDatastore struct {
	db *mongo.Database
}

//Collection returns the collection reference.
func (r *ServiceStateDatastore) Collection(collectionName string) *mongo.Collection {
	return r.db.Collection(collectionName)
}

// Create creates a new service state
func (r *ServiceStateDatastore) Create(ctx context.Context, collectionName string, serviceState *models.ServiceState) error {
	_, err := r.Collection(collectionName).InsertOne(ctx, serviceState)
	if err != nil {
		return err
	}
	return nil
}

// UpdateByID updates the service state for the given id.
func (r *ServiceStateDatastore) UpdateByID(ctx context.Context, collectionName string, id primitive.ObjectID, update bson.D) error {
	opts := options.FindOneAndUpdate().SetUpsert(false)
	filter := bson.D{{"_id", id}}
	var updatedDocument bson.M
	return r.Collection(collectionName).FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
}

// NewServiceStateDatastore - provides an instance of ServiceStateDatastore
func NewServiceStateDatastore(db *mongo.Database) *ServiceStateDatastore {
	return &ServiceStateDatastore{
		db: db,
	}
}

// GenerateCollectionName generates the service state collection name for the given application index.
func GenerateCollectionName(applicationIndex int) string {
	return fmt.Sprintf("service_state_%d", applicationIndex)
}
