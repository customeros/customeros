package test

import (
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/context"
)

func SetupRabbitMQTestContainer() (string, func(), error) {
	ctx := context.Background()

	// Set up RabbitMQ container request
	rabbitMQContainerReq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.9-management", // or a specific version
		ExposedPorts: []string{"5672/tcp", "15672/tcp"},
		WaitingFor:   wait.ForLog("Server startup complete"), // Wait for RabbitMQ to start
	}

	// Start the RabbitMQ container
	rabbitMQContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: rabbitMQContainerReq,
		Started:          true,
	})
	if err != nil {
		return "", nil, err
	}

	// Get the mapped port for RabbitMQ's AMQP interface
	hostPort, err := rabbitMQContainer.MappedPort(ctx, "5672")
	if err != nil {
		return "", nil, err
	}

	// Construct RabbitMQ URL for testing
	rabbitMQURL := fmt.Sprintf("amqp://guest:guest@localhost:%s", hostPort.Port())

	// Cleanup function to terminate the container
	cleanup := func() {
		rabbitMQContainer.Terminate(ctx)
	}

	return rabbitMQURL, cleanup, nil
}
