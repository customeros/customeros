package mapper

import (
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
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
