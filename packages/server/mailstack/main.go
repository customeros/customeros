package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/customeros/customeros/packages/server/mailstack/config"
	"github.com/customeros/customeros/packages/server/mailstack/internal/database"
	"github.com/customeros/customeros/packages/server/mailstack/internal/repository"
	"github.com/customeros/customeros/packages/server/mailstack/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mailstack <command>")
		fmt.Println("Commands:")
		fmt.Println("  migrate   Run database migrations")
		fmt.Println("  server    Start the application server")
		os.Exit(1)
	}

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Config initialization failed: %v", err)
	}
	if cfg == nil {
		log.Fatalf("config is empty")
	}

	// Setup the databases
	mailstackDB, err := database.InitMailstackDatabase(&database.DatabaseConfig{
		DBName:          cfg.MailstackDatabaseConfig.DBName,
		Host:            cfg.MailstackDatabaseConfig.Host,
		Port:            cfg.MailstackDatabaseConfig.Port,
		User:            cfg.MailstackDatabaseConfig.User,
		Password:        cfg.MailstackDatabaseConfig.Password,
		MaxConn:         cfg.MailstackDatabaseConfig.MaxConn,
		MaxIdleConn:     cfg.MailstackDatabaseConfig.MaxIdleConn,
		ConnMaxLifetime: cfg.MailstackDatabaseConfig.ConnMaxLifetime,
		LogLevel:        cfg.MailstackDatabaseConfig.LogLevel,
		SSLMode:         cfg.MailstackDatabaseConfig.SSLMode,
	})
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	switch os.Args[1] {
	case "migrate":

		err := repository.MigrateDB(mailstackDB)
		if err != nil {
			log.Fatalf("Database migration failed: %v", err)
		}
		log.Println("Database migration completed successfully")

	case "server":

		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println("MailStack starting up...")

		server, err := server.NewServer(cfg, mailstackDB)
		if err != nil {
			log.Fatalf("Server setup failed: %v", err)
		}

		err = server.Initialize(context.Background())
		if err != nil {
			log.Fatalf("Server initialization failed: %v", err)
		}

		err = server.Run()
		if err != nil {
			log.Fatalf("Server startup failed: %v", err)
		}

		log.Println("Shutdown complete")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Usage: mailstack <command>")
		fmt.Println("Commands:")
		fmt.Println("  migrate   Run database migrations")
		fmt.Println("  server    Start the application server")
		os.Exit(1)
	}
}
