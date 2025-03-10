package imap

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

// monitorFolder handles the monitoring of a single folder
func (s *IMAPService) monitorFolder(ctx context.Context, mailboxID string, c *client.Client, folderName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.monitorFolder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	// Select the mailbox (folder)
	mbox, err := s.selectFolder(ctx, c, mailboxID, folderName)
	if err != nil {
		return err
	}

	// Process initial messages with pagination
	err = s.processInitialMessages(ctx, c, mailboxID, folderName, mbox)
	if err != nil {
		log.Printf("[%s][%s] Warning: Error processing initial messages: %v", mailboxID, folderName, err)
		// Continue anyway - this is not fatal
	}

	// Set up IDLE monitoring
	err = s.setupIdleMonitoring(ctx, c, mailboxID, folderName, mbox.Messages)

	return err
}

// selectFolder selects a folder on the IMAP server
func (s *IMAPService) selectFolder(ctx context.Context, c *client.Client, mailboxID, folderName string) (*imap.MailboxStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.selectFolder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	// Set client timeout temporarily
	c.Timeout = 30 * time.Second
	mbox, err := c.Select(folderName, false)
	c.Timeout = 0 // Reset timeout

	if err != nil {
		err = fmt.Errorf("[%s][%s] Error selecting folder: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	log.Printf("[%s][%s] Selected folder - Messages: %d, Recent: %d, Unseen: %d",
		mailboxID, folderName, mbox.Messages, mbox.Recent, mbox.Unseen)

	span.SetTag("messages.total", mbox.Messages)
	span.SetTag("messages.recent", mbox.Recent)
	span.SetTag("messages.unseen", mbox.Unseen)

	// Update folder stats
	s.updateFolderStats(ctx, mailboxID, folderName, mbox)

	return mbox, nil
}

// updateFolderStats updates the statistics for a folder
func (s *IMAPService) updateFolderStats(ctx context.Context, mailboxID, folderName string, mbox *imap.MailboxStatus) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.updateFolderStats")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(context.Background(), span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

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
	unseenCount := uint32(0)

	// If we have a valid client, try to get actual unseen count
	s.clientsMutex.RLock()
	client, hasClient := s.clients[mailboxID]
	s.clientsMutex.RUnlock()

	if hasClient {
		// Create search criteria for unseen messages
		criteria := imap.NewSearchCriteria()
		criteria.WithoutFlags = []string{imap.SeenFlag}

		// Set client timeout temporarily
		client.Timeout = 10 * time.Second
		ids, err := client.Search(criteria)
		client.Timeout = 0 // Reset timeout

		if err == nil {
			unseenCount = uint32(len(ids))
			span.SetTag("unseen.count", unseenCount)
		} else {
			log.Printf("[%s][%s] Error counting unseen messages: %v", mailboxID, folderName, err)
			tracing.TraceErr(span, err)
			// Fall back to mbox.Unseen which may be less accurate
			unseenCount = mbox.Unseen
		}
	} else {
		// Fall back to mbox.Unseen if we don't have a client
		unseenCount = mbox.Unseen
	}

	// Update folder stats
	status.Folders[folderName] = interfaces.FolderStats{
		Total:    mbox.Messages,
		Unseen:   unseenCount,
		LastSeen: s.getLastSyncedUID(ctx, mailboxID, folderName),
		LastSync: time.Now(),
	}

	s.statuses[mailboxID] = status
}

