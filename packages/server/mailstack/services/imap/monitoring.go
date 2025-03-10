package imap

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

// runHealthChecks starts a background health check process
func (s *IMAPService) runHealthChecks() {
	span := opentracing.StartSpan("IMAPService.runHealthChecks")
	tracing.SetDefaultServiceSpanTags(context.Background(), span)
	span.SetTag("interval", "1m")
	span.Finish()

	log.Println("Starting IMAP health check scheduler")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAllConnections()
		case <-s.ctx.Done():
			log.Println("Stopping IMAP health checks due to context cancellation")
			return
		}
	}
}

// checkAllConnections performs a health check on all active connections
func (s *IMAPService) checkAllConnections() {
	span := opentracing.StartSpan("IMAPService.checkAllConnections")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)

	// Get a snapshot of current clients
	s.clientsMutex.RLock()
	clients := make(map[string]struct{})
	for id := range s.clients {
		clients[id] = struct{}{}
	}
	s.clientsMutex.RUnlock()

	// Track statistics for the health check
	var checkWg sync.WaitGroup
	var checkMutex sync.Mutex
	totalChecks := len(clients)
	healthyCount := 0
	unhealthyCount := 0

	span.SetTag("connections.total", totalChecks)

	// Don't check more than 5 connections at once to avoid overwhelming the system
	semaphore := make(chan struct{}, 5)

	for id := range clients {
		checkWg.Add(1)
		go func(mailboxID string) {
			defer checkWg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			isHealthy := s.checkConnection(mailboxID)

			// Update statistics
			checkMutex.Lock()
			if isHealthy {
				healthyCount++
			} else {
				unhealthyCount++
			}
			checkMutex.Unlock()
		}(id)
	}

	// Wait for all checks to complete
	checkWg.Wait()

	// Log and record metrics
	log.Printf("Health check completed - Total: %d, Healthy: %d, Unhealthy: %d",
		totalChecks, healthyCount, unhealthyCount)
	span.SetTag("connections.healthy", healthyCount)
	span.SetTag("connections.unhealthy", unhealthyCount)
}

// checkConnection checks the health of a single connection
func (s *IMAPService) checkConnection(mailboxID string) bool {
	span := opentracing.StartSpan("IMAPService.checkConnection")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)
	span.SetTag("mailbox.id", mailboxID)

	// Get the client
	s.clientsMutex.RLock()
	c, exists := s.clients[mailboxID]
	s.clientsMutex.RUnlock()

	if !exists {
		log.Printf("[%s] Client not found during health check", mailboxID)
		span.SetTag("status", "not_found")
		return false
	}

	// Set a timeout for the NOOP command
	c.Timeout = 10 * time.Second

	// Perform a NOOP to check connection
	err := c.Noop()
	c.Timeout = 0 // Reset timeout

	if err != nil {
		log.Printf("[%s] Health check failed: %v", mailboxID, err)
		span.SetTag("status", "unhealthy")
		span.SetTag("error", err.Error())
		tracing.TraceErr(span, err)

		// Update status
		s.updateStatusError(mailboxID, err)

		// Remove from clients map
		s.clientsMutex.Lock()
		delete(s.clients, mailboxID)
		s.clientsMutex.Unlock()

		// Trigger reconnection
		go s.reconnectWithBackoff(mailboxID)

		return false
	}

	// Update status to show connection is healthy
	s.updateConnectionHealth(mailboxID, true, "")
	span.SetTag("status", "healthy")

	return true
}

