package interfaces

import "context"

type OpensearchService interface {
	UpsertDocument(ctx context.Context, indexName string, documentId *string, document interface{}) error
}
