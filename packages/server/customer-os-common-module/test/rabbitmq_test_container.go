package test

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
)

func InitTestRabbitMQ() (testcontainers.Container, string) {
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

	return rabbitmqContainer, connString
}

func TerminateRabbitMQ(container testcontainers.Container, ctx context.Context) {
	err := container.Terminate(ctx)
	if err != nil {
		log.Fatal("Container should stop")
	}
}
