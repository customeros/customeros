package interfaces

import "context"

type OpensearchService interface {
	EmbeddingsIndexCheck(ctx context.Context, indexName string) error
	UpsertDocument(ctx context.Context, indexName string, documentId *string, document interface{}) error
	HybridSearch(ctx context.Context, searchParams HybridSearchRequest) ([]HybridSearchResult, error)
}

type HybridSearchRequest struct {
	Index          string
	Query          string
	EmbeddedQuery  []float64
	ResultsLimit   *int
	Filter         map[string]interface{}
	SemanticWeight *int // out of 100
	KeywordWeight  *int // out of 100
}

type HybridSearchResult struct {
	ID               string     `json:"id"`
	SourceContentID  string     `json:"sourceContentId"`
	SourceURL        string     `json:"sourceUrl"`
	ContentType      string     `json:"contentType"`
	Content          string     `json:"content"`
	ContentCreatedAt string     `json:"contentCreatedAt"`
	EmbeddingModel   string     `json:"embeddingModel"`
	Vector           []float64  `json:"vector"`
	EmbeddedAt       string     `json:"embeddedAt"`
	Tags             []Tag      `json:"tags"`
	Summary          string     `json:"summary"`
	SummaryVector    []float64  `json:"summaryVector"`
	Questions        []Question `json:"questions"`
	Score            float64    `json:"_score"`
}

type Tag struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Question struct {
	Text   string    `json:"text"`
	Vector []float64 `json:"vector"`
}
