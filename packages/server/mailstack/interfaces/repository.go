package interfaces

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
