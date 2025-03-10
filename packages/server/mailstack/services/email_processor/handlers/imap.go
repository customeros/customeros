package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/mailsherpa/mailvalidate"
	go_imap "github.com/emersion/go-imap"
	"github.com/jhillyerd/enmime"
	"github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/internal/models"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
)

// IMAPHandler processes events from IMAP sources
type IMAPHandler struct {
	repositories       *repository.Repositories
	eventService       *events.EventsService
	emailFilterService interfaces.EmailFilterService
}

// NewIMAPHandler creates a new IMAP email handler
func NewIMAPHandler(
	repositories *repository.Repositories,
	eventService *events.EventsService,
	emailFilterService interfaces.EmailFilterService,
) *IMAPHandler {
	return &IMAPHandler{
		repositories:       repositories,
		eventService:       eventService,
		emailFilterService: emailFilterService,
	}
}

// Handle processes an IMAP email event
func (h *IMAPHandler) Handle(ctx context.Context, event interfaces.MailEvent) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPHandler.Handle")
	defer span.Finish()

	switch msg := event.Message.(type) {
	case *go_imap.Message:
		err := h.processIMAPMessage(ctx, event.MailboxID, event.Folder, msg)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

	default:
		err := fmt.Errorf("  Unknown message type %T", msg)
		span.LogKV("messageId", event.MessageID)
		tracing.TraceErr(span, err)
		return
	}
}

// Main processing function
func (h *IMAPHandler) processIMAPMessage(ctx context.Context, mailboxID, folder string, msg *go_imap.Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IMAPHandler.processIMAPMessage")
	defer span.Finish()

	email := models.Email{
		MailboxID:  mailboxID,
		Provider:   h.determineProvider(ctx, mailboxID),
		Folder:     folder,
		ImapUID:    msg.Uid,
		ReceivedAt: utils.NowPtr(),
	}

	// Process envelope data
	h.processEnvelope(&email, msg.Envelope)

	// Process message content
	attachments := h.processMessageContent(&email, msg)

	err := h.emailFilterService.ScanEmail(ctx, &email)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Save the email entity to the database
	err = h.repositories.EmailRepository.Create(ctx, &email)
	if err != nil {
		err = errors.Wrap(err, "Error saving email")
		return err
	}

	// Create attachment records if any
	if email.HasAttachment && len(attachments) > 0 {
		h.processAttachments(email.ID, attachments)
	}

	// Publish event for downstream processing
	if email.Classification == enum.EmailOK {
		return h.eventService.Publisher.PublishFanoutEvent(ctx, email.ID, model.EMAIL, dto.EmailReceived{})
	}
	return nil
}

// Process envelope data - separate function for better readability
func (h *IMAPHandler) processEnvelope(email *models.Email, envelope *go_imap.Envelope) {
	if envelope == nil {
		return
	}

	// Basic envelope data
	if !envelope.Date.IsZero() {
		sentTime := envelope.Date
		email.SentAt = &sentTime
	}

	email.Subject = envelope.Subject
	email.InReplyTo = envelope.InReplyTo
	email.MessageID = envelope.MessageId

	// Extract References if available
	if envelope.InReplyTo != "" {
		// Many clients put References in the InReplyTo field separated by spaces
		references := strings.Split(envelope.InReplyTo, " ")
		email.References = pq.StringArray(references)
	}

	// Sender information
	if len(envelope.From) > 0 {
		sender := envelope.From[0]
		email.FromName = sender.PersonalName
		syntaxValidation := mailvalidate.ValidateEmailSyntax(sender.Address())
		if syntaxValidation.IsValid {
			email.FromAddress = syntaxValidation.CleanEmail
		}
	}

	// Recipients
	email.ToAddresses = h.convertAddressesToStringArray(envelope.To)
	email.CcAddresses = h.convertAddressesToStringArray(envelope.Cc)
	email.BccAddresses = h.convertAddressesToStringArray(envelope.Bcc)

	// Store raw envelope data for reference
	envelopeMap := make(map[string]interface{})
	envelopeMap["date"] = envelope.Date
	envelopeMap["subject"] = envelope.Subject
	envelopeMap["message_id"] = envelope.MessageId
	envelopeMap["in_reply_to"] = envelope.InReplyTo
	envelopeMap["from"] = h.addressesToMap(envelope.From)
	envelopeMap["to"] = h.addressesToMap(envelope.To)
	envelopeMap["cc"] = h.addressesToMap(envelope.Cc)
	envelopeMap["bcc"] = h.addressesToMap(envelope.Bcc)
	email.Envelope = models.JSONMap(envelopeMap)
}

