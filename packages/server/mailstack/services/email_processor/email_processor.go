package email_processor

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/services/email_processor/handlers"
)

type Processor struct {
	imapHandler *handlers.IMAPHandler
	// Add other handlers as needed
}

func NewProcessor(
	repositories *repository.Repositories,
	eventService *events.EventsService,
	emailFilterService interfaces.EmailFilterService,
) *Processor {
	return &Processor{
		imapHandler: handlers.NewIMAPHandler(repositories, eventService, emailFilterService),
	}
}

// ProcessMailEvent is the main entry point for processing all mail events
func (p *Processor) ProcessMailEvent(ctx context.Context, mailEvent interfaces.MailEvent) {
	// Determine source and route accordingly
	switch mailEvent.Source {
	case "imap":
		p.imapHandler.Handle(ctx, mailEvent)
	// Add cases for other sources as you implement them
	default:
		// Generic handling
	}
}