// reconnectWithBackoff attempts to reconnect with exponential backoff
func (s *IMAPService) reconnectWithBackoff(mailboxID string) {
	span := opentracing.StartSpan("IMAPService.reconnectWithBackoff")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)
	span.SetTag("mailbox.id", mailboxID)

	// Get the config
	s.clientsMutex.RLock()
	config, exists := s.configs[mailboxID]
	s.clientsMutex.RUnlock()

	if !exists {
		log.Printf("[%s] Config not found, cannot reconnect", mailboxID)
		span.SetTag("error", "config_not_found")
		return
	}

	// Set up backoff parameters
	backoff := time.Second
	maxBackoff := time.Minute * 5
	maxAttempts := 10
	attempts := 0

	for attempts < maxAttempts {
		// Check if service is shutting down
		if s.ctx.Err() != nil {
			log.Printf("[%s] Reconnection cancelled due to service shutdown", mailboxID)
			span.SetTag("cancelled", true)
			return
		}

		attempts++
		span.SetTag("attempt", attempts)
		span.SetTag("backoff", backoff.String())

		log.Printf("[%s] Reconnection attempt %d/%d after %v backoff",
			mailboxID, attempts, maxAttempts, backoff)

		// Update status to show reconnection attempt
		s.updateConnectionHealth(mailboxID, false, fmt.Sprintf("Reconnecting (attempt %d/%d)", attempts, maxAttempts))

		// Wait for backoff period
		select {
		case <-time.After(backoff):
			// Continue with reconnection
		case <-s.ctx.Done():
			log.Printf("[%s] Reconnection cancelled due to service shutdown", mailboxID)
			span.SetTag("cancelled", true)
			return
		}

		// Try to reconnect
		reconnectSpan := opentracing.StartSpan(
			"IMAPService.reconnectAttempt",
			opentracing.ChildOf(span.Context()),
		)
		reconnectSpan.SetTag("mailbox.id", mailboxID)
		reconnectSpan.SetTag("attempt", attempts)

		c, err := s.connectMailbox(context.Background(), config)
		if err != nil {
			log.Printf("[%s] Reconnection attempt %d failed: %v", mailboxID, attempts, err)
			reconnectSpan.SetTag("error", err.Error())
			tracing.TraceErr(reconnectSpan, err)
			reconnectSpan.Finish()

			// Increase backoff with jitter
			backoff = addJitter(backoff * 2)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
			continue
		}

		// Reconnection successful
		log.Printf("[%s] Successfully reconnected on attempt %d", mailboxID, attempts)
		reconnectSpan.SetTag("success", true)
		reconnectSpan.Finish()

		// Store the new client
		s.clientsMutex.Lock()
		s.clients[mailboxID] = c
		s.clientsMutex.Unlock()

		// Update status
		s.updateConnectionHealth(mailboxID, true, "")

		// Start monitoring folders
		go func() {
			for _, folder := range config.Folders {
				// Start monitoring each folder in a separate goroutine
				folderName := string(folder)
				go func(fn string) {
					err := s.monitorFolder(s.ctx, mailboxID, c, fn)
					if err != nil {
						log.Printf("[%s][%s] Error monitoring folder after reconnection: %v",
							mailboxID, fn, err)
					}
				}(folderName)
			}
		}()

		return
	}

	// Max attempts reached
	log.Printf("[%s] Failed to reconnect after %d attempts", mailboxID, maxAttempts)
	span.SetTag("max_attempts_reached", true)

	// Update status to show permanent failure
	s.updateConnectionHealth(mailboxID, false, fmt.Sprintf("Failed to reconnect after %d attempts", maxAttempts))
}

// updateConnectionHealth updates the connection health status
func (s *IMAPService) updateConnectionHealth(mailboxID string, isHealthy bool, message string) {
	s.statusMutex.Lock()
	defer s.statusMutex.Unlock()

	status, exists := s.statuses[mailboxID]
	if !exists {
		status = interfaces.MailboxStatus{
			Folders: make(map[string]interfaces.FolderStats),
		}
	}

	status.Connected = isHealthy
	status.LastChecked = time.Now()

	if !isHealthy && message != "" {
		status.LastError = message
	} else if isHealthy {
		status.LastError = ""
	}

	s.statuses[mailboxID] = status
}

// getMailboxStatuses gets the status of all mailboxes
func (s *IMAPService) getMailboxStatuses() map[string]interfaces.MailboxStatus {
	s.statusMutex.RLock()
	defer s.statusMutex.RUnlock()

	// Create a copy to avoid race conditions
	result := make(map[string]interfaces.MailboxStatus, len(s.statuses))
	for id, status := range s.statuses {
		// Deep copy the status
		foldersCopy := make(map[string]interfaces.FolderStats, len(status.Folders))
		for folderName, stats := range status.Folders {
			foldersCopy[folderName] = stats
		}

		statusCopy := interfaces.MailboxStatus{
			Connected:   status.Connected,
			LastError:   status.LastError,
			LastChecked: status.LastChecked,
			Folders:     foldersCopy,
		}

		result[id] = statusCopy
	}

	return result
}

// Helper function to add jitter to a duration
func addJitter(d time.Duration) time.Duration {
	// Add up to 20% jitter
	jitterFactor := 0.8 + 0.4*rand.Float64() // between 0.8 and 1.2
	return time.Duration(float64(d) * jitterFactor)
}
