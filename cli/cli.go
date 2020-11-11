package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/FelipeMP0/go-orchestration-choreography-comparison/cli/v2/operations"
)

func main() {
	args := os.Args

	switch args[1] {
	case "send":
		sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
		hostPtr := sendCmd.String("host", "localhost", "active mq host")
		portPtr := sendCmd.String("port", "8000", "active mq port")
		usernamePtr := sendCmd.String("username", "user", "active mq user")
		passwordPtr := sendCmd.String("password", "password", "active mq password")
		queuePtr := sendCmd.String("queue", "/example-queue", "active mq queue")
		stringMessagePtr := sendCmd.String("stringMessage", "example", "string message")

		sendCmd.Parse(args[2:])

		portInt, err := strconv.ParseInt(*portPtr, 0, 0)

		if err != nil {
			log.Panicf("error parsing string %s to int", *portPtr)
		}

		sendMessageInput := operations.SendMessageInput{
			Host:          *hostPtr,
			Port:          portInt,
			Username:      *usernamePtr,
			Password:      *passwordPtr,
			QueueName:     *queuePtr,
			StringMessage: *stringMessagePtr,
		}

		operations.Send(sendMessageInput)
	default:
		log.Println("command not found")
	}
}
