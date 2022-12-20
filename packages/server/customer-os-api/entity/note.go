package entity

import (
	"fmt"
)

type NoteEntity struct {
	Id   string
	Text string
}

func (note NoteEntity) ToString() string {
	return fmt.Sprintf("id: %s\ntext: %s", note.Id, note.Text)
}

type NoteEntities []NoteEntity
