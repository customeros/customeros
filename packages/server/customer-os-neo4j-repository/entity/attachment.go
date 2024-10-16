package entity

import (
	"time"
)

type AttachmentProperty string

const (
	AttachmentPropertyPublicUrl AttachmentProperty = "publicUrl"
)

type AttachmentEntity struct {
	DataLoaderKey

	Id        string
	CreatedAt *time.Time

	CdnUrl    string
	PublicUrl string
	BasePath  string
	FileName  string
	MimeType  string
	Size      int64

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type AttachmentEntities []AttachmentEntity
