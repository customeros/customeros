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

	// Keep track of the highest message sequence number we've seen
	lastSeqNum := mbox.Messages

	// Set up updates channel
	updates := make(chan client.Update, 100)
	c.Updates = updates

	// Start IDLE
	idleChan := make(chan error, 1)
	idleCmd, err := c.IdleWithFallback(idleChan, s.ctx)
	if err != nil {
		c.Updates = nil
		return fmt.Errorf("failed to start IDLE on folder %s: %w", folderName, err)
	}

	// Process updates
	for {
		select {
		case update := <-updates:
			switch u := update.(type) {
			case *client.MailboxUpdate:
				// If we have new messages
				if u.Mailbox.Messages > lastSeqNum {
					// Fetch new messages
					err := s.fetchNewMessages(mailboxID, c, folderName, lastSeqNum+1, u.Mailbox.Messages)
					if err != nil {
						log.Printf("Error fetching new messages: %v", err)
					}
					lastSeqNum = u.Mailbox.Messages
				}

				// Update folder stats
				s.updateFolderStats(mailboxID, folderName, u.Mailbox)
			}

		case err := <-idleChan:
			if err != nil {
				log.Printf("IDLE command error for mailbox %s, folder %s: %v", mailboxID, folderName, err)
				c.Updates = nil
				return err
			}

			// Restart IDLE
			idleCmd, err = c.IdleWithFallback(idleChan, s.ctx)
			if err != nil {
				c.Updates = nil
				return fmt.Errorf("failed to restart IDLE: %w", err)
			}

		case <-s.ctx.Done():
			// Cancel IDLE command
			idleCmd.Done()
			c.Updates = nil
			return nil
		}
	}
}

func (s *IMAPService) updateFolderStats(mailboxID, folderName string, mbox *interfaces.MailboxStatus) {
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

func (s *IMAPService) fetchNewMessages(mailboxID string, c *client.Client, folderName string, from, to uint32) error {
	if from > to {
		return nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	// Items to fetch
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchBodyStructure, "BODY.PEEK[]"}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

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
