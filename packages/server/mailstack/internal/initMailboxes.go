package internal

import (
	"fmt"
	"log"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/services"
)

// InitMailboxes initializes all mailbox connections from configuration
func InitMailboxes(s *services.Services) error {
	log.Println("Initializing mailbox connections...")

	mailbox := interfaces.MailboxConfig{
		ID:       "test",
		Server:   "mail.hostedemail.com",
		Port:     993,
		Username: "test@testcustomeros.com",
		Password: "admin123!",
		Folders:  []string{"INBOX"},
		TLS:      true,
	}
	mailboxes := []interfaces.MailboxConfig{mailbox}

	// Add each mailbox from configuration
	for _, mbConfig := range mailboxes {
		log.Printf("Adding mailbox: %s (%s)", mbConfig.ID, mbConfig.Username)

		if err := s.IMAPService.AddMailbox(mbConfig); err != nil {
			return fmt.Errorf("failed to add mailbox %s: %w", mbConfig.ID, err)
		}
	}

	log.Printf("Successfully initialized %d mailboxes", len(mailboxes))
	return nil
}
