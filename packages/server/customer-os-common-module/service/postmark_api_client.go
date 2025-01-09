package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"strings"
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

// ListAllServers uses paging to retrieve all servers from Postmark.
func (c *postmarkApiClient) ListAllServers(ctx context.Context) ([]ServerDetails, error) {
	var (
		allServers []ServerDetails
		offset     = 0
		count      = 100 // or whatever page size you want
	)

	for {
		// 1. Call ListServers with the current offset
		resp, err := c.ListServers(ctx, count, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to list servers (offset=%d): %w", offset, err)
		}

		// 2. Accumulate servers
		allServers = append(allServers, resp.Servers...)

		// 3. If we've retrieved all servers, break
		if offset+count >= resp.TotalCount {
			break
		}

		// 4. Otherwise, increment offset and continue
		offset += count
	}

	return allServers, nil
}

// ListServers retrieves a list of Postmark servers using the account token.
func (c *postmarkApiClient) ListServers(ctx context.Context, count, offset int) (*ListServersResponse, error) {
	url := fmt.Sprintf("%s/servers?count=%d&offset=%d", c.baseURL, count, offset)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("listServers: unexpected status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var serversList ListServersResponse
	if err := json.NewDecoder(resp.Body).Decode(&serversList); err != nil {
		return nil, err
	}

	return &serversList, nil
}

type ListServersResponse struct {
	TotalCount int             `json:"TotalCount"`
	Servers    []ServerDetails `json:"Servers"`
}

type ServerDetails struct {
	ID        int64  `json:"ID"`
	Name      string `json:"Name"`
	APITokens []struct {
		ServerToken string `json:"ServerToken"`
	} `json:"ApiTokens"`
	// ...other fields as needed
}

func (c *postmarkApiClient) GetServerByName(ctx context.Context, targetName string) (*ServerDetails, error) {
	// 1. Pull all servers
	servers, err := c.ListAllServers(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Find the one matching by name
	for _, srv := range servers {
		if strings.EqualFold(srv.Name, targetName) {
			return &srv, nil
		}
	}

	return nil, nil
}
