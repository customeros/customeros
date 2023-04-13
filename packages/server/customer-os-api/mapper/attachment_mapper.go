package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToAttachment(entity *entity.AttachmentEntity) *model.Attachment {
	return &model.Attachment{
		ID:        entity.Id,
		CreatedAt: *entity.CreatedAt,
		MimeType:  entity.MimeType,
		Size:      entity.Size,
		Name:      entity.Name,
		Extension: entity.Extension,

		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
}

func MapEntitiesToAttachment(entities *entity.AttachmentEntities) []*model.Attachment {
	var attachments []*model.Attachment
	for _, attachmentEntity := range *entities {
		attachments = append(attachments, MapEntityToAttachment(&attachmentEntity))
	}
	return attachments
}

func MapAttachmentInputToEntity(input *model.AttachmentInput) *entity.AttachmentEntity {
	return &entity.AttachmentEntity{
		MimeType:  input.MimeType,
		Size:      input.Size,
		Name:      input.Name,
		Extension: input.Extension,

		AppSource: input.AppSource,
	}
}
