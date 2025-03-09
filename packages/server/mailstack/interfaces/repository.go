package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

type MessageStateRepository interface {
	GetLastSeenUID(mailboxID, folderName string) (uint32, error)
	UpdateLastSeenUID(mailboxID, folderName string, uid uint32) error
}

type MailboxRepository interface {
	GetMailboxes() ([]MailboxConfig, error)
	GetMailbox(id string) (MailboxConfig, error)
	SaveMailbox(mailbox MailboxConfig) error
	DeleteMailbox(id string) error
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
