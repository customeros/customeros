package imap

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/emersion/go-imap/client"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
)

type IMAPService struct {
	repositories *repository.Repositories
	clients      map[string]*client.Client
	configs      map[string]*models.Mailbox
	eventHandler func(context.Context, interfaces.MailEvent)
	clientsMutex sync.RWMutex
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	statuses     map[string]interfaces.MailboxStatus
	statusMutex  sync.RWMutex
}

func NewIMAPService(repos *repository.Repositories) interfaces.IMAPService {
	return &IMAPService{
		repositories: repos,
		clients:      make(map[string]*client.Client),
		configs:      make(map[string]*models.Mailbox),
		statuses:     make(map[string]interfaces.MailboxStatus),
	}
}

const (
	DEFAULT_IMAP_LOGOUT     = 25 // minutes
	DEFAULT_POLLING_PERIOD  = 20 // minutes
	INITIAL_SYNC_BATCH_SIZE = 100
	INITIAL_SYNC_MAX_TOTAL  = 1000
)

func (s *IMAPService) SetEventHandler(handler func(context.Context, interfaces.MailEvent)) {
	s.eventHandler = handler
}

func (s *IMAPService) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	// start monitoring mailboxes
	for id, config := range s.configs {
		go s.monitorMailbox(ctx, id, config)
	}

	go s.runHealthChecks()

	go s.runPeriodicSyncMaintenance()

	return nil
}

func (s *IMAPService) Stop() error {
	log.Println("IMAPService: Stop called, cancelling context...")

	// Cancel main context to signal all goroutines
	if s.cancel != nil {
		s.cancel()
	}

	// Create a timeout context for shutdown operations
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Close all connections with timeout
	s.clientsMutex.Lock()
	clients := make(map[string]*client.Client)
	for id, c := range s.clients {
		clients[id] = c
		delete(s.clients, id)
	}
	s.clientsMutex.Unlock()

	// Logout clients with timeout
	for id, c := range clients {
		log.Printf("IMAPService: Logging out client %s...", id)

		// Create a goroutine to handle logout
		go func(client *client.Client, clientID string) {
			err := client.Logout()
			if err != nil {
				log.Printf("Error logging out %s: %v", clientID, err)
			}
		}(c, id)
	}

	// Wait for goroutines to finish or timeout
	log.Println("IMAPService: Waiting for goroutines to finish (max 5 seconds)...")

	// Use a channel to signal completion of waitgroup
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// Wait for either waitgroup completion or timeout
	select {
	case <-done:
		log.Println("IMAPService: All goroutines finished gracefully")
	case <-shutdownCtx.Done():
		log.Println("IMAPService: Timed out waiting for goroutines, forcing exit")
	}

	log.Println("IMAPService: Stop completed")
	return nil
}

func (s *IMAPService) AddMailbox(ctx context.Context, config *models.Mailbox) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.AddMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	_, exists := s.configs[config.ID]
	if exists {
		return fmt.Errorf("mailbox with ID %s already exists", config.ID)
	}

	s.configs[config.ID] = config
	s.updateStatus(config.ID, interfaces.MailboxStatus{
		Connected: false,
		Folders:   make(map[string]interfaces.FolderStats),
	})

	// Start monitoring if service is already running
	if s.ctx != nil {
		go s.monitorMailbox(ctx, config.ID, config)
	}

	return nil
}

func (s *IMAPService) RemoveMailbox(ctx context.Context, mailboxID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.AddMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if client, exists := s.clients[mailboxID]; exists {
		client.Logout()
		delete(s.clients, mailboxID)
	}

	delete(s.configs, mailboxID)

	// Remove status
	s.statusMutex.Lock()
	delete(s.statuses, mailboxID)
	s.statusMutex.Unlock()

	return nil
}

func (s *IMAPService) Status() map[string]interfaces.MailboxStatus {
	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]interfaces.MailboxStatus)
	for id, status := range s.statuses {
		result[id] = status
	}

	return result
}

