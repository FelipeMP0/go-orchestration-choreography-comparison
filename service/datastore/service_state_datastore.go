package datastore

import (
	"context"

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

// CollectionName returns the collection name that this repository interacts with.
func (r *ServiceStateDatastore) CollectionName() string {
	return "service_state"
}

//Collection returns the collection reference.
func (r *ServiceStateDatastore) Collection() *mongo.Collection {
	return r.db.Collection(r.CollectionName())
}

// Create creates a new service state
func (r *ServiceStateDatastore) Create(ctx context.Context, user *models.ServiceState) error {
	_, err := r.Collection().InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// UpdateByID updates the service state for the given id.
func (r *ServiceStateDatastore) UpdateByID(ctx context.Context, id primitive.ObjectID, update bson.D) error {
	opts := options.FindOneAndUpdate().SetUpsert(false)
	filter := bson.D{{"_id", id}}
	var updatedDocument bson.M
	return r.Collection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
}

// NewServiceStateDatastore - provides an instance of ServiceStateDatastore
func NewServiceStateDatastore(db *mongo.Database) *ServiceStateDatastore {
	return &ServiceStateDatastore{
		db: db,
	}
}
