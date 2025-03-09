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

func NewIMAPService() interfaces.IMAPService {
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
				if err := s.monitorFolder(mailboxID, c, string(folder)); err != nil {
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
	log.Printf("[%s][%s] Starting to monitor folder", mailboxID, folderName)

	// Select the mailbox (folder)
	mbox, err := c.Select(folderName, false)
	if err != nil {
		log.Printf("[%s][%s] Error selecting folder: %v", mailboxID, folderName, err)
		return fmt.Errorf("failed to select folder %s: %w", folderName, err)
	}

	log.Printf("[%s][%s] Selected folder - Messages: %d, Recent: %d, Unseen: %d",
		mailboxID, folderName, mbox.Messages, mbox.Recent, mbox.Unseen)

	// Check server capabilities
	caps, err := c.Capability()
	if err != nil {
		log.Printf("[%s][%s] Error getting capabilities: %v", mailboxID, folderName, err)
	} else {
		log.Printf("[%s][%s] Server capabilities: %v", mailboxID, folderName, caps)
	}

	// Update folder stats
	s.updateFolderStats(mailboxID, folderName, mbox)

	// Check for recent messages at startup
	if mbox.Recent > 0 {
		log.Printf("[%s][%s] Found %d recent messages, fetching them", mailboxID, folderName, mbox.Recent)

		// Fetch all recent messages
		criteria := imap.NewSearchCriteria()
		criteria.WithFlags = []string{imap.RecentFlag}

		uids, err := c.Search(criteria)
		if err != nil {
			log.Printf("[%s][%s] Error searching for recent messages: %v", mailboxID, folderName, err)
		} else if len(uids) > 0 {
			log.Printf("[%s][%s] Found %d UIDs with recent flag", mailboxID, folderName, len(uids))

			seqSet := new(imap.SeqSet)
			for _, uid := range uids {
				seqSet.AddNum(uid)
			}

			messages := make(chan *imap.Message, 10)
			done := make(chan error, 1)
			items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchBodyStructure, "BODY.PEEK[]", imap.FetchUid}

			go func() {
				done <- c.Fetch(seqSet, items, messages)
			}()

			for msg := range messages {
				log.Printf("[%s][%s] Processing recent message: %d", mailboxID, folderName, msg.SeqNum)
				if s.eventHandler != nil {
					s.eventHandler(interfaces.MailEvent{
						Source:    "imap",
						MailboxID: mailboxID,
						Folder:    folderName,
						MessageID: msg.SeqNum,
						EventType: "new",
						Message:   msg,
					})
				}
			}

			err = <-done
			if err != nil {
				log.Printf("[%s][%s] Error fetching recent messages: %v", mailboxID, folderName, err)
			}
		}
	}

	// Check for unseen messages as well
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	uids, err := c.Search(criteria)
	if err != nil {
		log.Printf("[%s][%s] Error searching for unseen messages: %v", mailboxID, folderName, err)
	} else {
		log.Printf("[%s][%s] Found %d unseen messages", mailboxID, folderName, len(uids))

		if len(uids) > 0 {
			// Get the latest few unseen messages
			maxUnseen := 5 // Limit to avoid processing too many messages at startup
			if len(uids) < maxUnseen {
				maxUnseen = len(uids)
			}

			// Take the most recent N unseen messages
			recentUids := uids[len(uids)-maxUnseen:]

			seqSet := new(imap.SeqSet)
			for _, uid := range recentUids {
				seqSet.AddNum(uid)
			}

			messages := make(chan *imap.Message, 10)
			done := make(chan error, 1)
			items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchBodyStructure, "BODY.PEEK[]", imap.FetchUid}

			go func() {
				done <- c.Fetch(seqSet, items, messages)
			}()

			for msg := range messages {
				log.Printf("[%s][%s] Processing unseen message: %d", mailboxID, folderName, msg.SeqNum)
				if s.eventHandler != nil {
					s.eventHandler(interfaces.MailEvent{
						Source:    "imap",
						MailboxID: mailboxID,
						Folder:    folderName,
						MessageID: msg.SeqNum,
						EventType: "new",
						Message:   msg,
					})
				}
			}

			err = <-done
			if err != nil {
				log.Printf("[%s][%s] Error fetching unseen messages: %v", mailboxID, folderName, err)
			}
		}
	}

	// Remember initial message count
	initialCount := mbox.Messages
	log.Printf("[%s][%s] Initial message count: %d", mailboxID, folderName, initialCount)

	// Set up updates channel
	updates := make(chan client.Update, 100)
	c.Updates = updates
	log.Printf("[%s][%s] Set up updates channel", mailboxID, folderName)

	// Test the updates channel
	log.Printf("[%s][%s] Testing updates channel...", mailboxID, folderName)
	select {
	case update := <-updates:
		log.Printf("[%s][%s] Got an immediate update on the channel: %T", mailboxID, folderName, update)
	default:
		log.Printf("[%s][%s] No immediate updates available", mailboxID, folderName)
	}

	// Check IDLE support
	supported, err := c.Support("IDLE")
	if err != nil {
		log.Printf("[%s][%s] Error checking IDLE support: %v", mailboxID, folderName, err)
	}
	log.Printf("[%s][%s] IDLE support: %v", mailboxID, folderName, supported)

	// Create a stop channel that will be closed when we need to stop IDLE
	stop := make(chan struct{})

	// Start a goroutine to handle context cancellation
	go func() {
		<-s.ctx.Done()
		log.Printf("[%s][%s] Context cancelled, stopping IDLE", mailboxID, folderName)
		close(stop)
	}()

	// Start a goroutine to send NOOPs periodically
	go func() {
		noopTicker := time.NewTicker(30 * time.Second)
		defer noopTicker.Stop()

		for {
			select {
			case <-noopTicker.C:
				err := c.Noop()
				if err != nil {
					log.Printf("[%s][%s] Error during NOOP: %v", mailboxID, folderName, err)
					return
				}
				log.Printf("[%s][%s] Sent NOOP command", mailboxID, folderName)
			case <-s.ctx.Done():
				return
			}
		}
	}()

	// Process updates while IDLE is running
	updateProcessor := make(chan struct{})
	go func() {
		defer close(updateProcessor)

		for {
			select {
			case update, ok := <-updates:
				if !ok {
					log.Printf("[%s][%s] Updates channel closed", mailboxID, folderName)
					return
				}

				log.Printf("[%s][%s] Received update: %T", mailboxID, folderName, update)

				switch u := update.(type) {
				case *client.MailboxUpdate:
					log.Printf("[%s][%s] Mailbox update - Messages: %d (was: %d)",
						mailboxID, folderName, u.Mailbox.Messages, initialCount)

					// If we have new messages
					if u.Mailbox.Messages > initialCount {
						newMessages := u.Mailbox.Messages - initialCount
						log.Printf("[%s][%s] 📥 Detected %d new message(s)", mailboxID, folderName, newMessages)

						// Fetch new messages
						err := s.fetchNewMessages(mailboxID, c, folderName, initialCount+1, u.Mailbox.Messages)
						if err != nil {
							log.Printf("[%s][%s] Error fetching new messages: %v", mailboxID, folderName, err)
						}
						initialCount = u.Mailbox.Messages
					}

					// Update folder stats
					s.updateFolderStats(mailboxID, folderName, u.Mailbox)

				case *client.ExpungeUpdate:
					log.Printf("[%s][%s] Message expunged: %d", mailboxID, folderName, u.SeqNum)
					if u.SeqNum <= initialCount {
						initialCount--
					}

				case *client.MessageUpdate:
					log.Printf("[%s][%s] Message updated: %v", mailboxID, folderName, u.Message)
					// Instead of u.SeqNum, we should use u.Message.SeqNum if available
					if u.Message != nil {
						log.Printf("[%s][%s] Message updated, SeqNum: %d", mailboxID, folderName, u.Message.SeqNum)
					}

				default:
					log.Printf("[%s][%s] Received update of unknown type: %T", mailboxID, folderName, update)
				}

			case <-s.ctx.Done():
				log.Printf("[%s][%s] Context cancelled in update processor", mailboxID, folderName)
				return
			}
		}
	}()

	// Start IDLE with proper timeout handling
	log.Printf("[%s][%s] Starting IDLE command", mailboxID, folderName)
	err = c.Idle(stop, &client.IdleOptions{
		LogoutTimeout: time.Duration(DEFAULT_IMAP_LOGOUT) * time.Minute,
		PollInterval:  time.Duration(DEFAULT_POLLING_PERIOD) * time.Minute,
	})

	// Wait for update processor to finish
	<-updateProcessor

	// If we get here, either there was an error or the context was canceled
	c.Updates = nil

	if err != nil && s.ctx.Err() == nil {
		// There was an error and it wasn't due to context cancellation
		log.Printf("[%s][%s] IDLE error: %v", mailboxID, folderName, err)
		return fmt.Errorf("IDLE error: %w", err)
	}

	log.Printf("[%s][%s] Stopped monitoring folder", mailboxID, folderName)
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

