//+build wireinject

package app

import (
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/datastore"
	"github.com/google/wire"
)

func Initialize() *App {
	wire.Build(NewApp, datastore.NewMongoClient, datastore.NewDatabase, datastore.WireSet)
	return &App{}
}
