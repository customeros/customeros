package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type EmbeddingService interface {
	Segment(ctx context.Context, input string) (*ContentSegments, error)
	EmbedWebpage(ctx context.Context, webpage postgres_entity.GlobalOrganizationWebpages) ([]string, error)
}

type ContentSegments struct {
	NumTokens int      `json:"num_tokens"`
	Tokenizer string   `json:"tokenizer"`
	NumChunks int      `json:"num_chunks"`
	Chunks    []string `json:"chunks"`
}
