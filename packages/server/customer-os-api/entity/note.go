package entity

import (
	"fmt"
	"time"
)

type NoteEntity struct {
	Id            string
	Html          string
	CreatedAt     *time.Time
	Source        DataSource
	SourceOfTruth DataSource
}

func (note NoteEntity) ToString() string {
	return fmt.Sprintf("id: %s\nhtml: %s", note.Id, note.Html)
}

type NoteEntities []NoteEntity

func (note NoteEntity) Labels() []string {
	return []string{"Note"}
}
