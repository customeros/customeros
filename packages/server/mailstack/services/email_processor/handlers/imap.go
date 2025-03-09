package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	go_imap "github.com/emersion/go-imap"
	"github.com/jhillyerd/enmime"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

// IMAPHandler processes events from IMAP sources
type IMAPHandler struct {
	eventService *events.EventsService
}

// NewIMAPHandler creates a new IMAP email handler
func NewIMAPHandler(eventService *events.EventsService) *IMAPHandler {
	return &IMAPHandler{
		eventService: eventService,
	}
}

// Handle processes an IMAP email event
func (h *IMAPHandler) Handle(event interfaces.MailEvent) {
	log.Printf("📩 NEW EMAIL EVENT RECEIVED - Mailbox: %s, Folder: %s",
		event.MailboxID,
		event.Folder)

	// Process based on message type
	switch msg := event.Message.(type) {
	case *go_imap.Message:
		h.processIMAPMessage(event.MailboxID, event.Folder, msg)
	default:
		log.Printf("  Unknown message type: %T", msg)
		log.Printf("  Raw message: %+v", msg)
	}
}

// processIMAPMessage handles a go-imap Message
func (h *IMAPHandler) processIMAPMessage(mailboxID, folder string, msg *go_imap.Message) {
	// Create a comprehensive emailData structure
	emailData := make(map[string]interface{})

	// Basic message metadata
	emailData["mailbox_id"] = mailboxID
	emailData["folder"] = folder
	emailData["event_type"] = "new"
	emailData["timestamp"] = time.Now().Format(time.RFC3339)
	emailData["uid"] = msg.Uid
	emailData["flags"] = msg.Flags
	emailData["size"] = msg.Size
	emailData["seq_num"] = msg.SeqNum

	// Process envelope information
	if msg.Envelope != nil {
		envelope := make(map[string]interface{})
		if !msg.Envelope.Date.IsZero() {
			envelope["date"] = msg.Envelope.Date.Format(time.RFC3339)
		}
		envelope["subject"] = msg.Envelope.Subject
		envelope["message_id"] = msg.Envelope.MessageId
		envelope["in_reply_to"] = msg.Envelope.InReplyTo

		// Process address fields
		envelope["from"] = formatAddressesForJSON(msg.Envelope.From)
		envelope["sender"] = formatAddressesForJSON(msg.Envelope.Sender)
		envelope["reply_to"] = formatAddressesForJSON(msg.Envelope.ReplyTo)
		envelope["to"] = formatAddressesForJSON(msg.Envelope.To)
		envelope["cc"] = formatAddressesForJSON(msg.Envelope.Cc)
		envelope["bcc"] = formatAddressesForJSON(msg.Envelope.Bcc)

		emailData["envelope"] = envelope
	}

	// Get the full message content
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
			}
			break
		}
	}

	// If we have the full message, parse it with enmime for better structure
	if fullMessageBuffer.Len() > 0 {
		emailParser, err := enmime.ReadEnvelope(&fullMessageBuffer)
		if err == nil {
			// Extract headers
			headers := make(map[string]interface{})
			for _, v := range emailParser.GetHeaderKeys() {
				headerVals := emailParser.GetHeaderValues(v)
				if len(headerVals) > 0 {
					headers[v] = headerVals
				}
			}
			emailData["headers"] = headers

			// Extract text and HTML parts
			if emailParser.Text != "" {
				emailData["body_text"] = emailParser.Text
			}
			if emailParser.HTML != "" {
				emailData["body_html"] = emailParser.HTML
			}

			// Process attachments
			attachments := make([]map[string]interface{}, 0)
			for _, attachment := range emailParser.Attachments {
				attachmentInfo := map[string]interface{}{
					"filename":     attachment.FileName,
					"content_type": attachment.ContentType,
					"disposition":  attachment.Disposition,
					"size":         len(attachment.Content),
				}
				attachments = append(attachments, attachmentInfo)
			}

			// Process inline attachments (like embedded images)
			for _, inlineAttachment := range emailParser.Inlines {
				attachmentInfo := map[string]interface{}{
					"filename":     inlineAttachment.FileName,
					"content_type": inlineAttachment.ContentType,
					"disposition":  "inline",
					"content_id":   inlineAttachment.ContentID,
					"size":         len(inlineAttachment.Content),
				}
				attachments = append(attachments, attachmentInfo)
			}

			if len(attachments) > 0 {
				emailData["attachments"] = attachments
				emailData["has_attachments"] = true
			} else {
				emailData["has_attachments"] = false
			}

			// Include the raw message for completeness
			emailData["raw_message"] = fullMessageBuffer.String()
		}
	} else {
		// Fallback to the original body structure approach if we don't have the full message
		if msg.BodyStructure != nil {
			emailData["body_structure"] = parseBodyStructure(msg.BodyStructure)

			// Determine if there are attachments from body structure
			attachments := extractAttachments(msg.BodyStructure)
			if len(attachments) > 0 {
				emailData["attachments"] = attachments
				emailData["has_attachments"] = true
			} else {
				emailData["has_attachments"] = false
			}
		}

		// Try to extract content from individual parts
		contentMap := make(map[string]interface{})
		for section, literal := range msg.Body {
			sectionKey := fmt.Sprintf("%v", section)

			data, err := io.ReadAll(literal)
			if err != nil {
				continue
			}

			// For text/plain parts, extract as body_text
			if strings.Contains(sectionKey, "TEXT") || strings.Contains(sectionKey, "text/plain") {
				emailData["body_text"] = string(data)
			}

			// For text/html parts, extract as body_html
			if strings.Contains(sectionKey, "HTML") || strings.Contains(sectionKey, "text/html") {
				emailData["body_html"] = string(data)
			}

			// Keep the section content in content map as well
			contentMap[sectionKey] = string(data)
		}

		if len(contentMap) > 0 {
			emailData["content"] = contentMap
		}
	}

	// Send to webhook
	h.eventService.Publisher.PublishFanoutEvent()
}

