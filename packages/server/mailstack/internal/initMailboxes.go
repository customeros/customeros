package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services"
)

// InitMailboxes initializes all mailbox connections from configuration
func InitMailboxes(s *services.Services, r *repository.Repositories) error {
	log.Println("Initializing mailbox connections...")
	ctx := context.Background()

	// get mailboxes from database
	mailboxes, err := r.MailboxRepository.GetMailboxes(ctx)
	if err != nil {
		return err
	}

	// Add each mailbox from configuration
	for _, mailbox := range mailboxes {
		if err := s.IMAPService.AddMailbox(ctx, mailbox); err != nil {
			return fmt.Errorf("failed to add mailbox %s: %w", mailbox.ID, err)
		}
	}

	log.Printf("Successfully initialized %d mailboxes", len(mailboxes))
	return nil
}