func (s *IMAPService) updateStatus(mailboxID string, status interfaces.MailboxStatus) {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()
	s.statuses[mailboxID] = status
}

func (s *IMAPService) updateStatusError(mailboxID string, err error) {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	status := s.statuses[mailboxID]
	status.Connected = false
	status.LastError = err.Error()
	s.statuses[mailboxID] = status
}

func (s *IMAPService) monitorMailbox(ctx context.Context, mailboxID string, config *models.Mailbox) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.monitorMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)

	s.wg.Add(1)
	defer s.wg.Done()

	backoff := time.Second
	maxBackoff := time.Minute * 5

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			// Create a child span for connection attempt
			connectSpan := opentracing.StartSpan(
				"IMAPService.connectMailbox",
				opentracing.ChildOf(span.Context()),
			)
			connectSpan.SetTag("mailbox.id", mailboxID)

			// Try to connect
			c, err := s.connectMailbox(ctx, config)
			if err != nil {
				err := fmt.Errorf("Error connecting to mailbox %s: %v", mailboxID, err)
				tracing.TraceErr(connectSpan, err)
				connectSpan.SetTag("error", true)
				connectSpan.Finish()

				s.updateStatusError(mailboxID, err)

				// Backoff before retrying
				backoffSpan := opentracing.StartSpan(
					"IMAPService.connectionBackoff",
					opentracing.ChildOf(span.Context()),
				)
				backoffSpan.SetTag("mailbox.id", mailboxID)
				backoffSpan.SetTag("backoff.duration_ms", backoff.Milliseconds())

				select {
				case <-time.After(backoff):
					backoff = min(backoff*2, maxBackoff)
				case <-s.ctx.Done():
					backoffSpan.Finish()
					return
				}

				backoffSpan.Finish()
				continue
			}

			connectSpan.SetTag("success", true)
			connectSpan.Finish()

			// Reset backoff on successful connection
			backoff = time.Second

			// Store client
			clientSpan := opentracing.StartSpan(
				"IMAPService.storeClient",
				opentracing.ChildOf(span.Context()),
			)

			s.clientsMutex.Lock()
			s.clients[mailboxID] = c
			s.clientsMutex.Unlock()

			clientSpan.Finish()

			// Update status
			statusSpan := opentracing.StartSpan(
				"IMAPService.updateStatus",
				opentracing.ChildOf(span.Context()),
			)

			s.updateStatus(mailboxID, interfaces.MailboxStatus{
				Connected: true,
				Folders:   make(map[string]interfaces.FolderStats),
			})

			statusSpan.Finish()

			// Monitor folders
			foldersSpan := opentracing.StartSpan(
				"IMAPService.monitorFolders",
				opentracing.ChildOf(span.Context()),
			)
			foldersSpan.SetTag("mailbox.id", mailboxID)
			foldersSpan.SetTag("folders.count", len(config.Folders))

			for _, folder := range config.Folders {
				folderSpan := opentracing.StartSpan(
					"IMAPService.monitorFolder",
					opentracing.ChildOf(foldersSpan.Context()),
				)
				folderSpan.SetTag("mailbox.id", mailboxID)
				folderSpan.SetTag("folder.name", folder)

				err := s.monitorFolder(ctx, mailboxID, c, string(folder))
				if err != nil {
					err = fmt.Errorf("Error monitoring folder %s: %v", folder, err)
					tracing.TraceErr(folderSpan, err)
					folderSpan.SetTag("error", true)
				}

				folderSpan.Finish()
			}

			foldersSpan.Finish()

			// If we're here, the connection was lost
			disconnectSpan := opentracing.StartSpan(
				"IMAPService.disconnectClient",
				opentracing.ChildOf(span.Context()),
			)
			disconnectSpan.SetTag("mailbox.id", mailboxID)

			s.clientsMutex.Lock()
			delete(s.clients, mailboxID)
			s.clientsMutex.Unlock()

			disconnectSpan.Finish()
		}
	}
}
