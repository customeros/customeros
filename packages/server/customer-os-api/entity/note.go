package entity

import (
	"fmt"
	"time"
)

type NoteEntity struct {
	Id        string
	Html      string
	CreatedAt *time.Time
}

func (note NoteEntity) ToString() string {
	return fmt.Sprintf("id: %s\nhtml: %s", note.Id, note.Html)
}

type NoteEntities []NoteEntity
