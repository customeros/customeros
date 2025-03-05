package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type opensearchService struct {
	eventsClient *opensearch.Client
	aiClient     *opensearch.Client
}

const (
	CONTEXT_SEARCH_LIMIT   = 20
	SEMANTIC_SEARCH_WEIGHT = 70
	KEYWORD_SEARCH_WEIGHT  = 30
)

func NewOpensearchService(logger logger.Logger, config *config.OpensearchConfig) interfaces.OpensearchService {
	if config == nil {
		return &opensearchService{}
	}

	return &opensearchService{
		eventsClient: createOpensearchClient(logger, config.EventsUrl, config.EventsUsername, config.EventsPassword),
		aiClient:     createOpensearchClient(logger, config.AIUrl, config.AIUsername, config.AIPassword),
	}
}

func createOpensearchClient(logger logger.Logger, url, username, password string) *opensearch.Client {
	if url == "" || username == "" || password == "" {
		return nil
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{url},
		Username:  username,
		Password:  password,
	})
	if err != nil {
		logger.Error("Failed to create client", "error", err)
		return nil
	}
	return client
}

// add all supported indexes here
func (c *opensearchService) getClientForIndex(indexName string) (*opensearch.Client, error) {
	switch {
	case strings.HasPrefix(indexName, "events-") ||
		strings.HasPrefix(indexName, "llm-"):
		return c.eventsClient, nil

	case strings.HasPrefix(indexName, "webpage-") ||
		strings.HasPrefix(indexName, "email-"):
		return c.aiClient, nil

	default:
		return nil, errors.New("opensearch indexName unsupported")
	}
}

