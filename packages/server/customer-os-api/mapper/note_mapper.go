package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapNoteInputToEntity(input *model.NoteInput) *entity.NoteEntity {
	if input == nil {
		return nil
	}
	noteEntity := entity.NoteEntity{
		Text: input.Text,
	}
	return &noteEntity
}

func MapNoteUpdateInputToEntity(input *model.NoteUpdateInput) *entity.NoteEntity {
	if input == nil {
		return nil
	}
	emailEntity := entity.NoteEntity{
		Id:   input.ID,
		Text: input.Text,
	}
	return &emailEntity
}

func MapEntityToNote(entity *entity.NoteEntity) *model.Note {
	return &model.Note{
		ID:   entity.Id,
		Text: entity.Text,
	}
}

func MapEntitiesToNotes(entities *entity.NoteEntities) []*model.Note {
	var notes []*model.Note
	for _, noteEntity := range *entities {
		notes = append(notes, MapEntityToNote(&noteEntity))
	}
	return notes
}
