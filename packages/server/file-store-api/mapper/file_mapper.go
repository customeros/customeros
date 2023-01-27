package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/repository/entity"
)

func MapFileEntityToDTO(input *entity.File, serviceUrl string) *dto.File {
	if input == nil {
		return nil
	}
	file := dto.File{
		Id:          input.ID,
		Name:        input.Name,
		Extension:   input.Extension,
		Mime:        input.MIME,
		Length:      input.Length,
		MetadataUrl: fmt.Sprintf("%s/file/%s", serviceUrl, input.ID),
		DownloadUrl: fmt.Sprintf("%s/file/%s/download", serviceUrl, input.ID),
	}
	return &file
}
