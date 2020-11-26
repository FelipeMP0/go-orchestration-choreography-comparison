package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/models"
)

// AMQPListener is a listener to a specific topic from an AMQP server.
type AMQPSender struct {
}

// Configuration represents the configuration needed to connect and to send to Queue.
type Configuration struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Send sends a message to ActimeMQ.
func (s *AMQPSender) Send(config Configuration, queueName string, serviceMessage models.ServiceMessage) {
	log.Println("Sending message")

	address := fmt.Sprintf("amqp://%s:%d", config.Host, config.Port)
	client, err := amqp.Dial(address,
		amqp.ConnSASLPlain(config.Username, config.Password))

	if err != nil {
		log.Println("Dialing AMQP server:", err)
	}

	defer client.Close()

	session, err := client.NewSession()

	if err != nil {
		log.Fatalln("Creating AMQP session:", err)
	}

	ctx := context.Background()

	sender, err := session.NewSender(amqp.LinkTargetAddress(queueName))

	if err != nil {
		log.Fatalln("Creating sender link:", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	data, _ := json.Marshal(serviceMessage)

	err = sender.Send(ctx, amqp.NewMessage([]byte(data)))

	if err != nil {
		log.Fatalln("Sending message:", err)
	}

	sender.Close(ctx)
	cancel()
}

// NewAMQPSender provides an instance of AMQPSender.
func NewAMQPSender() *AMQPSender {
	return &AMQPSender{}
}
