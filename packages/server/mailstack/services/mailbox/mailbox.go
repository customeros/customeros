package mailbox

import (
	"context"
	"log"
	"sync"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

type MailService struct {
	imapService interfaces.MailboxService
	handlers    []func(interfaces.MailEvent)
	handlerMu   sync.RWMutex
}

func NewMailService(imapService interfaces.MailboxService) *MailService {
	service := &MailService{
		imapService: imapService,
		handlers:    make([]func(interfaces.MailEvent), 0),
	}

	// Set event handler on IMAP service
	imapService.SetEventHandler(service.handleMailEvent)

	return service
}

func (s *MailService) Start(ctx context.Context) error {
	return s.imapService.Start(ctx)
}

func (s *MailService) Stop() error {
	return s.imapService.Stop()
}

func (s *MailService) AddMailbox(config interfaces.MailboxConfig) error {
	return s.imapService.AddMailbox(config)
}

func (s *MailService) RemoveMailbox(mailboxID string) error {
	return s.imapService.RemoveMailbox(mailboxID)
}

func (s *MailService) AddEventHandler(handler func(interfaces.MailEvent)) {
	s.handlerMu.Lock()
	defer s.handlerMu.Unlock()
	s.handlers = append(s.handlers, handler)
}

func (s *MailService) handleMailEvent(event interfaces.MailEvent) {
	s.handlerMu.RLock()
	defer s.handlerMu.RUnlock()

	for _, handler := range s.handlers {
		// Execute each handler in a goroutine
		go func(h func(interfaces.MailEvent), e interfaces.MailEvent) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic in event handler: %v", r)
				}
			}()
			h(e)
		}(handler, event)
	}
}

func (s *MailService) Status() map[string]interfaces.MailboxStatus {
	return s.imapService.Status()
}
