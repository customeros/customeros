package interfaces

type OpensearchService interface {
	IndexDocument(indexName string, document interface{}) error
}
