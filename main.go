package main

import (
	"context"
	"log"
	"time"

	"github.com/Azure/go-amqp"
)

func main() {
	log.Println("Application started")
	time.Sleep(5 * time.Second)
	log.Println("Waited 5 seconds")

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

	{
		sender, err := session.NewSender(amqp.LinkTargetAddress("/example-queue"))

		if err != nil {
			log.Fatalln("Creating sender link:", err)
		}

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

		err = sender.Send(ctx, amqp.NewMessage([]byte("Hello!")))

		if err != nil {
			log.Fatalln("Sending message:", err)
		}

		sender.Close(ctx)
		cancel()
	}

	{
		receiver, err := session.NewReceiver(
			amqp.LinkSourceAddress("/example-queue"),
			amqp.LinkCredit(10))

		if err != nil {
			log.Fatalln("Creating receiver link:", err)
		}

		defer func() {
			ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			receiver.Close(ctx)
			cancel()
		}()

		for {
			log.Println("Reading messages")
			msg, err := receiver.Receive(ctx)

			if err != nil {
				log.Fatalln("Reading message from AMQP:", err)
			}

			msg.Accept(context.Background())

			log.Printf("Message received: %s\n", msg.GetData())
		}
	}
}