// listFolders lists all available folders on the server
func (s *IMAPService) listFolders(ctx context.Context, c *client.Client, mailboxID string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.listFolders")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)

	// Set client timeout temporarily
	c.Timeout = 30 * time.Second
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	var folders []string
	for m := range mailboxes {
		folders = append(folders, m.Name)
	}

	c.Timeout = 0 // Reset timeout
	err := <-done
	if err != nil {
		err = fmt.Errorf("[%s] Error listing folders: %v", mailboxID, err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Sort folders alphabetically for consistency
	sort.Strings(folders)

	span.SetTag("folders.count", len(folders))
	log.Printf("[%s] Found %d folders: %v", mailboxID, len(folders), folders)

	return folders, nil
}

// processInitialMessages handles the initial sync of messages with pagination
func (s *IMAPService) processInitialMessages(ctx context.Context, c *client.Client, mailboxID, folderName string, mbox *imap.MailboxStatus) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.processInitialMessages")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	// Get the last synchronized UID
	lastSyncedUID := s.getLastSyncedUID(ctx, mailboxID, folderName)
	span.SetTag("last_synced_uid", lastSyncedUID)

	// If we have a previous sync point, use it
	if lastSyncedUID > 0 {
		log.Printf("[%s][%s] Resuming sync from UID %d", mailboxID, folderName, lastSyncedUID)
		return s.syncNewMessagesSinceUID(ctx, c, mailboxID, folderName, lastSyncedUID)
	}

	// For first-time sync, paginate through unseen messages
	return s.paginatedInitialSync(ctx, c, mailboxID, folderName)
}

// paginatedInitialSync performs a paginated sync of unseen messages
func (s *IMAPService) paginatedInitialSync(ctx context.Context, c *client.Client, mailboxID, folderName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.paginatedInitialSync")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)

	// Search for unseen messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	// Set client timeout temporarily
	c.Timeout = 30 * time.Second
	uids, err := c.UidSearch(criteria)
	c.Timeout = 0 // Reset timeout

	if err != nil {
		err = fmt.Errorf("[%s][%s] Error searching for unseen messages: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		return err
	}

	totalMessages := len(uids)
	span.SetTag("unseen.total", totalMessages)

	if totalMessages == 0 {
		log.Printf("[%s][%s] No unseen messages to sync", mailboxID, folderName)
		return nil
	}

	// Limit the total number of messages to process
	maxToProcess := INITIAL_SYNC_MAX_TOTAL
	if totalMessages > maxToProcess {
		log.Printf("[%s][%s] Limiting initial sync to %d of %d unseen messages",
			mailboxID, folderName, maxToProcess, totalMessages)

		// Sort UIDs in descending order (newest first)
		sort.Slice(uids, func(i, j int) bool {
			return uids[i] > uids[j]
		})

		uids = uids[:maxToProcess]
	}

	// Process in batches
	batchSize := INITIAL_SYNC_BATCH_SIZE

	log.Printf("[%s][%s] Starting paginated sync of %d messages in batches of %d",
		mailboxID, folderName, len(uids), batchSize)

	for i := 0; i < len(uids); i += batchSize {
		// Check if context is cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}

		end := i + batchSize
		if end > len(uids) {
			end = len(uids)
		}

		batchUIDs := uids[i:end]

		log.Printf("[%s][%s] Processing batch %d-%d of %d",
			mailboxID, folderName, i+1, end, len(uids))

		// Create a sequence set for this batch
		seqSet := new(imap.SeqSet)
		for _, uid := range batchUIDs {
			seqSet.AddNum(uid)
		}

		// Fetch and process this batch with timeout
		batchCtx, batchCancel := context.WithTimeout(ctx, 60*time.Second)
		err := s.fetchAndProcessMessages(batchCtx, c, mailboxID, folderName, seqSet, true)
		batchCancel()

		if err != nil {
			log.Printf("[%s][%s] Error processing batch: %v", mailboxID, folderName, err)
			// Continue with next batch instead of failing completely
		}

		// Short pause between batches to prevent overloading
		select {
		case <-ctx.Done():
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
		log.Printf("[%s][%s] Updated last synced UID to %d", mailboxID, folderName, highestUID)
	}

	log.Printf("[%s][%s] Completed initial sync of %d messages", mailboxID, folderName, len(uids))

	return nil
}

// syncNewMessagesSinceUID syncs messages that arrived since the last sync
func (s *IMAPService) syncNewMessagesSinceUID(ctx context.Context, c *client.Client, mailboxID, folderName string, lastUID uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.syncNewMessagesSinceUID")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)
	span.SetTag("last_uid", lastUID)

	// Create search criteria for messages with UID greater than lastUID
	criteria := imap.NewSearchCriteria()
	uidRange := new(imap.SeqSet)
	uidRange.AddRange(lastUID+1, 0) // From lastUID+1 to infinity
	criteria.Uid = uidRange

	// Set client timeout temporarily
	c.Timeout = 30 * time.Second
	uids, err := c.UidSearch(criteria)
	c.Timeout = 0 // Reset timeout

	if err != nil {
		err = fmt.Errorf("[%s][%s] Error searching for new messages: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		return err
	}

	if len(uids) == 0 {
		log.Printf("[%s][%s] No new messages since UID %d", mailboxID, folderName, lastUID)
		return nil
	}

	log.Printf("[%s][%s] Found %d new messages since UID %d", mailboxID, folderName, len(uids), lastUID)
	span.SetTag("new_messages", len(uids))

	// Create a sequence set for fetching
	seqSet := new(imap.SeqSet)
	for _, uid := range uids {
		seqSet.AddNum(uid)
	}

	// Fetch the new messages
	err = s.fetchAndProcessMessages(ctx, c, mailboxID, folderName, seqSet, true)
	if err != nil {
		return err
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
		log.Printf("[%s][%s] Updated last synced UID to %d", mailboxID, folderName, highestUID)
	}

	return nil
}

// fetchAndProcessMessages fetches and processes messages with the given sequence set
func (s *IMAPService) fetchAndProcessMessages(
	ctx context.Context,
	c *client.Client,
	mailboxID,
	folderName string,
	seqSet *imap.SeqSet,
	isUID bool,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.fetchAndProcessMessages")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)
	span.SetTag("is_uid", isUID)

	// Items to fetch - balance between getting enough data and not overloading
	items := []imap.FetchItem{
		imap.FetchEnvelope,
		imap.FetchFlags,
		imap.FetchBodyStructure,
		"BODY.PEEK[]",
		imap.FetchUid,
	}

	// Create channels for fetching
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	// Set client timeout for fetch
	c.Timeout = 60 * time.Second

	// Start fetching
	go func() {
		if isUID {
			done <- c.UidFetch(seqSet, items, messages)
		} else {
			done <- c.Fetch(seqSet, items, messages)
		}
	}()

	// Process messages as they arrive
	var highestUID uint32
	messageCount := 0

	// Create a context with timeout for processing
	processCtx, processCancel := context.WithCancel(ctx)
	defer processCancel()

	// Set up processing with timeout protection
	go func() {
		select {
		case <-ctx.Done():
			return // Parent context was cancelled
		case <-time.After(60 * time.Second):
			log.Printf("[%s][%s] Warning: Message processing timeout", mailboxID, folderName)
			processCancel() // Cancel processing if it takes too long
		}
	}()

	for msg := range messages {
		select {
		case <-processCtx.Done():
			// If processing is cancelled, drain the channel and exit
			for range messages {
				// Drain remaining messages
			}
			return fmt.Errorf("message processing cancelled")
		default:
			// Process the message
			messageCount++

			if msg.Uid > highestUID {
				highestUID = msg.Uid
			}

			// Add processing span
			msgSpan := opentracing.StartSpan(
				"IMAPService.processMessage",
				opentracing.ChildOf(span.Context()),
			)
			msgSpan.SetTag("mailbox.id", mailboxID)
			msgSpan.SetTag("folder.name", folderName)
			msgSpan.SetTag("uid", msg.Uid)
			msgSpan.SetTag("seq_num", msg.SeqNum)

			// Process the message if we have a handler
			if s.eventHandler != nil {
				s.eventHandler(ctx, interfaces.MailEvent{
					Source:    "imap",
					MailboxID: mailboxID,
					Folder:    folderName,
					MessageID: msg.SeqNum,
					EventType: "new",
					Message:   msg,
				})
			}

			msgSpan.Finish()
		}
	}

	// Reset client timeout
	c.Timeout = 0

	// Wait for fetch to complete
	err := <-done
	if err != nil {
		err = fmt.Errorf("[%s][%s] Error fetching messages: %v", mailboxID, folderName, err)
		tracing.TraceErr(span, err)
		return err
	}

	log.Printf("[%s][%s] Processed %d messages", mailboxID, folderName, messageCount)
	span.SetTag("messages.processed", messageCount)

	// Update last synced UID if this was a UID fetch
	if isUID && highestUID > 0 {
		s.saveLastSyncedUID(ctx, mailboxID, folderName, highestUID)
	}

	return nil
}

// fetchNewMessages fetches new messages by sequence number
func (s *IMAPService) fetchNewMessages(ctx context.Context, mailboxID string, c *client.Client, folderName string, from, to uint32) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPService.fetchNewMessages")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag("mailbox.id", mailboxID)
	span.SetTag("folder.name", folderName)
	span.SetTag("from", from)
	span.SetTag("to", to)

	if from > to {
		return nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	return s.fetchAndProcessMessages(ctx, c, mailboxID, folderName, seqSet, false)
}
