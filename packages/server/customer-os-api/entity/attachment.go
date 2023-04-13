package entity

import (
	"fmt"
	"time"
)

type AttachmentEntity struct {
	Id        string
	CreatedAt *time.Time

	MimeType  string
	Name      string
	Extension string
	Size      int64

	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string

	DataloaderKey string
}

func (attachmentEntity AttachmentEntity) ToString() string {
	return fmt.Sprintf("id: %s", attachmentEntity.Id)
}

type AttachmentEntities []AttachmentEntity

func (attachmentEntity AttachmentEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Attachment,
		NodeLabel_Attachment + "_" + tenant,
	}
}
