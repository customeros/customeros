package interfaces

import "context"

type OpensearchService interface {
	EmbeddingsIndexCheck(ctx context.Context, indexName string) error
	UpsertDocument(ctx context.Context, indexName string, documentId *string, document interface{}) error
}
