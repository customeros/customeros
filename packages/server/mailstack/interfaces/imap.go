package interfaces

import (
	"context"
	"time"

	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type IMAPService interface {
	Start(ctx context.Context) error
	Stop() error
	AddMailbox(ctx context.Context, mailbox *models.Mailbox) error
	RemoveMailbox(ctx context.Context, mailboxID string) error
	Status() map[string]MailboxStatus
	SetEventHandler(handler func(ctx context.Context, mailEvent MailEvent))
}

type MailboxStatus struct {
	Connected   bool
	LastError   string
	Folders     map[string]FolderStats
	LastChecked time.Time
}

type FolderStats struct {
	Total    uint32
	Unseen   uint32
	LastSeen uint32
	LastSync time.Time
}

type MailEvent struct {
	Source    string
	MailboxID string
	Folder    string
	MessageID uint32
	EventType string
	Message   interface{}
}
