package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type EmbeddingService interface {
	Segment(ctx context.Context, input string) (*ContentSegments, error)
	EmbedWebpage(ctx context.Context, webpage postgres_entity.GlobalOrganizationWebpages) ([]string, error)
	GetEmbedding(ctx context.Context, content string, task enum.EmbeddingTask) ([]float64, error)
}

type ContentSegments struct {
	NumTokens int      `json:"num_tokens"`
	Tokenizer string   `json:"tokenizer"`
	NumChunks int      `json:"num_chunks"`
	Chunks    []string `json:"chunks"`
}
