package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type AttachmentEntity struct {
	Id        string
	CreatedAt *time.Time

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
		neo4jutil.NodeLabelAttachment,
		neo4jutil.NodeLabelAttachment + "_" + tenant,
	}
}
