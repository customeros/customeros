package listeners

import (
	"github.com/rabbitmq/amqp091-go"
	"log"
)

// MessageHandler defines the function type for message handlers
type MessageHandler func(delivery amqp091.Delivery)

// startListener is a generic listener function that consumes messages from a given queue
func startListener(conn *amqp091.Connection, queueName string, handler MessageHandler) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName,
		"",   // no consumer tag
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)

	if err != nil {

	}

	log.Printf("Listening for messages on %s...", queueName)
	for d := range msgs {
		handler(d)
	}
}
