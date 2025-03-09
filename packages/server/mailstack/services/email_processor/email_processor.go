package email_processor

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/services/email_processor/handlers"
)

type Processor struct {
	imapHandler *handlers.IMAPHandler
	// Add other handlers as needed
}

func NewProcessor(eventService *events.EventsService) *Processor {
	return &Processor{
		imapHandler: handlers.NewIMAPHandler(eventService),
	}
}

// ProcessMailEvent is the main entry point for processing all mail events
func (p *Processor) ProcessMailEvent(event interfaces.MailEvent) {
	// Determine source and route accordingly
	switch event.Source {
	case "imap":
		p.imapHandler.Handle(event)
	// Add cases for other sources as you implement them
	default:
		// Generic handling
	}
}
