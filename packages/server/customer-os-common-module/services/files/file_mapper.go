package files

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapFileEntityToDTO(input *interfaces.File, serviceUrl string) *interfaces.FileDTO {
	if input == nil {
		return nil
	}
	file := interfaces.FileDTO{
		Id:          input.ID,
		FileName:    input.FileName,
		MimeType:    input.MimeType,
		Size:        input.Size,
		MetadataUrl: fmt.Sprintf("%s/files/%s", serviceUrl, input.ID),
		DownloadUrl: fmt.Sprintf("%s/files/%s/download", serviceUrl, input.ID),
		CdnUrl:      input.CdnUrl,
	}
	return &file
}

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