// Process message content
func (h *IMAPHandler) processMessageContent(email *models.Email, msg *go_imap.Message) []map[string]interface{} {
	// Determine thread ID (could be based on References, Subject, etc.)
	email.ThreadID = h.determineThreadID(email)

	// Get the full message content
	fullMessageData := h.extractFullMessage(msg)

	if len(fullMessageData) > 0 {
		// Parse with enmime for better email parsing
		return h.parseWithEnmime(email, fullMessageData)
	} else {
		// Fallback to manual extraction
		return h.extractContentManually(email, msg)
	}
}

// Extract full message data
func (h *IMAPHandler) extractFullMessage(msg *go_imap.Message) []byte {
	var fullMessageBuffer bytes.Buffer

	for section, literal := range msg.Body {
		if section.Peek {
			continue // Skip PEEK sections to avoid duplicates
		}

		// Check if this is the full message section
		if len(section.Path) == 0 && section.Specifier == go_imap.EntireSpecifier {
			data, err := io.ReadAll(literal)
			if err == nil {
				fullMessageBuffer.Write(data)
				break
			}
		}
	}

	return fullMessageBuffer.Bytes()
}

// Parse message using enmime
func (h *IMAPHandler) parseWithEnmime(email *models.Email, messageData []byte) []map[string]interface{} {
	emailParser, err := enmime.ReadEnvelope(bytes.NewReader(messageData))
	if err != nil {
		log.Printf("Error parsing email with enmime: %v", err)
		return nil
	}

	// Extract headers
	headers := make(map[string]interface{})
	for _, key := range emailParser.GetHeaderKeys() {
		values := emailParser.GetHeaderValues(key)
		if len(values) > 0 {
			headers[key] = values
		}
	}
	email.RawHeaders = models.JSONMap(headers)

	// Extract body content
	email.BodyText = emailParser.Text
	email.BodyHTML = emailParser.HTML

	// Process attachments
	attachments := make([]map[string]interface{}, 0)

	// Regular attachments
	for _, attachment := range emailParser.Attachments {
		attachmentInfo := map[string]interface{}{
			"filename":     attachment.FileName,
			"content_type": attachment.ContentType,
			"disposition":  attachment.Disposition,
			"size":         len(attachment.Content),
			"content":      attachment.Content, // Include the actual content
		}
		attachments = append(attachments, attachmentInfo)
	}

	// Inline attachments
	for _, inlineAttachment := range emailParser.Inlines {
		attachmentInfo := map[string]interface{}{
			"filename":     inlineAttachment.FileName,
			"content_type": inlineAttachment.ContentType,
			"disposition":  "inline",
			"content_id":   inlineAttachment.ContentID,
			"size":         len(inlineAttachment.Content),
			"content":      inlineAttachment.Content,
		}
		attachments = append(attachments, attachmentInfo)
	}

	if len(attachments) > 0 {
		email.HasAttachment = true
	}
	return attachments
}

// Manual content extraction as fallback
func (h *IMAPHandler) extractContentManually(email *models.Email, msg *go_imap.Message) []map[string]interface{} {
	// Store body structure if available
	var attachments []map[string]interface{}
	if msg.BodyStructure != nil {
		bodyStructure := h.parseBodyStructure(msg.BodyStructure)
		email.BodyStructure = models.JSONMap(bodyStructure)

		// Check for attachments in body structure
		attachments = h.extractAttachmentsFromStructure(msg.BodyStructure)
		if len(attachments) > 0 {
			email.HasAttachment = true
		}
	}

	// Try to extract content from individual parts
	for section, literal := range msg.Body {
		sectionKey := fmt.Sprintf("%v", section)

		data, err := io.ReadAll(literal)
		if err != nil {
			continue
		}

		// Extract text and HTML content
		if strings.Contains(strings.ToLower(sectionKey), "text/plain") {
			email.BodyText = string(data)
		} else if strings.Contains(strings.ToLower(sectionKey), "text/html") {
			email.BodyHTML = string(data)
		}
	}
	return attachments
}

// Process attachments - create attachment records
func (h *IMAPHandler) processAttachments(emailID string, attachmentsData []map[string]interface{}) {
	for _, attachmentData := range attachmentsData {
		// Create attachment entity
		attachment := models.EmailAttachment{
			EmailID:     emailID,
			Filename:    attachmentData["filename"].(string),
			ContentType: attachmentData["content_type"].(string),
			Size:        attachmentData["size"].(int),
			IsInline:    attachmentData["disposition"] == "inline",
		}

		// Set ContentID for inline attachments
		if attachment.IsInline && attachmentData["content_id"] != nil {
			attachment.ContentID = attachmentData["content_id"].(string)
		}

		// Save attachment metadata
		if err := h.repositories.EmailAttachmentRepository.Create(context.Background(), &attachment); err != nil {
			log.Printf("Error saving attachment: %v", err)
			continue
		}

		// Upload attachment content to storage
		if content, ok := attachmentData["content"].([]byte); ok && len(content) > 0 {
			err := h.repositories.EmailAttachmentRepository.Store(
				context.Background(),
				&attachment,
				content,
			)
			if err != nil {
				log.Printf("Error storing attachment content: %v", err)
			}
		}
	}
}

