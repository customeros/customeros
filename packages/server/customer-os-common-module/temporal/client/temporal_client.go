package temporal_client

import (
	"fmt"

	"go.temporal.io/sdk/client"
)

func TemporalClient(hostPort, namespace string) (client.Client, error) {
	// Initialize a Temporal Client
	// Specify the Namespace in the Client options
	clientOptions := client.Options{
		HostPort:  hostPort, // Temporal Server "Host:Port"
		Namespace: namespace,
	}
	temporalClient, err := client.Dial(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to create a Temporal Client: %s", err.Error())
	}
	return temporalClient, nil
}