func (s *IMAPService) fetchNewMessages(mailboxID string, c *client.Client, folderName string, from, to uint32) error {
	if from > to {
		return nil
	}

	log.Printf("[%s][%s] Fetching messages from sequence %d to %d", mailboxID, folderName, from, to)

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	// Items to fetch
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchBodyStructure, "BODY.PEEK[]", imap.FetchUid}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	for msg := range messages {
		log.Printf("[%s][%s] Received message: UID=%d, Seq=%d, Subject=%s",
			mailboxID, folderName, msg.Uid, msg.SeqNum, msg.Envelope.Subject)

		if s.eventHandler != nil {
			log.Printf("[%s][%s] Triggering event handler for message %d", mailboxID, folderName, msg.SeqNum)
			s.eventHandler(interfaces.MailEvent{
				Source:    "imap",
				MailboxID: mailboxID,
				Folder:    folderName,
				MessageID: msg.SeqNum,
				EventType: "new",
				Message:   msg,
			})
		} else {
			log.Printf("[%s][%s] Warning: No event handler registered", mailboxID, folderName)
		}
	}

	err := <-done
	if err != nil {
		log.Printf("[%s][%s] Error fetching messages: %v", mailboxID, folderName, err)
		return err
	}

	log.Printf("[%s][%s] Successfully fetched messages", mailboxID, folderName)
	return nil
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
				Source:    "imap",
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
	// if highestUID > 0 {
	// 	err := s.tracker.UpdateLastSeenUID(mailboxID, folderName, highestUID)
	// 	if err != nil {
	// 		log.Printf("Warning: Failed to update last seen UID: %v", err)
	// 	}
	// }

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
