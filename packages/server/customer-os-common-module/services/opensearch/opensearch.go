package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"strings"
)

type opensearchService struct {
	client *opensearch.Client
}

func NewOpensearchService(logger logger.Logger, config *config.OpensearchConfig) interfaces.OpensearchService {
	opensearchConfig := opensearch.Config{
		Addresses: []string{config.Url},
		Username:  config.Username,
		Password:  config.Password,
	}

	client, err := opensearch.NewClient(opensearchConfig)
	if err != nil {
		logger.Fatalf("error creating opensearch client: %s", err)
		return nil
	}

	return &opensearchService{
		client: client,
	}
}

// IndexDocument indexes a single document
func (c *opensearchService) IndexDocument(indexName string, document interface{}) error {
	jsonDoc, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:   indexName,
		Body:    strings.NewReader(string(jsonDoc)),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), c.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}
