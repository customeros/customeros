package imap

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/emersion/go-imap/client"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

// connectMailbox establishes a connection to an IMAP server for a given mailbox configuration
func (s *IMAPService) connectMailbox(ctx context.Context, config *models.Mailbox) (*client.Client, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.connectMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", config.ID)
	span.SetTag("server", config.ImapServer)
	span.SetTag("port", config.ImapPort)
	span.SetTag("tls", config.ImapTLS)

	// Format server address with port
	serverAddr := fmt.Sprintf("%s:%d", config.ImapServer, config.ImapPort)

	// Set up connection with timeout
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	var c *client.Client
	var err error

	if config.ImapTLS {
		// Create custom TLS config
		tlsConfig := &tls.Config{
			ServerName: config.ImapServer,
		}

		c, err = client.DialWithDialerTLS(dialer, serverAddr, tlsConfig)
	} else {
		c, err = client.DialWithDialer(dialer, serverAddr)
	}

	if err != nil {
		span.SetTag("error", true)
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to connect to %s: %w", serverAddr, err)
	}

	// Check server capabilities
	caps, err := c.Capability()
	if err != nil {
		// Close the connection before returning
		c.Logout()

		span.SetTag("error", true)
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get capabilities: %w", err)
	}

	// Log server capabilities for debugging
	log.Printf("[%s] Server capabilities: %v", config.ID, caps)
	span.SetTag("server.capabilities", fmt.Sprintf("%v", caps))

	// Add login tracing
	loginSpan := opentracing.StartSpan(
		"IMAPService.login",
		opentracing.ChildOf(span.Context()),
	)
	loginSpan.SetTag("username", config.ImapUsername)

	// Set client timeout for login
	c.Timeout = 30 * time.Second

	// Perform login
	err = c.Login(config.ImapUsername, config.ImapPassword)
	if err != nil {
		// Close the connection before returning
		c.Logout()

		loginSpan.SetTag("error", true)
		tracing.TraceErr(loginSpan, err)
		loginSpan.Finish()

		return nil, fmt.Errorf("failed to login as %s: %w", config.ImapUsername, err)
	}

	loginSpan.SetTag("success", true)
	loginSpan.Finish()

	// Reset client timeout to default
	c.Timeout = 0 // No timeout for normal operations

	// Log successful connection
	log.Printf("[%s] Successfully connected and logged in to %s", config.ID, serverAddr)
	span.SetTag("success", true)

	return c, nil
}

// reconnectMailbox attempts to reconnect to a mailbox with exponential backoff
func (s *IMAPService) reconnectMailbox(ctx context.Context, mailboxID string) (*client.Client, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.reconnectMailbox")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)

	// Get the config for this mailbox
	s.clientsMutex.RLock()
	config, exists := s.configs[mailboxID]
	s.clientsMutex.RUnlock()

	if !exists {
		err := fmt.Errorf("no configuration found for mailbox %s", mailboxID)
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Try to connect
	c, err := s.connectMailbox(ctx, config)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Store the new client
	s.clientsMutex.Lock()
	s.clients[mailboxID] = c
	s.clientsMutex.Unlock()

	// Update status
	s.updateStatus(mailboxID, interfaces.MailboxStatus{
		Connected: true,
		Folders:   make(map[string]interfaces.FolderStats),
	})

	return c, nil
}

// getClient gets an existing client or creates a new one if needed
func (s *IMAPService) getClient(ctx context.Context, mailboxID string) (*client.Client, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.getClient")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)

	// Try to get an existing client
	s.clientsMutex.RLock()
	c, exists := s.clients[mailboxID]
	s.clientsMutex.RUnlock()

	// If client exists, verify it's still working
	if exists {
		// Perform a NOOP to check connection
		err := c.Noop()
		if err == nil {
			return c, nil
		}

		// Connection is broken, remove it
		log.Printf("[%s] Existing connection is broken: %v", mailboxID, err)
		s.clientsMutex.Lock()
		delete(s.clients, mailboxID)
		s.clientsMutex.Unlock()
	}

	// Need to create a new connection
	return s.reconnectMailbox(ctx, mailboxID)
}

// disconnectClient safely disconnects an IMAP client
func (s *IMAPService) disconnectClient(mailboxID string, c *client.Client) {
	span := opentracing.StartSpan("IMAPService.disconnectClient")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)
	span.SetTag("mailbox.id", mailboxID)

	if c == nil {
		return
	}

	// Create a timeout context for logout
	logoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Set client timeout for logout
	c.Timeout = 5 * time.Second

	// Create a channel to signal logout completion
	done := make(chan error, 1)

	// Perform logout in a goroutine
	go func() {
		done <- c.Logout()
		close(done)
	}()

	// Wait for logout or timeout
	select {
	case err := <-done:
		if err != nil {
			log.Printf("[%s] Error during logout: %v", mailboxID, err)
			span.SetTag("error", true)
			tracing.TraceErr(span, err)
		} else {
			log.Printf("[%s] Successfully logged out", mailboxID)
		}
	case <-logoutCtx.Done():
		log.Printf("[%s] Logout timed out", mailboxID)
		span.SetTag("timeout", true)
	}

	// Remove from clients map
	s.clientsMutex.Lock()
	delete(s.clients, mailboxID)
	s.clientsMutex.Unlock()

	// Update status
	s.statusMutex.Lock()
	if status, ok := s.statuses[mailboxID]; ok {
		status.Connected = false
		s.statuses[mailboxID] = status
	}
	s.statusMutex.Unlock()
}

// disconnectAllClients safely disconnects all IMAP clients
func (s *IMAPService) disconnectAllClients() {
	span := opentracing.StartSpan("IMAPService.disconnectAllClients")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)

	// Get a copy of the clients map
	s.clientsMutex.Lock()
	clients := make(map[string]*client.Client)
	for id, c := range s.clients {
		clients[id] = c
		delete(s.clients, id)
	}
	s.clientsMutex.Unlock()

	// Disconnect each client
	for id, c := range clients {
		s.disconnectClient(id, c)
	}
}
