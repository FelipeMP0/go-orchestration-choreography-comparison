package models

// ServiceMessage is the representation of messages used in the communication between all the microservices.
type ServiceMessage struct {
	ServiceState string `json:"serviceState"`
}
