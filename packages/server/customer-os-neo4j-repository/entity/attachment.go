package entity

import (
	"time"
)

type AttachmentEntity struct {
	DataLoaderKey

	Id        string
	CreatedAt *time.Time

	CdnUrl   string
	BasePath string
	FileName string
	MimeType string
	Size     int64

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type AttachmentEntities []AttachmentEntity
