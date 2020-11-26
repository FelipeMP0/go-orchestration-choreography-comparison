package main

import (
	"log"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/app"
)

func main() {
	log.Println("Application started")
	application := app.Initialize()
	application.Start()
}
