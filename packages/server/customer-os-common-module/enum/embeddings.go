package enum

type EmbeddingContentType string

const (
	EmbeddingWebpage EmbeddingContentType = "WEBPAGE"
)

func (t EmbeddingContentType) String() string {
	return string(t)
}

type EmbeddingTask string

const (
	EmbeddingClassification   EmbeddingTask = "classification"
	EmbeddingGeneric          EmbeddingTask = "text-matching"
	EmbeddingPassageRetrieval EmbeddingTask = "retrieval.passage"
	EmbeddingQuery            EmbeddingTask = "retrieval.query"
)

func (t EmbeddingTask) String() string {
	return string(t)
}
