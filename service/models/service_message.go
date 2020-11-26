package models

// ServiceMessage is the representation of messages used in the communication between all the microservices.
type ServiceMessage struct {
	ApplicationIndex int    `json:"applicationIndex"`
	ServiceState     string `json:"serviceState"`
}
