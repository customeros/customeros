package entity

import (
	"time"
)

type AttachmentProperty string

const (
	AttachmentPropertyPublicUrl          AttachmentProperty = "publicUrl"
	AttachmentPropertyPublicUrlExpiresAt AttachmentProperty = "publicUrlExpiresAt"
)

type AttachmentEntity struct {
	DataLoaderKey

	Id        string
	CreatedAt *time.Time

	CdnUrl             string
	BasePath           string
	FileName           string
	MimeType           string
	Size               int64
	PublicUrl          string
	PublicUrlExpiresAt *time.Time

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
}

type AttachmentEntities []AttachmentEntity
