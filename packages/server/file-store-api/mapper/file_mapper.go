package mapper

import (
	"fmt"
	fs "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
)

func MapFileEntityToDTO(input *model.File, serviceUrl string) *fs.FileDTO {
	if input == nil {
		return nil
	}
	file := fs.FileDTO{
		Id:          input.ID,
		FileName:    input.FileName,
		MimeType:    input.MimeType,
		Size:        input.Size,
		MetadataUrl: fmt.Sprintf("%s/file/%s", serviceUrl, input.ID),
		DownloadUrl: fmt.Sprintf("%s/file/%s/download", serviceUrl, input.ID),
		CdnUrl:      input.CdnUrl,
	}
	return &file
}

func MapAttachmentResponseToFileEntity(input *model.Attachment) *model.File {
	if input == nil {
		return nil
	}
	return &model.File{
		ID:       input.Id,
		FileName: input.FileName,
		MimeType: input.MimeType,
		BasePath: input.BasePath,
		Size:     input.Size,
		CdnUrl:   input.CdnUrl,
	}
}
