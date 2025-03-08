package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	go_imap "github.com/emersion/go-imap"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/services/imap"
)

func main() {
	// Configure logging to include timestamps and file information
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("MailStack starting up...")

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create IMAP service
	log.Println("Creating IMAP service...")
	imapService := imap.NewIMAPService()

	// Set up webhook handler for email events
	webhookURL := "https://webhook.site/9efaff8f-b23e-4874-9750-e0089cc092ab"
	log.Printf("Webhook URL configured: %s", webhookURL)

	webhookHandler := func(event interfaces.MailEvent) {
		log.Printf("📩 NEW EMAIL EVENT RECEIVED - Mailbox: %s, Folder: %s",
			event.MailboxID,
			event.Folder)

		// Log message details based on type
		switch msg := event.Message.(type) {
		case *go_imap.Message:
			if msg.Envelope != nil {
				log.Printf("  Subject: %s", msg.Envelope.Subject)
				if len(msg.Envelope.From) > 0 {
					log.Printf("  From: %s", msg.Envelope.From[0].Address())
				}
				log.Printf("  Message-ID: %s", msg.Envelope.MessageId)
			} else {
				log.Printf("  Message has no envelope: %+v", msg)
			}
		default:
			log.Printf("  Unknown message type: %T", msg)
			log.Printf("  Raw message: %+v", msg)
		}

		// Send to webhook
		log.Printf("🌐 Sending to webhook: %s", webhookURL)
		if err := sendToWebhook(webhookURL, event); err != nil {
			log.Printf("❌ Error sending to webhook: %v", err)
		} else {
			log.Printf("✅ Successfully sent webhook notification")
		}
	}

	// Register webhook handler
	log.Println("Registering event handler...")
	imapService.SetEventHandler(webhookHandler)

	// Hardcoded mailbox configuration
	mailboxConfig := interfaces.MailboxConfig{
		ID:       "test",
		Server:   "mail.hostedemail.com",
		Port:     993,
		Username: "test@testcustomeros.com",
		Password: "admin123!",
		Folders:  []string{"INBOX"},
		TLS:      true,
	}

	log.Printf("Adding mailbox configuration - Server: %s, Username: %s, Folders: %v",
		mailboxConfig.Server,
		mailboxConfig.Username,
		mailboxConfig.Folders)

	// Add mailbox to the service
	if err := imapService.AddMailbox(mailboxConfig); err != nil {
		log.Fatalf("❌ Failed to add mailbox: %v", err)
	}
	log.Println("✅ Mailbox configuration added successfully")

	// Start the IMAP service
	log.Println("Starting IMAP service...")
	if err := imapService.Start(ctx); err != nil {
		log.Fatalf("❌ Failed to start IMAP service: %v", err)
	}
	log.Println("✅ IMAP service started successfully")

	// Set up HTTP server for status checking
	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Status request received from %s", r.RemoteAddr)
		status := imapService.Status()
		w.Header().Set("Content-Type", "application/json")

		// Pretty print the status for logging
		statusJSON, _ := json.MarshalIndent(status, "", "  ")
		log.Printf("Current mailbox status: \n%s", string(statusJSON))

		// Send normal response
		json.NewEncoder(w).Encode(status)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ HTTP server error: %v", err)
		}
	}()
	log.Println("✅ HTTP server started successfully")

	log.Println("MailStack is now running. Press Ctrl+C to exit.")

	// Set up signal handling for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	<-stop
	log.Println("Shutdown signal received...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shut down HTTP server
	log.Println("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("❌ HTTP server shutdown error: %v", err)
	} else {
		log.Println("✅ HTTP server shut down successfully")
	}

	// Stop IMAP service
	log.Println("Stopping IMAP service...")
	if err := imapService.Stop(); err != nil {
		log.Printf("❌ IMAP service shutdown error: %v", err)
	} else {
		log.Println("✅ IMAP service stopped successfully")
	}

	log.Println("Shutdown complete")
}

