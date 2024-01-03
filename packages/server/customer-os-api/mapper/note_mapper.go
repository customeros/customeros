package mapper

import (
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapNoteInputToEntity(input *model.NoteInput) *entity.NoteEntity {
	if input == nil {
		return nil
	}
	noteEntity := entity.NoteEntity{
		Content:       utils.IfNotNilString(input.Content),
		ContentType:   utils.IfNotNilString(input.ContentType),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	return &noteEntity
}

func MapNoteUpdateInputToEntity(input *model.NoteUpdateInput) *entity.NoteEntity {
	if input == nil {
		return nil
	}
	noteEntity := entity.NoteEntity{
		Id:            input.ID,
		Content:       utils.IfNotNilString(input.Content),
		ContentType:   utils.IfNotNilString(input.ContentType),
		SourceOfTruth: neo4jentity.DataSourceOpenline,
	}
	return &noteEntity
}

func MapEntityToNote(entity *entity.NoteEntity) *model.Note {
	note := model.Note{
		ID:            entity.Id,
		Content:       utils.StringPtr(entity.Content),
		ContentType:   utils.StringPtr(entity.ContentType),
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
	return &note
}

func MapEntitiesToNotes(entities *entity.NoteEntities) []*model.Note {
	var notes []*model.Note
	for _, noteEntity := range *entities {
		notes = append(notes, MapEntityToNote(&noteEntity))
	}
	return notes
}