// IndexDocument indexes a single document
func (c *opensearchService) UpsertDocument(ctx context.Context, indexName string, documentId *string, document interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpensearchService.UpsertDocument")
	defer span.Finish()

	if c == nil {
		err := errors.New("Opensearch service is not initialized")
		tracing.TraceErr(span, err)
		return err
	}

	client, err := c.getClientForIndex(indexName)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to select opensearch client"))
		return err
	}
	if client == nil {
		err := fmt.Errorf("opensearch client is nil")
		return err
	}

	jsonDoc, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:   indexName,
		Body:    strings.NewReader(string(jsonDoc)),
		Refresh: "true",
	}

	if documentId != nil {
		req.DocumentID = *documentId
	}

	res, err := req.Do(ctx, client)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		tracing.TraceErr(span, fmt.Errorf("error indexing document: %s", res.String()))
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// HybridSearch performs both vector and keyword search and combines the results
func (c *opensearchService) HybridSearch(ctx context.Context, searchParams interfaces.HybridSearchRequest) ([]interfaces.HybridSearchResult, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "opensearchService.hybridSearch")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if c.aiClient == nil {
		err := fmt.Errorf("opensearch client is nil")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Set default limit if not provided
	if searchParams.ResultsLimit == nil {
		limit := CONTEXT_SEARCH_LIMIT
		searchParams.ResultsLimit = &limit
	}
	if searchParams.KeywordWeight == nil {
		keyword := KEYWORD_SEARCH_WEIGHT
		searchParams.KeywordWeight = &keyword
	}
	if searchParams.SemanticWeight == nil {
		semantic := SEMANTIC_SEARCH_WEIGHT
		searchParams.SemanticWeight = &semantic
	}

	keywordBoost := float64(*searchParams.KeywordWeight) / 100.0
	semanticBoost := float64(*searchParams.SemanticWeight) / 100.0

	// Create the hybrid search query
	searchBody := map[string]interface{}{
		"size": *searchParams.ResultsLimit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						// Vector search component
						"knn": map[string]interface{}{
							"vector": map[string]interface{}{
								"vector": searchParams.EmbeddedQuery,
								"k":      *searchParams.ResultsLimit,
								"boost":  semanticBoost, // Vector search weight
							},
						},
					},
					{
						// Keyword search component
						"multi_match": map[string]interface{}{
							"query":  searchParams.Query,
							"fields": []string{"content^2", "summary^1.5", "questions.text^1"},
							"type":   "best_fields",
							"boost":  keywordBoost, // Keyword search weight
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
	}

	// Add filters if provided
	if searchParams.Filter != nil && len(searchParams.Filter) > 0 {
		mustClauses := make([]map[string]interface{}, 0)

		for key, value := range searchParams.Filter {
			// Handle special case for tags which is a nested field
			if key == "tags" {
				if tags, ok := value.([]map[string]string); ok {
					for _, tag := range tags {
						mustClauses = append(mustClauses, map[string]interface{}{
							"nested": map[string]interface{}{
								"path": "tags",
								"query": map[string]interface{}{
									"bool": map[string]interface{}{
										"must": []map[string]interface{}{
											{
												"match": map[string]interface{}{
													"tags.type": tag["type"],
												},
											},
											{
												"match": map[string]interface{}{
													"tags.value": tag["value"],
												},
											},
										},
									},
								},
							},
						})
					}
				}
			} else {
				// Standard field filter
				mustClauses = append(mustClauses, map[string]interface{}{
					"match": map[string]interface{}{
						key: value,
					},
				})
			}
		}

		if len(mustClauses) > 0 {
			searchBody["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = mustClauses
		}
	}

	// Convert query to JSON
	jsonBody, err := json.Marshal(searchBody)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("error marshaling search query: %w", err)
	}

	// Create and execute search request
	req := opensearchapi.SearchRequest{
		Index: []string{searchParams.Index},
		Body:  strings.NewReader(string(jsonBody)),
	}

	res, err := req.Do(ctx, c.aiClient)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		err = fmt.Errorf("search error: %s", res.String())
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Parse the search results
	var searchResult struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Score  float64                       `json:"_score"`
				ID     string                        `json:"_id"`
				Source interfaces.HybridSearchResult `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("error parsing search response: %w", err)
	}

	// Extract results
	results := make([]interfaces.HybridSearchResult, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		// Add the score to the result
		hit.Source.Score = hit.Score
		results = append(results, hit.Source)
	}

	return results, nil
}

func (c *opensearchService) EmbeddingsIndexCheck(ctx context.Context, indexName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "opensearchService.EmbeddingsIndexCheck")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	mapping := `{
      "mappings": {
        "properties": {
          "id": { "type": "keyword" },
          "sourceContentId": { "type": "keyword" },
          "sourceUrl": { "type": "keyword" },
          "contentType": { "type": "keyword" },
          "content": { "type": "text", "analyzer": "standard" },
          "contentCreatedAt": { "type": "date" },
          "embeddingModel": { "type": "keyword" },
          "vector": {
            "type": "knn_vector", 
            "dimension": 1024,
            "method": {
              "name": "hnsw",
              "space_type": "cosinesimil",
              "engine": "lucene",
              "parameters": {
                "ef_construction": 256,
                "m": 24
              }
            }
          },
          "embeddedAt": { "type": "date" },
          "tags": {
            "type": "nested",
            "properties": {
              "type": { "type": "keyword" },
              "value": { "type": "keyword" }
            }
          },
          "summary": { "type": "text", "analyzer": "standard" },
          "summaryVector": {
            "type": "knn_vector",
            "dimension": 1024,
            "method": {
              "name": "hnsw",
              "space_type": "cosinesimil",
              "engine": "lucene"
            }
          },
          "questions": {
            "type": "nested",
            "properties": {
              "text": { "type": "text", "analyzer": "standard" },
              "vector": {
                "type": "knn_vector",
                "dimension": 1024,
                "method": {
                  "name": "hnsw",
                  "space_type": "cosinesimil",
                  "engine": "lucene"
                }
              }
            }
          }
        }
      },
      "settings": {
        "index": {
          "knn": true,
          "knn.algo_param.ef_search": 250,
          "number_of_shards": 3,
          "number_of_replicas": 1,
          "refresh_interval": "10s"
        }
      }
    }`

	return c.ensureIndexExists(ctx, indexName, mapping)
}

func (c *opensearchService) LLMObservabilityIndexCheck(ctx context.Context, indexName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpensearchService.LLMObservabilityIndexCheck")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	mapping := `{
      "mappings": {
        "properties": {
          "request_id": { "type": "keyword" },
          "timestamp": { "type": "date" },
          "user_id": { "type": "keyword" },
          "tenant": { "type": "keyword" },
          "model": { "type": "keyword" },
          "temperature": { "type": "float" },
          "prompt_tokens": { "type": "integer" },
          "response_tokens": { "type": "integer" },
          "total_tokens": { "type": "integer" },
          "cost_usd": { "type": "float" },
          "success": { "type": "boolean" },
          "error_message": { 
            "type": "text",
            "fields": {
              "keyword": { "type": "keyword", "ignore_above": 256 }
            }
          },
          "trace_id": { "type": "keyword" },
          "systemPrompt": { 
            "type": "text",
            "fields": {
              "keyword": { "type": "keyword", "ignore_above": 256 }
            }
          },
          "prompt": { 
            "type": "text",
            "fields": {
              "keyword": { "type": "keyword", "ignore_above": 256 }
            }
          },
          "response": { 
            "type": "text",
            "fields": {
              "keyword": { "type": "keyword", "ignore_above": 256 }
            }
          }
        }
      },
      "settings": {
        "index": {
          "number_of_shards": 3,
          "number_of_replicas": 1,
          "refresh_interval": "10s"
        }
      }
    }`

	return c.ensureIndexExists(ctx, indexName, mapping)
}

func (c *opensearchService) ensureIndexExists(ctx context.Context, indexName, mapping string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "opensearchService.ensureIndexExists")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("indexName", indexName)

	if c == nil {
		err := errors.New("Opensearch service is not initialized")
		tracing.TraceErr(span, err)
		return err
	}

	client, err := c.getClientForIndex(indexName)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to select opensearch client"))
		return err
	}
	if client == nil {
		err := fmt.Errorf("opensearch client is nil")
		return err
	}

	// Check if index exists
	existsReq := opensearchapi.IndicesExistsRequest{
		Index: []string{indexName},
	}
	existsRes, err := existsReq.Do(ctx, client)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("error checking if index exists: %w", err)
	}

	// Guard against nil response
	if existsRes == nil {
		err := fmt.Errorf("received nil response when checking if index exists")
		tracing.TraceErr(span, err)
		return err
	}

	// Safely close the response body
	if existsRes.Body != nil {
		defer existsRes.Body.Close()
	}

	// If index doesn't exist, create it
	if existsRes.StatusCode == 404 {
		createReq := opensearchapi.IndicesCreateRequest{
			Index: indexName,
			Body:  strings.NewReader(mapping),
		}
		createRes, err := createReq.Do(ctx, client)
		if err != nil {
			tracing.TraceErr(span, err)
			return fmt.Errorf("error creating index: %w", err)
		}

		// Guard against nil response
		if createRes == nil {
			err := fmt.Errorf("received nil response when creating index")
			tracing.TraceErr(span, err)
			return err
		}

		// Safely close the response body
		if createRes.Body != nil {
			defer createRes.Body.Close()
		}

		if createRes.StatusCode >= 300 {
			responseBody, _ := io.ReadAll(createRes.Body)
			errorMsg := fmt.Sprintf("failed to create index, status: %d, response: %s",
				createRes.StatusCode, string(responseBody))
			err := fmt.Errorf(errorMsg)
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
