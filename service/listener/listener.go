package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/FelipeMP0/go-orchestration-choreography-comparison/service/v2/models"
)

// AMQPListener is a listener to a specific topic from an AMQP server.
type AMQPListener struct {
}

// Configuration represents the configuration needed to connect and to listen to Queue.
type Configuration struct {
	Host      string
	Port      int
	Username  string
	Password  string
	QueueName string
}

// StartListener starts an AMQP listener.
func (listener *AMQPListener) StartListener(config Configuration, c chan<- models.ServiceMessage, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(c)

	log.Println("Listener started")

	time.Sleep(10 * time.Second)

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

	receiver, err := session.NewReceiver(
		amqp.LinkSourceAddress(config.QueueName),
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

		data := msg.GetData()

		log.Printf("Message received: %s\n", data)

		serviceMessage := models.ServiceMessage{}

		json.Unmarshal(data, &serviceMessage)

		c <- serviceMessage
	}
}

// NewAMQPListener provides an instance of AMQPListener.
func NewAMQPListener() *AMQPListener {
	return &AMQPListener{}
}
