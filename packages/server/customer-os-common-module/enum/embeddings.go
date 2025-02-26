package enum

type EmbeddingContentType string

const (
	EmbeddingWebpage EmbeddingContentType = "WEBPAGE"
)

func (t EmbeddingContentType) String() string {
	return string(t)
}
