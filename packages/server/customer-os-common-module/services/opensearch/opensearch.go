package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type opensearchService struct {
	client *opensearch.Client
}

const (
	SEARCH_LIMIT           = 20
	SEMANTIC_SEARCH_WEIGHT = 70
	KEYWORD_SEARCH_WEIGHT  = 30
)

func NewOpensearchService(logger logger.Logger, config *config.OpensearchConfig) interfaces.OpensearchService {
	opensearchConfig := opensearch.Config{
		Addresses: []string{config.Url},
		Username:  config.Username,
		Password:  config.Password,
	}

	if config.Url == "" {
		logger.Fatalf("opensearch_url not set in env")
		return nil
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
func (c *opensearchService) UpsertDocument(ctx context.Context, indexName string, documentId *string, document interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "OpensearchService.UpsertDocument")
	defer span.Finish()

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

	res, err := req.Do(ctx, c.client)
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

	// Set default limit if not provided
	if searchParams.ResultsLimit == nil {
		limit := SEARCH_LIMIT
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

	res, err := req.Do(ctx, c.client)
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
              "engine": "nmslib",
              "parameters": {
                "ef_construction": 128,
                "m": 16
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
              "engine": "nmslib"
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
                  "engine": "nmslib"
                }
              }
            }
          }
        }
      },
      "settings": {
        "index": {
          "knn": true,
          "knn.algo_param.ef_search": 100,
          "number_of_shards": 5,
          "number_of_replicas": 1
        }
      }
    }`

	// Check if index exists
	existsReq := opensearchapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	// Use the passed context instead of creating a new one
	existsRes, err := existsReq.Do(ctx, c.client)
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

		// Use the passed context
		createRes, err := createReq.Do(ctx, c.client)
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
			err := fmt.Errorf("failed to create index, status: %d", createRes.StatusCode)
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}
