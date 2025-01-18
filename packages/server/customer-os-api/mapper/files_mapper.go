package mapper

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
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
		MetadataUrl: fmt.Sprintf("%s/file/%s", serviceUrl, input.ID),
		DownloadUrl: fmt.Sprintf("%s/file/%s/download", serviceUrl, input.ID),
		CdnUrl:      input.CdnUrl,
	}
	return &file
}