// Function to send events to webhook
func sendToWebhook(webhookURL string, event interfaces.MailEvent) error {
	// Extract relevant information from the message
	var emailData map[string]interface{}

	// Parse the message based on its type
	switch msg := event.Message.(type) {
	case *go_imap.Message:
		// Extract information from go-imap Message
		emailData = extractEmailData(msg)
	default:
		emailData = map[string]interface{}{
			"mailbox_id": event.MailboxID,
			"folder":     event.Folder,
			"event_type": event.EventType,
			"message_id": event.MessageID,
			"raw":        fmt.Sprintf("%v", event.Message),
		}
	}

	// Add event metadata
	emailData["timestamp"] = time.Now().Format(time.RFC3339)

	// Convert to JSON
	data, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	// Create request
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook error: %s - %s", resp.Status, string(body))
	}

	log.Printf("Webhook notification sent successfully for message ID: %d", event.MessageID)
	return nil
}

// Helper function to extract data from IMAP message
func extractEmailData(msg *go_imap.Message) map[string]interface{} {
	data := map[string]interface{}{
		"uid":    msg.Uid,
		"seq_id": msg.SeqNum,
		"flags":  msg.Flags,
	}

	// Add envelope information if available
	if msg.Envelope != nil {
		env := msg.Envelope
		data["subject"] = env.Subject
		data["message_id"] = env.MessageId
		data["date"] = env.Date.Format(time.RFC3339)

		if len(env.From) > 0 {
			fromAddrs := make([]string, len(env.From))
			for i, addr := range env.From {
				fromAddrs[i] = addr.Address()
			}
			data["from"] = fromAddrs
		}

		if len(env.To) > 0 {
			toAddrs := make([]string, len(env.To))
			for i, addr := range env.To {
				toAddrs[i] = addr.Address()
			}
			data["to"] = toAddrs
		}
	}

	// Try to extract message body if available
	// Note: This is a simplified approach, full parsing would be more complex
	var section go_imap.BodySectionName
	r := msg.GetBody(&section)
	if r != nil {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(r); err == nil {
			data["body_preview"] = buf.String()[:min(1000, buf.Len())] // First 1000 chars
		}
	}

	return data
}

// func main() {
// 	// Parse command-line flags
// 	listenAddr := flag.String("addr", ":8080", "HTTP listen address")
// 	flag.Parse()
//
// 	// Create context for graceful shutdown
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	// Setup the database
// 	dbConfig := database.Config{
// 		Host:     "localhost",
// 		Port:     5432,
// 		User:     "mailstack_user",
// 		Password: "your_secure_password",
// 		DBName:   "mailstack",
// 		SSLMode:  "disable", // Use "require" in production
// 	}
//
// 	db, err := database.NewConnection(dbConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
//
// 	// Auto-migrate the db schema
// 	if err := db.AutoMigrate(&models.Mailbox{}, &models.MessageState{}); err != nil {
// 		log.Fatalf("Failed to migrate database schema: %v", err)
// 	}
//
// 	// Create repositories
// 	mailboxRepo := repository.NewMailboxRepository(db)
// 	messageStateRepo := repository.NewMessageStateRepository(db)
//
// 	// Create Services
// 	imapService := imap.NewIMAPService()
// 	mailService := mailbox.NewMailService(imapService)
//
// 	// Set up HTTP server
// 	mux := http.NewServeMux()
//
// 	// Create and set up API server
// 	mailServer := api.NewMailServer(mailService)
// 	mailServer.SetupRoutes(mux)
//
// 	// Create HTTP server
// 	server := &http.Server{
// 		Addr:    *listenAddr,
// 		Handler: mux,
// 	}
//
// 	// Start mail service
// 	if err := mailService.Start(ctx); err != nil {
// 		log.Fatalf("Failed to start mail service: %v", err)
// 	}
//
// 	// Start HTTP server in a goroutine
// 	go func() {
// 		log.Printf("Starting HTTP server on %s", *listenAddr)
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("HTTP server error: %v", err)
// 		}
// 	}()
//
// 	// Set up signal handling for graceful shutdown
// 	stop := make(chan os.Signal, 1)
// 	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
//
// 	// Wait for interrupt signal
// 	<-stop
// 	log.Println("Shutting down...")
//
// 	// Create shutdown context with timeout
// 	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer shutdownCancel()
//
// 	// Shut down HTTP server
// 	if err := server.Shutdown(shutdownCtx); err != nil {
// 		log.Printf("HTTP server shutdown error: %v", err)
// 	}
//
// 	// Stop mail service
// 	if err := mailService.Stop(); err != nil {
// 		log.Printf("Mail service shutdown error: %v", err)
// 	}
//
// 	log.Println("Shutdown complete")
// }
