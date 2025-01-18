package cosapi_interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
)

type NoteService interface {
	GetById(ctx context.Context, id string) (*entity.NoteEntity, error)
	NoteLinkAttachment(ctx context.Context, noteID string, attachmentID string) error
	NoteUnlinkAttachment(ctx context.Context, noteID string, attachmentID string) error

	CreateNoteForMeeting(ctx context.Context, meetingId string, entity *entity.NoteEntity) (*entity.NoteEntity, error)
	GetNotesForMeetings(ctx context.Context, ids []string) (*entity.NoteEntities, error)

	UpdateNote(ctx context.Context, entity *entity.NoteEntity) (*entity.NoteEntity, error)
	DeleteNote(ctx context.Context, noteId string) (bool, error)
}
