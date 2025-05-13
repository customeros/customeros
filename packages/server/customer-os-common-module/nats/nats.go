package nats_common

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

type NATSConnections struct {
	Conn *nats.Conn
	JS   nats.JetStreamContext
}

func (n *NATSConnections) Close() {
	if n.Conn != nil {
		n.Conn.Close()
		log.Println("NATS connection closed")
	}
}

type NatsHeader string

const (
	CORE_STREAM           = "core"
	LEADS_STREAM          = "leads"
	MAX_STREAM_RECONNECTS = -1 // never stop trying to reconnect

	NATS_HEADER_TENANT NatsHeader = "X-Tenant"
)

// InitNats initializes the NATS connection and sets up streams
func InitNats(config *config.NATSConfig, environment string) (*NATSConnections, error) {
	if config == nil {
		return nil, fmt.Errorf("NATS config is nil")
	}
	replicas := 3
	if environment != "production" {
		replicas = 1
	}

	opts := []nats.Option{
		nats.Name(CORE_STREAM),
		nats.MaxReconnects(MAX_STREAM_RECONNECTS),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("NATS error: %v", err)
		}),
	}

	// Setup NATS connection
	servers := []string{config.Node1}
	if config.Node2 != "" {
		servers = append(servers, config.Node2)
	}
	if config.Node3 != "" {
		servers = append(servers, config.Node3)
	}

	conn, err := nats.Connect(strings.Join(servers, ","), opts...)
	if err != nil {
		log.Printf("NATS error: %v", err)
		return nil, err
	}

	// Create JetStream context
	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %v", err)
	}

	// Set up streams centrally
	err = setupNATSStreams(js, replicas)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to set up NATS streams: %v", err)
	}

	log.Println("NATS initialized successfully")

	return &NATSConnections{
		Conn: conn,
		JS:   js,
	}, nil
}

func setupNATSStreams(js nats.JetStreamContext, replicas int) error {
	persistedSubjects := []string{
		"core.>",
	}
	return setupWorkQueueStream(js, CORE_STREAM, persistedSubjects, replicas)
}

func setupWorkQueueStream(js nats.JetStreamContext, streamName string, subjects []string, replicas int) error {
	streamInfo, err := js.StreamInfo(streamName)
	if err != nil {
		// Stream doesn't exist, create it
		_, err = js.AddStream(&nats.StreamConfig{
			Name:      streamName,
			Subjects:  subjects,
			Retention: nats.WorkQueuePolicy,
			Storage:   nats.FileStorage,
			Replicas:  replicas,
			MaxAge:    168 * time.Hour,
			Discard:   nats.DiscardOld,
		})
		if err != nil {
			return fmt.Errorf("failed to create stream: %w", err)
		}
		log.Printf("Created stream %s with subjects %v", streamName, subjects)
	} else {
		// Stream exists, verify subjects
		currentSubjects := streamInfo.Config.Subjects
		missingSubjects := findMissingSubjects(currentSubjects, subjects)

		if len(missingSubjects) > 0 {
			// Update stream with new subjects
			updatedSubjects := append(currentSubjects, missingSubjects...)
			_, err = js.UpdateStream(&nats.StreamConfig{
				Name:      streamName,
				Subjects:  updatedSubjects,
				Retention: streamInfo.Config.Retention,
				Storage:   streamInfo.Config.Storage,
				Replicas:  streamInfo.Config.Replicas,
				MaxAge:    streamInfo.Config.MaxAge,
				Discard:   streamInfo.Config.Discard,
			})
			if err != nil {
				return fmt.Errorf("failed to update stream subjects: %w", err)
			}
			log.Printf("Updated stream %s with additional subjects %v", streamName, missingSubjects)
		}
	}
	return nil
}

// Helper function to find subjects that are in the desired list but not in the current list
func findMissingSubjects(currentSubjects, desiredSubjects []string) []string {
	var missing []string
	for _, desired := range desiredSubjects {
		found := false
		for _, current := range currentSubjects {
			if current == desired {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, desired)
		}
	}
	return missing
}
