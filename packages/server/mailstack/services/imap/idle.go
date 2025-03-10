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
)

// setupIdleMonitoring sets up IDLE monitoring for a folder
func (s *IMAPService) setupIdleMonitoring(ctx context.Context, c *client.Client, mailboxID, folderName string, initialCount uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.setupIdleMonitoring")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)
	span.SetTag("initial_count", initialCount)

	// Use sync.Once to prevent multiple channel closes
	var stopOnce sync.Once
	stop := make(chan struct{}, 1)

	// Safe channel closer
	safeClose := func() {
		stopOnce.Do(func() {
			close(stop)
		})
	}
	defer safeClose()

	// Recover from potential panics
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[%s][%s] Recovered from panic: %v", mailboxID, folderName, r)
		}
	}()

	// Check IDLE support (existing code)
	supportSpan := opentracing.StartSpan(
		"IMAPService.checkIdleSupport",
		opentracing.ChildOf(span.Context()),
	)
	supported, err := c.Support("IDLE")
	if err != nil {
		log.Printf("[%s][%s] Error checking IDLE support: %v", mailboxID, folderName, err)
		tracing.TraceErr(supportSpan, err)
		supportSpan.SetTag("error", true)
	}
	supportSpan.SetTag("idle_supported", supported)
	supportSpan.Finish()

	// Updates channel with buffer
	updates := make(chan client.Update, 100)
	c.Updates = updates

	// Set up monitoring goroutines
	errChan := s.setupIdleGoroutines(ctx, c, mailboxID, folderName, updates, initialCount, stop)

	// IDLE command
	idleSpan := opentracing.StartSpan(
		"IMAPService.idleCommand",
		opentracing.ChildOf(span.Context()),
	)
	idleSpan.SetTag("mailbox.id", mailboxID)
	idleSpan.SetTag("folder.name", folderName)

	c.Timeout = 0
	idleErr := c.Idle(stop, &client.IdleOptions{
		LogoutTimeout: DEFAULT_IMAP_LOGOUT,
		PollInterval:  DEFAULT_POLLING_PERIOD,
	})

	if idleErr != nil && ctx.Err() == nil {
		log.Printf("[%s][%s] IDLE error: %v", mailboxID, folderName, idleErr)
		tracing.TraceErr(idleSpan, idleErr)
		idleSpan.SetTag("error", true)
	}
	idleSpan.Finish()

	// Signal stop
	safeClose()

	// Wait for monitoring goroutines
	select {
	case err := <-errChan:
		if err != nil {
			log.Printf("[%s][%s] Monitoring error: %v", mailboxID, folderName, err)
			return err
		}
	case <-time.After(5 * time.Second):
		log.Printf("[%s][%s] Timed out waiting for monitoring to finish", mailboxID, folderName)
	}

	c.Updates = nil
	log.Printf("[%s][%s] Stopped monitoring folder", mailboxID, folderName)
	return nil
}

