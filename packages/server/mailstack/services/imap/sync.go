package imap

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/emersion/go-imap"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
)

// getLastSyncedUID gets the last synced UID for a mailbox folder
func (s *IMAPService) getLastSyncedUID(ctx context.Context, mailboxID, folderName string) uint32 {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.getLastSyncedUID")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	state, err := s.repositories.MailboxSyncRepository.GetSyncState(ctx, mailboxID, folderName)
	if err != nil {
		log.Printf("[%s][%s] Error getting sync state: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		return 0
	}

	if state == nil {
		span.SetTag("found", false)
		return 0 // No sync state yet
	}

	span.SetTag("found", true)
	span.SetTag("last_uid", state.LastUID)
	return state.LastUID
}

// saveLastSyncedUID saves the last synced UID for a mailbox folder
func (s *IMAPService) saveLastSyncedUID(ctx context.Context, mailboxID, folderName string, uid uint32) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.saveLastSyncedUID")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)
	span.SetTag("uid", uid)

	state := &models.MailboxSyncState{
		ID:         fmt.Sprintf("%s-%s", mailboxID, folderName),
		MailboxID:  mailboxID,
		FolderName: folderName,
		LastUID:    uid,
		LastSync:   time.Now(),
	}

	err := s.repositories.MailboxSyncRepository.SaveSyncState(ctx, state)
	if err != nil {
		log.Printf("[%s][%s] Error saving sync state: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		span.SetTag("error", true)
	}

	log.Printf("[%s][%s] Updated last synced UID to %d", mailboxID, folderName, uid)
}

// resetSyncState resets the sync state for a specific mailbox and folder
func (s *IMAPService) resetSyncState(ctx context.Context, mailboxID, folderName string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.resetSyncState")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	err := s.repositories.MailboxSyncRepository.DeleteSyncState(ctx, mailboxID, folderName)
	if err != nil {
		log.Printf("[%s][%s] Error resetting sync state: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		span.SetTag("error", true)
	}

	log.Printf("[%s][%s] Reset sync state", mailboxID, folderName)
}

// resetAllSyncState resets the sync state for all mailboxes
func (s *IMAPService) resetAllSyncState(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.resetAllSyncState")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get configured mailboxes
	s.clientsMutex.RLock()
	mailboxIDs := make([]string, 0, len(s.configs))
	for id := range s.configs {
		mailboxIDs = append(mailboxIDs, id)
	}
	s.clientsMutex.RUnlock()

	// Reset sync state for each mailbox
	for _, id := range mailboxIDs {
		err := s.repositories.MailboxSyncRepository.DeleteMailboxSyncStates(ctx, id)
		if err != nil {
			log.Printf("[%s] Error resetting mailbox sync state: %v", id, err)
			tracing.TraceErr(span, err)
		}
	}

	log.Printf("Reset all sync states")
}

// isInInitialSync checks if a mailbox folder is currently in initial sync
func (s *IMAPService) isInInitialSync(ctx context.Context, mailboxID, folderName string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.isInInitialSync")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	state, err := s.repositories.MailboxSyncRepository.GetSyncState(ctx, mailboxID, folderName)
	if err != nil {
		log.Printf("[%s][%s] Error checking sync state: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		span.SetTag("error", true)
		return true // Assume initial sync on error
	}

	result := state == nil || state.LastUID == 0
	span.SetTag("is_initial_sync", result)
	return result
}

// cleanSyncState removes entries for mailboxes that no longer exist
func (s *IMAPService) cleanSyncState(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.cleanSyncState")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get the set of currently configured mailboxes
	s.clientsMutex.RLock()
	activeMailboxes := make(map[string]bool)
	for id := range s.configs {
		activeMailboxes[id] = true
	}
	s.clientsMutex.RUnlock()

	// Get all sync states
	states, err := s.repositories.MailboxSyncRepository.GetAllSyncStates(ctx)
	if err != nil {
		log.Printf("Error getting all sync states: %v", err)
		tracing.TraceErr(span, err)
		span.SetTag("error", true)
		return
	}

	// Find and remove orphaned sync states
	for mailboxID := range states {
		if !activeMailboxes[mailboxID] {
			log.Printf("Cleaning up sync state for removed mailbox: %s", mailboxID)
			err := s.repositories.MailboxSyncRepository.DeleteMailboxSyncStates(ctx, mailboxID)
			if err != nil {
				log.Printf("[%s] Error deleting orphaned sync state: %v", mailboxID, err)
				tracing.TraceErr(span, err)
			}
		}
	}
}

// runPeriodicSyncMaintenance runs periodic maintenance on the sync state
func (s *IMAPService) runPeriodicSyncMaintenance() {
	log.Println("Starting periodic sync maintenance")
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			s.cleanSyncState(ctx)
			cancel()
		case <-s.ctx.Done():
			log.Println("Stopping periodic sync maintenance")
			return
		}
	}
}

// fullResync performs a full resynchronization for a mailbox
func (s *IMAPService) fullResync(ctx context.Context, mailboxID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.fullResync")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)

	// Reset sync state for this mailbox
	err := s.repositories.MailboxSyncRepository.DeleteMailboxSyncStates(ctx, mailboxID)
	if err != nil {
		log.Printf("[%s] Error resetting sync state: %v", mailboxID, err)
		tracing.TraceErr(span, err)
		// Continue anyway
	}

	// Get the mailbox configuration
	s.clientsMutex.RLock()
	config, exists := s.configs[mailboxID]
	s.clientsMutex.RUnlock()

	if !exists {
		err := fmt.Errorf("mailbox configuration not found: %s", mailboxID)
		tracing.TraceErr(span, err)
		return err
	}

	// Get a client or create a new one
	client, err := s.getClient(ctx, mailboxID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Process each folder
	for _, folder := range config.Folders {
		folderName := string(folder)
		folderSpan := opentracing.StartSpan(
			"IMAPService.resyncFolder",
			opentracing.ChildOf(span.Context()),
		)
		folderSpan.SetTag("mailbox.id", mailboxID)
		folderSpan.SetTag("folder.name", folderName)

		// Select the folder
		_, err := client.Select(folderName, false)
		if err != nil {
			err = fmt.Errorf("error selecting folder %s: %w", folderName, err)
			tracing.TraceErr(folderSpan, err)
			folderSpan.Finish()
			continue
		}

		// Create a search criteria for all messages
		criteria := imap.NewSearchCriteria()

		// Search for all messages
		uids, err := client.UidSearch(criteria)
		if err != nil {
			err = fmt.Errorf("error searching messages in folder %s: %w", folderName, err)
			tracing.TraceErr(folderSpan, err)
			folderSpan.Finish()
			continue
		}

		// Check if we have any messages
		if len(uids) == 0 {
			log.Printf("[%s][%s] No messages to resync", mailboxID, folderName)
			folderSpan.Finish()
			continue
		}

		log.Printf("[%s][%s] Starting full resync of %d messages", mailboxID, folderName, len(uids))

		// Process in batches with pagination
		batchSize := INITIAL_SYNC_BATCH_SIZE
		totalMessages := len(uids)
		maxMessages := INITIAL_SYNC_MAX_TOTAL

		// Limit the total number of messages
		if totalMessages > maxMessages {
			log.Printf("[%s][%s] Limiting resync to %d of %d messages",
				mailboxID, folderName, maxMessages, totalMessages)

			// Sort UIDs in descending order to get newest first
			sort.Slice(uids, func(i, j int) bool {
				return uids[i] > uids[j]
			})

			uids = uids[:maxMessages]
		}

		// Process in batches
		for i := 0; i < len(uids); i += batchSize {
			// Check if context is cancelled
			if ctx.Err() != nil {
				folderSpan.Finish()
				return ctx.Err()
			}

			end := i + batchSize
			if end > len(uids) {
				end = len(uids)
			}

			batchUIDs := uids[i:end]
			log.Printf("[%s][%s] Processing batch %d-%d of %d",
				mailboxID, folderName, i+1, end, len(uids))

			// Create sequence set for this batch
			seqSet := new(imap.SeqSet)
			for _, uid := range batchUIDs {
				seqSet.AddNum(uid)
			}

			// Fetch and process batch
			batchCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
			err := s.fetchAndProcessMessages(batchCtx, client, mailboxID, folderName, seqSet, true)
			cancel()

			if err != nil {
				log.Printf("[%s][%s] Error processing batch: %v", mailboxID, folderName, err)
				// Continue with next batch
			}

			// Add a small delay between batches
			select {
			case <-ctx.Done():
				folderSpan.Finish()
				return ctx.Err()
			case <-time.After(100 * time.Millisecond):
				// Continue with next batch
			}
		}

		// Find highest UID and save it
		var highestUID uint32
		for _, uid := range uids {
			if uint32(uid) > highestUID {
				highestUID = uint32(uid)
			}
		}

		if highestUID > 0 {
			s.saveLastSyncedUID(ctx, mailboxID, folderName, highestUID)
		}

		folderSpan.Finish()
	}

	return nil
}

// waitForInitialSync waits for initial sync to complete with timeout
func (s *IMAPService) waitForInitialSync(ctx context.Context, mailboxID, folderName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.waitForInitialSync")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	// Create a timeout context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// Check every second if the initial sync is complete
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if !s.isInInitialSync(ctx, mailboxID, folderName) {
				return nil // Initial sync is complete
			}
		}
	}
}