// Helper function to convert addresses to string array
func (h *IMAPHandler) convertAddressesToStringArray(addresses []*go_imap.Address) pq.StringArray {
	if len(addresses) == 0 {
		return pq.StringArray{}
	}

	result := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		if addr.MailboxName != "" && addr.HostName != "" {
			emailAddr := addr.Address()
			validation := mailvalidate.ValidateEmailSyntax(emailAddr)
			if validation.IsValid {
				result = append(result, validation.CleanEmail)
			}
		}
	}

	return pq.StringArray(result)
}

// Helper to convert addresses to map for JSON storage
func (h *IMAPHandler) addressesToMap(addresses []*go_imap.Address) []map[string]string {
	result := make([]map[string]string, 0, len(addresses))

	for _, addr := range addresses {
		addrMap := map[string]string{
			"name":    addr.PersonalName,
			"address": addr.Address(),
		}
		result = append(result, addrMap)
	}

	return result
}

// Parse body structure recursively
func (h *IMAPHandler) parseBodyStructure(bs *go_imap.BodyStructure) map[string]interface{} {
	if bs == nil {
		return nil
	}

	result := make(map[string]interface{})

	result["mime_type"] = bs.MIMEType
	result["mime_subtype"] = bs.MIMESubType
	result["parameters"] = bs.Params
	result["id"] = bs.Id
	result["description"] = bs.Description
	result["encoding"] = bs.Encoding
	result["size"] = bs.Size
	result["lines"] = bs.Lines

	if bs.Disposition != "" {
		result["disposition"] = bs.Disposition
		result["disposition_params"] = bs.DispositionParams
	}

	if bs.Language != nil {
		result["language"] = bs.Language
	}

	if bs.Location != nil {
		result["location"] = bs.Location
	}

	if bs.MD5 != "" {
		result["md5"] = bs.MD5
	}

	if len(bs.Parts) > 0 {
		parts := make([]map[string]interface{}, 0, len(bs.Parts))
		for _, part := range bs.Parts {
			parts = append(parts, h.parseBodyStructure(part))
		}
		result["parts"] = parts
	}

	return result
}

// Extract attachments from body structure
func (h *IMAPHandler) extractAttachmentsFromStructure(bs *go_imap.BodyStructure) []map[string]interface{} {
	attachments := []map[string]interface{}{}

	// Check if this part is an attachment
	if bs.Disposition == "attachment" || bs.Disposition == "inline" {
		attachment := make(map[string]interface{})

		// Get filename
		filename := ""
		if bs.DispositionParams != nil {
			if fname, ok := bs.DispositionParams["filename"]; ok {
				filename = fname
			}
		}

		if filename == "" && bs.Params != nil {
			if name, ok := bs.Params["name"]; ok {
				filename = name
			}
		}

		if filename == "" {
			filename = fmt.Sprintf("attachment.%s", strings.ToLower(bs.MIMESubType))
		}

		attachment["filename"] = filename
		attachment["content_type"] = fmt.Sprintf("%s/%s", bs.MIMEType, bs.MIMESubType)
		attachment["size"] = bs.Size
		attachment["disposition"] = bs.Disposition

		attachments = append(attachments, attachment)
	}

	// Recursively check all parts
	if len(bs.Parts) > 0 {
		for _, part := range bs.Parts {
			partAttachments := h.extractAttachmentsFromStructure(part)
			attachments = append(attachments, partAttachments...)
		}
	}

	return attachments
}

// Determine thread ID based on message references or subject
func (h *IMAPHandler) determineThreadID(email *models.Email) string {
	// If we have a reference, use the first reference as thread ID
	if len(email.References) > 0 {
		return email.References[0]
	}

	// If we have an in-reply-to, use that
	if email.InReplyTo != "" {
		return email.InReplyTo
	}

	// Otherwise use the message ID itself (starting a new thread)
	return email.MessageID
}

// Determine email provider based on mailbox configuration
func (h *IMAPHandler) determineProvider(ctx context.Context, mailboxID string) enum.EmailProvider {
	// Get the mailbox configuration from repository
	mailbox, err := h.repositories.MailboxRepository.GetMailbox(ctx, mailboxID)
	if err != nil || mailbox == nil || mailbox.Provider == "" {
		return ""
	}

	return mailbox.Provider
}