// setupIdleGoroutines sets up goroutines for IDLE monitoring
func (s *IMAPService) setupIdleGoroutines(
	ctx context.Context,
	c *client.Client,
	mailboxID, folderName string,
	updates chan client.Update,
	initialCount uint32,
	stop chan struct{},
) chan error {
	span := opentracing.StartSpan("IMAPService.setupIdleGoroutines")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	errChan := make(chan error, 1)

	// Context cancel monitoring
	go func() {
		select {
		case <-ctx.Done():
			log.Printf("[%s][%s] Context cancelled, stopping IDLE", mailboxID, folderName)
			select {
			case <-stop: // Stop channel already closed
			default:
				close(stop)
			}
		case <-stop: // Stop initiated elsewhere
		}
	}()

	// NOOP keepalive (as a safety net in addition to IDLE's built-in polling)
	go func() {
		// Use a longer interval since IDLE already does polling
		noopTicker := time.NewTicker(5 * time.Minute)
		defer noopTicker.Stop()

		for {
			select {
			case <-noopTicker.C:
				noopSpan := opentracing.StartSpan(
					"IMAPService.idleNoop",
					opentracing.ChildOf(span.Context()),
				)
				noopSpan.SetTag("mailbox.id", mailboxID)
				noopSpan.SetTag("folder.name", folderName)

				// Set timeout for NOOP
				c.Timeout = 10 * time.Second
				err := c.Noop()
				c.Timeout = 0 // Reset timeout

				if err != nil {
					log.Printf("[%s][%s] Error during NOOP: %v", mailboxID, folderName, err)
					tracing.TraceErr(noopSpan, err)
					noopSpan.SetTag("error", true)
					noopSpan.Finish()

					select {
					case errChan <- err:
					default:
					}

					// Try to close stop channel if not already closed
					select {
					case <-stop:
					default:
						close(stop)
					}
					return
				}
				noopSpan.Finish()
			case <-stop:
				return
			}
		}
	}()

	// Update processor
	go func() {
		currentCount := initialCount
		messageSpan := opentracing.StartSpan(
			"IMAPService.processIdleUpdates",
			opentracing.ChildOf(span.Context()),
		)
		messageSpan.SetTag("mailbox.id", mailboxID)
		messageSpan.SetTag("folder.name", folderName)
		defer messageSpan.Finish()

		// Track when we last saw activity
		lastActivity := time.Now()
		staleThreshold := 30 * time.Minute

		// Set up a ticker to check for stale connections
		staleTicker := time.NewTicker(5 * time.Minute)
		defer staleTicker.Stop()

		for {
			select {
			case update, ok := <-updates:
				if !ok {
					log.Printf("[%s][%s] Updates channel closed", mailboxID, folderName)
					messageSpan.SetTag("reason", "channel_closed")
					close(errChan)
					return
				}

				// Update last activity time
				lastActivity = time.Now()

				updateSpan := opentracing.StartSpan(
					"IMAPService.processIdleUpdate",
					opentracing.ChildOf(messageSpan.Context()),
				)
				updateSpan.SetTag("mailbox.id", mailboxID)
				updateSpan.SetTag("folder.name", folderName)
				updateSpan.SetTag("update_type", fmt.Sprintf("%T", update))

				log.Printf("[%s][%s] Received update: %T", mailboxID, folderName, update)

				switch u := update.(type) {
				case *client.MailboxUpdate:
					updateSpan.SetTag("messages.previous", currentCount)
					updateSpan.SetTag("messages.current", u.Mailbox.Messages)

					log.Printf("[%s][%s] Mailbox update - Messages: %d (was: %d)",
						mailboxID, folderName, u.Mailbox.Messages, currentCount)

					// If we have new messages
					if u.Mailbox.Messages > currentCount {
						newMessages := u.Mailbox.Messages - currentCount
						log.Printf("[%s][%s] 📥 Detected %d new message(s)", mailboxID, folderName, newMessages)
						updateSpan.SetTag("new_messages", newMessages)

						// Fetch new messages with error handling
						fetchSpan := opentracing.StartSpan(
							"IMAPService.fetchNewIdleMessages",
							opentracing.ChildOf(updateSpan.Context()),
						)
						fetchSpan.SetTag("mailbox.id", mailboxID)
						fetchSpan.SetTag("folder.name", folderName)
						fetchSpan.SetTag("from", currentCount+1)
						fetchSpan.SetTag("to", u.Mailbox.Messages)

						fetchErr := s.fetchNewMessages(context.Background(), mailboxID, c, folderName, currentCount+1, u.Mailbox.Messages)
						if fetchErr != nil {
							log.Printf("[%s][%s] Error fetching new messages: %v",
								mailboxID, folderName, fetchErr)
							tracing.TraceErr(fetchSpan, fetchErr)
							fetchSpan.SetTag("error", true)
						}
						fetchSpan.Finish()

						currentCount = u.Mailbox.Messages
					}

					// Update folder stats
					s.updateFolderStats(ctx, mailboxID, folderName, u.Mailbox)

				case *client.ExpungeUpdate:
					log.Printf("[%s][%s] Message expunged: %d", mailboxID, folderName, u.SeqNum)
					updateSpan.SetTag("expunged_seq", u.SeqNum)

					if u.SeqNum <= currentCount {
						currentCount--
						updateSpan.SetTag("messages.current", currentCount)
					}

				case *client.MessageUpdate:
					if u.Message != nil {
						log.Printf("[%s][%s] Message updated: %d", mailboxID, folderName, u.Message.SeqNum)
						updateSpan.SetTag("updated_seq", u.Message.SeqNum)
						updateSpan.SetTag("updated_uid", u.Message.Uid)

						// If flagged as recent or unseen, process it
						isRecent := false
						for _, flag := range u.Message.Flags {
							if flag == "\\Recent" {
								isRecent = true
								break
							}
						}

						if isRecent {
							log.Printf("[%s][%s] Processing updated message marked as recent: %d",
								mailboxID, folderName, u.Message.SeqNum)

							// We would normally fetch more details here,
							// but since we're just getting a flag update, log it
							updateSpan.SetTag("is_recent", true)
						}
					}

				default:
					log.Printf("[%s][%s] Received update of unknown type: %T", mailboxID, folderName, update)
				}

				updateSpan.Finish()

			case <-staleTicker.C:
				// Check if the connection appears stale
				if time.Since(lastActivity) > staleThreshold {
					log.Printf("[%s][%s] Connection appears stale (no activity for %v), reconnecting",
						mailboxID, folderName, time.Since(lastActivity))

					messageSpan.SetTag("reason", "stale_connection")
					messageSpan.SetTag("stale_duration", time.Since(lastActivity).String())

					// Signal reconnection by returning an error
					errChan <- fmt.Errorf("stale connection detected")

					// Try to close stop channel if not already closed
					select {
					case <-stop:
					default:
						close(stop)
					}
					return
				}

			case <-stop:
				messageSpan.SetTag("reason", "stopped")
				close(errChan)
				return
			}
		}
	}()

	return errChan
}

