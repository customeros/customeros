package entity

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

// Deprecated, use neo4j module instead
type AttachmentEntity struct {
	Id        string
	CreatedAt *time.Time

	CdnUrl   string
	BasePath string
	FileName string
	MimeType string
	Size     int64

	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	AppSource     string

	DataloaderKey string
}

func (attachmentEntity AttachmentEntity) ToString() string {
	return fmt.Sprintf("id: %s", attachmentEntity.Id)
}

type AttachmentEntities []AttachmentEntity

func (attachmentEntity AttachmentEntity) Labels(tenant string) []string {
	return []string{
		model.NodeLabelAttachment,
		model.NodeLabelAttachment + "_" + tenant,
	}
}
