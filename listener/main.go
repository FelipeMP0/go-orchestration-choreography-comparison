package main

import (
	"context"
	"log"
	"time"

	"github.com/Azure/go-amqp"
)

func main() {
	log.Println("Listener started")

	time.Sleep(10 * time.Second)
	log.Println("Waited 10 seconds")

	client, err := amqp.Dial("amqp://activemq-artemis:61616",
		amqp.ConnSASLPlain("user-amq", "password-amq"))

	if err != nil {
		log.Println("Dialing AMQP server:", err)
	}

	defer client.Close()

	session, err := client.NewSession()

	if err != nil {
		log.Fatalln("Creating AMQP session:", err)
	}

	ctx := context.Background()

	receiver, err := session.NewReceiver(
		amqp.LinkSourceAddress("/test1/queues/test2"),
		amqp.LinkCredit(10))

	if err != nil {
		log.Fatalln("Creating receiver link:", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		receiver.Close(ctx)
		cancel()
	}()

	log.Println("Reading messages")

	for {
		msg, err := receiver.Receive(ctx)

		if err != nil {
			log.Fatalln("Reading message from AMQP:", err)
		}

		msg.Accept(context.Background())

		log.Printf("Message received: %s\n", msg.GetData())
	}
}