// Helper function to format addresses for JSON
func formatAddressesForJSON(addresses []*go_imap.Address) []map[string]string {
	result := make([]map[string]string, 0, len(addresses))

	for _, addr := range addresses {
		addressMap := make(map[string]string)
		if addr.PersonalName != "" {
			addressMap["name"] = addr.PersonalName
		}
		addressMap["mailbox"] = addr.MailboxName
		addressMap["host"] = addr.HostName
		addressMap["address"] = fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)

		result = append(result, addressMap)
	}

	return result
}

// Helper function to recursively parse body structure into a map
func parseBodyStructure(bs *go_imap.BodyStructure) map[string]interface{} {
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
			parts = append(parts, parseBodyStructure(part))
		}
		result["parts"] = parts
	}

	return result
}

func extractAttachments(bs *go_imap.BodyStructure) []map[string]interface{} {
	attachments := []map[string]interface{}{}

	// Process this part if it's an attachment
	if bs.Disposition == "attachment" || bs.Disposition == "inline" {
		attachment := make(map[string]interface{})

		// Get filename from disposition parameters
		if bs.DispositionParams != nil {
			if filename, ok := bs.DispositionParams["filename"]; ok {
				attachment["filename"] = filename
			}
		}

		// If filename wasn't in disposition params, check content type params
		if _, hasFilename := attachment["filename"]; !hasFilename && bs.Params != nil {
			if name, ok := bs.Params["name"]; ok {
				attachment["filename"] = name
			}
		}

		// If we still don't have a filename, generate one based on content type
		if _, hasFilename := attachment["filename"]; !hasFilename {
			attachment["filename"] = fmt.Sprintf("attachment.%s", strings.ToLower(bs.MIMESubType))
		}

		attachment["mime_type"] = fmt.Sprintf("%s/%s", bs.MIMEType, bs.MIMESubType)
		attachment["size"] = bs.Size
		attachment["disposition"] = bs.Disposition

		attachments = append(attachments, attachment)
	}

	// Recursively check all parts
	if len(bs.Parts) > 0 {
		for _, part := range bs.Parts {
			partAttachments := extractAttachments(part)
			attachments = append(attachments, partAttachments...)
		}
	}

	return attachments
}
