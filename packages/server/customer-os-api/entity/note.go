package entity

import (
	"fmt"
)

type NoteEntity struct {
	Id   string
	Html string
}

func (note NoteEntity) ToString() string {
	return fmt.Sprintf("id: %s\nhtml: %s", note.Id, note.Html)
}

type NoteEntities []NoteEntity
