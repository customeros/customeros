package interfaces

import (
	"context"

	"github.com/lib/pq"
)

type IMAPService interface {
	Start(ctx context.Context) error
	Stop() error
	AddMailbox(config MailboxConfig) error
	RemoveMailbox(mailboxID string) error
	Status() map[string]MailboxStatus
	SetEventHandler(handler func(event MailEvent))
	// DB calls
	// GetLastSeenUID(mailboxID, folderName string) (uint32, error)
	// UpdateLastSeenUID(mailboxID, folderName string, uid uint32) error
}

type MailboxConfig struct {
	ID       string
	Server   string
	Port     int
	Username string
	Password string
	Folders  pq.StringArray
	TLS      bool
}

type MailboxStatus struct {
	Connected bool
	LastError string
	Folders   map[string]FolderStats
}

type FolderStats struct {
	Total    uint32
	Unseen   uint32
	LastSeen uint32
}

type MailEvent struct {
	Source    string
	MailboxID string
	Folder    string
	MessageID uint32
	EventType string
	Message   interface{}
}
