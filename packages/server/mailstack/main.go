package main

import (
	"log"
	"os"

	"github.com/customeros/customeros/packages/server/mailstack/server"
)

func main() {
	// Configure logging to include timestamps and file information
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("MailStack starting up...")

	mailstackServer, err := server.NewServer()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	err = mailstackServer.Run()
	if err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}

	log.Println("Shutdown complete")
}
