package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	go_imap "github.com/emersion/go-imap"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
)

// IMAPHandler processes events from IMAP sources
type IMAPHandler struct {
	webhookURL string
}

// NewIMAPHandler creates a new IMAP email handler
func NewIMAPHandler(webhookURL string) *IMAPHandler {
	return &IMAPHandler{
		webhookURL: webhookURL,
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
	// Extract email details
	emailData := make(map[string]interface{})
	emailData["mailbox_id"] = mailboxID
	emailData["folder"] = folder
	emailData["event_type"] = "new"
	emailData["timestamp"] = time.Now().Format(time.RFC3339)

	if msg.Envelope != nil {
		log.Printf("  Subject: %s", msg.Envelope.Subject)
		emailData["subject"] = msg.Envelope.Subject

		if len(msg.Envelope.From) > 0 {
			fromAddr := msg.Envelope.From[0].Address()
			log.Printf("  From: %s", fromAddr)
			emailData["from"] = fromAddr
		}

		log.Printf("  Message-ID: %s", msg.Envelope.MessageId)
		emailData["message_id"] = msg.Envelope.MessageId
		emailData["date"] = msg.Envelope.Date.Format(time.RFC3339)
	} else {
		log.Printf("  Message has no envelope: %+v", msg)
	}

	// Send to webhook
	h.sendToWebhook(emailData)
}

// sendToWebhook sends the email data to the configured webhook
func (h *IMAPHandler) sendToWebhook(emailData map[string]interface{}) {
	log.Printf("🌐 Sending to webhook: %s", h.webhookURL)

	// Convert to JSON
	data, err := json.Marshal(emailData)
	if err != nil {
		log.Printf("❌ Error marshaling data: %v", err)
		return
	}

	// Create request
	req, err := http.NewRequest("POST", h.webhookURL, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("❌ Error creating request: %v", err)
		return
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Error sending to webhook: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode >= 400 {
		log.Printf("❌ Webhook error: %s", resp.Status)
		return
	}

	log.Printf("✅ Successfully sent webhook notification")
}
