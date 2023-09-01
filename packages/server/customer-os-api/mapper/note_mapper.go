package mapper

import (
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
		Html:          utils.IfNotNilString(input.HTML),
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	if noteEntity.Content != "" {
		noteEntity.Html = noteEntity.Content
	} else {
		noteEntity.Content = noteEntity.Html
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
		Html:          utils.IfNotNilString(input.HTML),
		SourceOfTruth: entity.DataSourceOpenline,
	}
	if noteEntity.Content != "" {
		noteEntity.Html = noteEntity.Content
	} else {
		noteEntity.Content = noteEntity.Html
	}
	return &noteEntity
}

func MapEntityToNote(entity *entity.NoteEntity) *model.Note {
	note := model.Note{
		ID:            entity.Id,
		Content:       entity.Content,
		ContentType:   entity.ContentType,
		HTML:          entity.Html,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
	if note.Content == "" && note.HTML != "" {
		note.Content = note.HTML
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
