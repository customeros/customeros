package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type MessageStateRepository interface {
	GetLastSeenUID(ctx context.Context, mailboxID, folderName string) (uint32, error)
	UpdateLastSeenUID(ctx context.Context, mailboxID, folderName string, uid uint32) error
}

type MailboxRepository interface {
	GetMailboxes(ctx context.Context) ([]*models.Mailbox, error)
	GetMailbox(ctx context.Context, id string) (*models.Mailbox, error)
	SaveMailbox(ctx context.Context, mailbox models.Mailbox) error
	DeleteMailbox(ctx context.Context, id string) error
}

type MailboxSyncRepository interface {
	GetSyncState(ctx context.Context, mailboxID, folderName string) (*models.MailboxSyncState, error)
	SaveSyncState(ctx context.Context, state *models.MailboxSyncState) error
	DeleteSyncState(ctx context.Context, mailboxID, folderName string) error
	DeleteMailboxSyncStates(ctx context.Context, mailboxID string) error
	GetAllSyncStates(ctx context.Context) (map[string]map[string]uint32, error)
	GetMailboxSyncStates(ctx context.Context, mailboxID string) (map[string]uint32, error)
}

type EmailAttachmentRepository interface {
	Create(ctx context.Context, attachment *models.EmailAttachment) error
	GetByID(ctx context.Context, id string) (*models.EmailAttachment, error)
	ListByEmail(ctx context.Context, emailID string) ([]*models.EmailAttachment, error)
	Store(ctx context.Context, attachment *models.EmailAttachment, data []byte) error
	GetData(ctx context.Context, id string) ([]byte, error)
	Delete(ctx context.Context, id string) error
}

type EmailRepository interface {
	Create(ctx context.Context, email *models.Email) error
	GetByID(ctx context.Context, id string) (*models.Email, error)
	GetByUID(ctx context.Context, mailboxID, folder string, uid uint32) (*models.Email, error)
	GetByMessageID(ctx context.Context, messageID string) (*models.Email, error)
	ListByMailbox(ctx context.Context, mailboxID string, limit, offset int) ([]*models.Email, int64, error)
	ListByFolder(ctx context.Context, mailboxID, folder string, limit, offset int) ([]*models.Email, int64, error)
	ListByThread(ctx context.Context, threadID string) ([]*models.Email, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*models.Email, int64, error)
	Update(ctx context.Context, email *models.Email) error
}
