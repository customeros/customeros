package nats_common

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/enums"

	"github.com/nats-io/nats.go"
)

type NatsHeader string

const (
	MAX_STREAM_RECONNECTS            = -1 // never stop trying to reconnect
	NATS_HEADER_TENANT    NatsHeader = "X-Tenant"
)

// NATSConnection represents a single NATS connection with its associated JetStream context
type NATSConnection struct {
	Conn *nats.Conn
	JS   nats.JetStreamContext
	Name enums.NatsStream // Stream name for logging
}

// NATSConnections holds all NATS connections
type NATSConnections struct {
	Streams map[enums.NatsStream]*NATSConnection // Map of stream name to connection details
}

func (n *NATSConnection) Close() error {
	if n.Conn != nil {
		n.Conn.Close()
		log.Printf("✅ NATS connection for stream %s closed", n.Name)
	}
	return nil
}

// Close gracefully closes all NATS connections
func (n *NATSConnections) Close() {
	if n.Streams == nil {
		return
	}

	// Close each stream connection gracefully
	for name, conn := range n.Streams {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Printf("❌ Error closing NATS connection for stream %s: %v", name, err)
			}
		}
	}
	log.Printf("✅ All NATS connections closed gracefully")
}

// InitNatsStream initializes a NATS connection for a specific stream
func InitNatsStream(serverList string, streamName enums.NatsStream) (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name(streamName.String()),
		nats.MaxReconnects(MAX_STREAM_RECONNECTS),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("NATS %s disconnected: %v", streamName, err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("NATS %s reconnected to %s", streamName, nc.ConnectedUrl())
		}),
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			log.Printf("NATS %s error: %v", streamName, err)
		}),
	}

	conn, err := nats.Connect(serverList, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS for %s stream: %v", streamName, err)
	}

	return conn, nil
}

// InitNats initializes the NATS connections and sets up streams
func InitNats(config *config.NATSConfig, environment string, streams []enums.NatsStream) (*NATSConnections, error) {
	if config == nil {
		return nil, fmt.Errorf("❌ NATS config is nil")
	}
	replicas := 3
	if environment != "production" {
		replicas = 1
	}

	// Setup NATS servers
	servers := []string{config.Node1}
	if config.Node2 != "" {
		servers = append(servers, config.Node2)
	}
	if config.Node3 != "" {
		servers = append(servers, config.Node3)
	}
	serverList := strings.Join(servers, ",")

	// Initialize connections map
	streamConns := make(map[enums.NatsStream]*NATSConnection)

	// Initialize all streams
	for _, streamName := range streams {
		// Create connection
		conn, err := InitNatsStream(serverList, streamName)
		if err != nil {
			// Clean up any existing connections
			for _, c := range streamConns {
				c.Close()
			}
			return nil, fmt.Errorf("❌ Failed to initialize NATS stream %s: %v", streamName, err)
		}
		log.Printf("✅ NATS connection established for stream %s", streamName)

		// Create JetStream context for each stream
		js, err := conn.JetStream()
		if err != nil {
			for _, c := range streamConns {
				c.Close()
			}
			return nil, fmt.Errorf("❌ Failed to create JetStream context for %s: %v", streamName, err)
		}
		log.Printf("✅ JetStream context created for stream %s", streamName)

		// Create NATSConnection instance
		streamConns[streamName] = &NATSConnection{
			Conn: conn,
			JS:   js,
			Name: streamName,
		}

		// Set up stream
		err = setupWorkQueueStream(js, streamName.String(), []string{streamName.String() + ".>"}, replicas)
		if err != nil {
			for _, c := range streamConns {
				c.Close()
			}
			return nil, fmt.Errorf("❌ Failed to set up stream %s: %v", streamName, err)
		}
		log.Printf("✅ Stream %s configured successfully", streamName)
	}

	log.Println("✅ NATS initialized successfully")

	return &NATSConnections{
		Streams: streamConns,
	}, nil
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

// GetConnection returns the NATS connection for a given stream
func (n *NATSConnections) GetNatsConnection(streamName enums.NatsStream) (*NATSConnection, error) {
	if conn, exists := n.Streams[streamName]; exists && conn != nil {
		return conn, nil
	}
	return nil, fmt.Errorf("❌ No connection found for stream %s", streamName)
}

// GetConnection returns the NATS connection for a given stream
func (n *NATSConnections) GetConnection(streamName enums.NatsStream) (*nats.Conn, error) {
	if conn, exists := n.Streams[streamName]; exists && conn != nil {
		return conn.Conn, nil
	}
	return nil, fmt.Errorf("❌ No connection found for stream %s", streamName)
}

// GetJetStream returns the JetStream context for a given stream
func (n *NATSConnections) GetJetStream(streamName enums.NatsStream) (nats.JetStreamContext, error) {
	if conn, exists := n.Streams[streamName]; exists && conn != nil {
		return conn.JS, nil
	}
	return nil, fmt.Errorf("❌ No JetStream context found for stream %s", streamName)
}

// GetStream returns the NATSConnection for a given stream
func (n *NATSConnections) GetStream(streamName enums.NatsStream) (*NATSConnection, error) {
	if conn, exists := n.Streams[streamName]; exists && conn != nil {
		return conn, nil
	}
	return nil, fmt.Errorf("❌ No stream found for name %s", streamName)
}
