package operations

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-amqp"
)

// SendMessageInput represents the configuration to connect to a ActiveMQ instance,
// and the message content to be send.
type SendMessageInput struct {
	Host          string
	Port          int64
	Username      string
	Password      string
	QueueName     string
	StringMessage string
}

// Send sends a message to ActimeMQ.
func Send(input SendMessageInput) {
	log.Println("Sending message: ", input)

	hostAndPort := fmt.Sprintf("amqp://%s:%d", input.Host, input.Port)
	client, err := amqp.Dial(hostAndPort,
		amqp.ConnSASLPlain(input.Username, input.Password))

	if err != nil {
		log.Println("Dialing AMQP server:", err)
	}

	defer client.Close()

	session, err := client.NewSession()

	if err != nil {
		log.Fatalln("Creating AMQP session:", err)
	}

	ctx := context.Background()

	sender, err := session.NewSender(amqp.LinkTargetAddress(input.QueueName))

	if err != nil {
		log.Fatalln("Creating sender link:", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	err = sender.Send(ctx, amqp.NewMessage([]byte(input.StringMessage)))

	if err != nil {
		log.Fatalln("Sending message:", err)
	}

	sender.Close(ctx)
	cancel()
}
