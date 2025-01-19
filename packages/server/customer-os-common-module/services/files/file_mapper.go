package files

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	neo4jEntity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

func MapAttachmentResponseToFileEntity(input *neo4jEntity.AttachmentEntity) *interfaces.File {
	if input == nil {
		return nil
	}
	return &interfaces.File{
		ID:       input.Id,
		FileName: input.FileName,
		MimeType: input.MimeType,
		BasePath: input.BasePath,
		Size:     input.Size,
		CdnUrl:   input.CdnUrl,
	}
}
