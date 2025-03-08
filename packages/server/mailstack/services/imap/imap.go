package imap

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

type IMAPService struct {
	clients      map[string]*client.Client
	configs      map[string]interfaces.MailboxConfig
	eventHandler func(interfaces.MailEvent)
	clientsMutex sync.RWMutex
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	statuses     map[string]interfaces.MailboxStatus
	statusMutex  sync.RWMutex
}

func NewIMAPService() *IMAPService {
	return &IMAPService{
		clients:  make(map[string]*client.Client),
		configs:  make(map[string]interfaces.MailboxConfig),
		statuses: make(map[string]interfaces.MailboxStatus),
	}
}

const (
	DEFAULT_IMAP_LOGOUT    = 25 // minutes
	DEFAULT_POLLING_PERIOD = 20 // minutes
)

func (s *IMAPService) SetEventHandler(handler func(interfaces.MailEvent)) {
	s.eventHandler = handler
}

func (s *IMAPService) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)

	// start monitoring mailboxes
	for id, config := range s.configs {
		go s.monitorMailbox(id, config)
	}

	go s.runHealthChecks()

	return nil
}

func (s *IMAPService) Stop() error {
	if s.cancel != nil {
		s.cancel()
	}

	// Close all connections
	s.clientsMutex.Lock()
	for id, client := range s.clients {
		client.Logout()
		delete(s.clients, id)
	}
	s.clientsMutex.Unlock()

	s.wg.Wait()
	return nil
}

func (s *IMAPService) AddMailbox(config interfaces.MailboxConfig) error {
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
		go s.monitorMailbox(config.ID, config)
	}

	return nil
}

func (s *IMAPService) RemoveMailbox(mailboxID string) error {
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

func (s *IMAPService) connectMailbox(config interfaces.MailboxConfig) (*client.Client, error) {
	// Format server address with port
	serverAddr := fmt.Sprintf("%s:%d", config.Server, config.Port)

	var c *client.Client
	var err error

	if config.TLS {
		c, err = client.DialTLS(serverAddr, nil)
	} else {
		c, err = client.Dial(serverAddr)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", serverAddr, err)
	}

	// Login
	if err := c.Login(config.Username, config.Password); err != nil {
		c.Logout()
		return nil, fmt.Errorf("failed to login as %s: %w", config.Username, err)
	}

	return c, nil
}

func (s *IMAPService) monitorMailbox(mailboxID string, config interfaces.MailboxConfig) {
	s.wg.Add(1)
	defer s.wg.Done()

	backoff := time.Second
	maxBackoff := time.Minute * 5

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			// Try to connect
			c, err := s.connectMailbox(config)
			if err != nil {
				log.Printf("Error connecting to mailbox %s: %v", mailboxID, err)
				s.updateStatusError(mailboxID, err)

				// Backoff before retrying
				select {
				case <-time.After(backoff):
					backoff = min(backoff*2, maxBackoff)
				case <-s.ctx.Done():
					return
				}
				continue
			}

			// Reset backoff on successful connection
			backoff = time.Second

			// Store client
			s.clientsMutex.Lock()
			s.clients[mailboxID] = c
			s.clientsMutex.Unlock()

			// Update status
			s.updateStatus(mailboxID, interfaces.MailboxStatus{
				Connected: true,
				Folders:   make(map[string]interfaces.FolderStats),
			})

			// Monitor folders
			for _, folder := range config.Folders {
				if err := s.monitorFolder(mailboxID, c, folder); err != nil {
					log.Printf("Error monitoring folder %s: %v", folder, err)
				}
			}

			// If we're here, the connection was lost
			s.clientsMutex.Lock()
			delete(s.clients, mailboxID)
			s.clientsMutex.Unlock()
		}
	}
}

