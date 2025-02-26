package dto

import "time"

type EmbeddingRecord struct {
	ID               string     `json:"id"`
	SourceContentID  string     `json:"sourceContentId"`
	SourceUrl        string     `json:"sourceUrl"`
	ContentType      string     `json:"contentType"`
	Content          string     `json:"content"`
	ContentCreatedAt time.Time  `json:"contentCreatedAt"`
	EmbeddingModel   string     `json:"embeddingModel"`
	Vector           []float64  `json:"vector"`
	EmbeddedAt       time.Time  `json:"embeddedAt"`
	Tags             []Tag      `json:"tags"`
	Summary          string     `json:"summary"`
	SummaryVector    []float64  `json:"summaryVector"`
	Questions        []Question `json:"questions"`
}

type Tag struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Question struct {
	Text   string    `json:"text"`
	Vector []float64 `json:"vector"`
}
