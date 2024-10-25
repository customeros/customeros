package test

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

func InitTestRabbitMQ() (testcontainers.Container, *amqp091.Connection) {
	var ctx = context.Background()

	// Set up RabbitMQ container
	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.8-management",         // use a RabbitMQ image with management plugin
		ExposedPorts: []string{"5672/tcp", "15672/tcp"}, // 5672 for AMQP, 15672 for management interface
		WaitingFor:   wait.ForListeningPort("5672/tcp"), // wait until port 5672 is accessible
	}

	rabbitmqContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Panic("Failed to start RabbitMQ container:", err)
	}

	// Retrieve connection details
	host, err := rabbitmqContainer.Host(ctx)
	if err != nil {
		panic(err)
	}
	port, err := rabbitmqContainer.MappedPort(ctx, "5672")
	if err != nil {
		panic(err)
	}

	// Create connection string
	connString := fmt.Sprintf("amqp://guest:guest@%s:%s/", host, port.Port())

	// Attempt to connect to RabbitMQ with retries
	var rabbitConn *amqp091.Connection
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		rabbitConn, err = amqp091.Dial(connString)
		if err == nil {
			// success, break out of loop
			break
		}
		if i == maxRetries-1 {
			// last attempt failed, panic
			panic("Failed to connect to RabbitMQ after multiple attempts: " + err.Error())
		}
		// sleep before retrying
		time.Sleep(1 * time.Second)
	}

	return rabbitmqContainer, rabbitConn
}

func createQueues(conn *amqp091.Connection) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Set up the exchange
	exchangeName := "customeros"
	exchangeType := "fanout"
	if err := declareExchange(channel, exchangeName, exchangeType); err != nil {
		return fmt.Errorf("Failed to declare exchange: %v", err)
	}

	// Set up the queue
	queueName := "events"
	queue, err := declareQueue(channel, queueName)
	if err != nil {
		return fmt.Errorf("Failed to declare queue: %v", err)
	}

	// Bind the queue to the exchange
	routingKey := "*"
	if err := bindQueue(channel, queue.Name, exchangeName, routingKey); err != nil {
		return fmt.Errorf("Failed to bind queue to exchange: %v", err)
	}

	return nil
}

func TerminateRabbitMq(container testcontainers.Container, ctx context.Context) {
	err := container.Terminate(ctx)
	if err != nil {
		log.Fatal("Container should stop")
	}
}

func declareExchange(channel *amqp091.Channel, exchangeName, exchangeType string) error {
	return channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type of exchange (e.g., "direct", "fanout", "topic")
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func declareQueue(channel *amqp091.Channel, queueName string) (amqp091.Queue, error) {
	return channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func bindQueue(channel *amqp091.Channel, queueName, exchangeName, routingKey string) error {
	return channel.QueueBind(
		queueName,    // name of the queue
		routingKey,   // binding key
		exchangeName, // source exchange
		false,        // no-wait
		nil,          // arguments
	)
}
