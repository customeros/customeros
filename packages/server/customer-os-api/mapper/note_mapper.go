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
		Html:          input.HTML,
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
		Public:        utils.IfNotNilBool(input.Public),
	}
	return &noteEntity
}

func MapNoteUpdateInputToEntity(input *model.NoteUpdateInput) *entity.NoteEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.NoteEntity{
		Id:            input.ID,
		Html:          input.HTML,
		SourceOfTruth: entity.DataSourceOpenline,
		Public:        utils.IfNotNilBool(input.Public),
	}
	return &emailEntity
}

func MapEntityToNote(entity *entity.NoteEntity) *model.Note {
	return &model.Note{
		ID:            entity.Id,
		HTML:          entity.Html,
		Public:        entity.Public,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
	}
}

func MapEntitiesToNotes(entities *entity.NoteEntities) []*model.Note {
	var notes []*model.Note
	for _, noteEntity := range *entities {
		notes = append(notes, MapEntityToNote(&noteEntity))
	}
	return notes
}
