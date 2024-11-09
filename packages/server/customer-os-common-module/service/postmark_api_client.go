package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
)

type postmarkApiClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewAPIClient creates a new instance of postmarkApiClient.
func NewPostmarkAPIClient(baseURL, apiKey string) *postmarkApiClient {
	return &postmarkApiClient{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

// CreateServer sends a request to create a new Postmark server.
func (c *postmarkApiClient) CreateServer(ctx context.Context, serverName, inboundWebhook, inboundForwardingDomain string) (*ServerResponse, error) {
	url := c.baseURL + "/servers"

	request := ServerRequest{
		Name:           serverName,
		InboundHookUrl: inboundWebhook,
		InboundDomain:  inboundForwardingDomain,
	}
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Account-Token", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var serverResponse ServerResponse
	if err = json.Unmarshal(body, &serverResponse); err != nil {
		return nil, err
	}

	return &serverResponse, nil
}

// ServerResponse represents the response structure from creating a server.
type ServerResponse struct {
	ID   int64  `json:"ID"`
	Name string `json:"Name"`
}

type ServerRequest struct {
	Name           string `json:"Name"`
	InboundHookUrl string `json:"InboundHookUrl"`
	InboundDomain  string `json:"InboundDomain"`
}