// requestIDLE sends an IDLE command to the server
func (s *IMAPService) requestIDLE(c *client.Client, stop chan struct{}) error {
	span := opentracing.StartSpan("IMAPService.requestIDLE")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)

	// Clear any pending timeout
	c.Timeout = 0

	// Start IDLE command
	err := c.Idle(stop, &client.IdleOptions{
		LogoutTimeout: DEFAULT_IMAP_LOGOUT,
		PollInterval:  DEFAULT_POLLING_PERIOD,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		span.SetTag("error", true)
		return err
	}

	return nil
}

// pollFolder is a fallback for servers that don't support IDLE
func (s *IMAPService) pollFolder(ctx context.Context, c *client.Client, mailboxID, folderName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.pollFolder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	log.Printf("[%s][%s] Starting polling mode (fallback for non-IDLE servers)", mailboxID, folderName)

	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	lastKnownCount := uint32(0)
	var lastError time.Time

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Don't poll too frequently after errors
			if !lastError.IsZero() && time.Since(lastError) < time.Minute {
				continue
			}

			// Select the mailbox to get current status
			mbox, err := c.Select(folderName, true) // Read-only mode
			if err != nil {
				log.Printf("[%s][%s] Error selecting folder during poll: %v", mailboxID, folderName, err)
				lastError = time.Now()
				continue
			}

			// Check for new messages
			if mbox.Messages > lastKnownCount {
				if lastKnownCount > 0 { // Skip initial poll
					newCount := mbox.Messages - lastKnownCount
					log.Printf("[%s][%s] Poll detected %d new message(s)", mailboxID, folderName, newCount)

					// Fetch the new messages
					err := s.fetchNewMessages(ctx, mailboxID, c, folderName, lastKnownCount+1, mbox.Messages)
					if err != nil {
						log.Printf("[%s][%s] Error fetching new messages during poll: %v", mailboxID, folderName, err)
						lastError = time.Now()
					}
				}
				lastKnownCount = mbox.Messages
			}

			// Update folder stats
			s.updateFolderStats(ctx, mailboxID, folderName, mbox)
		}
	}
}
