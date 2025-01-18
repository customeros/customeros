package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapEntityToAttachment(entity *entity.AttachmentEntity) *model.Attachment {
	if entity == nil {
		return nil
	}
	return &model.Attachment{
		ID:        entity.Id,
		CreatedAt: utils.IfNotNilTimeWithDefault(*entity.CreatedAt, utils.Now()),
		CdnURL:    entity.CdnUrl,
		BasePath:  entity.BasePath,
		MimeType:  entity.MimeType,
		FileName:  entity.FileName,
		Size:      entity.Size,

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
		Id:        utils.IfNotNilStringWithDefault(input.ID, ""),
		CreatedAt: input.CreatedAt,
		CdnUrl:    input.CdnURL,
		BasePath:  input.BasePath,
		FileName:  input.FileName,
		MimeType:  input.MimeType,
		Size:      input.Size,

		AppSource: input.AppSource,
	}
}