func (s *IMAPService) monitorFolder(mailboxID string, c *client.Client, folderName string) error {
	// Select the mailbox (folder)
	mbox, err := c.Select(folderName, false)
	if err != nil {
		return fmt.Errorf("failed to select folder %s: %w", folderName, err)
	}

	// Update folder stats
	s.updateFolderStats(mailboxID, folderName, mbox)

	// Get the last seen UID from persistent storage
	lastSeenUID, err := s.tracker.GetLastSeenUID(mailboxID, folderName)
	if err != nil {
		log.Printf("Warning: Could not get last seen UID: %v, starting from current state", err)
		lastSeenUID = 0
	}

	// If this is a new folder or we don't have history, use current state
	if lastSeenUID == 0 {
		// Store the current highest UID
		if mbox.UidNext > 1 {
			lastSeenUID = mbox.UidNext - 1
			err = s.tracker.UpdateLastSeenUID(mailboxID, folderName, lastSeenUID)
			if err != nil {
				log.Printf("Warning: Failed to update last seen UID: %v", err)
			}
		}
	} else {
		// Fetch any messages that arrived while we were disconnected
		if mbox.UidNext > lastSeenUID+1 {
			s.fetchNewMessagesByUID(mailboxID, c, folderName, lastSeenUID+1, mbox.UidNext-1)
		}
	}

	// Set up updates channel
	updates := make(chan client.Update, 100)
	c.Updates = updates

	// Create a stop channel that will be closed when we need to stop IDLE
	stop := make(chan struct{})

	// Start a goroutine to handle context cancellation
	go func() {
		<-s.ctx.Done()
		close(stop)
	}()

	// Start IDLE with proper timeout handling
	err = c.Idle(stop, &client.IdleOptions{
		LogoutTimeout: DEFAULT_IMAP_LOGOUT,
		PollInterval:  DEFAULT_POLLING_PERIOD,
	})

	// If we get here, either there was an error or the context was canceled
	c.Updates = nil

	if err != nil && s.ctx.Err() == nil {
		// There was an error and it wasn't due to context cancellation
		return fmt.Errorf("IDLE error: %w", err)
	}

	return nil
}

func (s *IMAPService) updateFolderStats(mailboxID, folderName string, mbox *imap.MailboxStatus) {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	status, exists := s.statuses[mailboxID]
	if !exists {
		status = interfaces.MailboxStatus{
			Connected: true,
			Folders:   make(map[string]interfaces.FolderStats),
		}
	}

	// Count unseen messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	s.clientsMutex.RLock()
	client := s.clients[mailboxID]
	s.clientsMutex.RUnlock()

	var unseenCount uint32 = 0
	if client != nil {
		ids, err := client.Search(criteria)
		if err == nil {
			unseenCount = uint32(len(ids))
		}
	}

	status.Folders[folderName] = interfaces.FolderStats{
		Total:    mbox.Messages,
		Unseen:   unseenCount,
		LastSeen: mbox.Messages,
	}

	s.statuses[mailboxID] = status
}

func (s *IMAPService) fetchNewMessagesByUID(mailboxID string, c *client.Client, folderName string, fromUID, toUID uint32) error {
	if fromUID > toUID {
		return nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(fromUID, toUID)

	// Items to fetch
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchBodyStructure, "BODY.PEEK[]", imap.FetchUid}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.UidFetch(seqSet, items, messages)
	}()

	var highestUID uint32 = 0

	for msg := range messages {
		if s.eventHandler != nil {
			s.eventHandler(interfaces.MailEvent{
				MailboxID: mailboxID,
				Folder:    folderName,
				MessageID: msg.SeqNum,
				EventType: "new",
				Message:   msg,
			})
		}

		if msg.Uid > highestUID {
			highestUID = msg.Uid
		}
	}

	// Update the last seen UID
	if highestUID > 0 {
		err := s.tracker.UpdateLastSeenUID(mailboxID, folderName, highestUID)
		if err != nil {
			log.Printf("Warning: Failed to update last seen UID: %v", err)
		}
	}

	return <-done
}

func (s *IMAPService) runHealthChecks() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAllConnections()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *IMAPService) checkAllConnections() {
	s.clientsMutex.RLock()
	clients := make(map[string]*client.Client)
	for id, c := range s.clients {
		clients[id] = c
	}
	s.clientsMutex.RUnlock()

	for id, c := range clients {
		if err := c.Noop(); err != nil {
			log.Printf("Health check failed for mailbox %s: %v", id, err)

			// Update status
			s.updateStatusError(id, err)

			// Remove from clients map
			s.clientsMutex.Lock()
			delete(s.clients, id)
			s.clientsMutex.Unlock()
		}
	}
}
