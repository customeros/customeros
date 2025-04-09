package main

import (
	"fmt"
	"log"
	"os"

	"github.com/customeros/customeros/packages/server/scraper/internal/config"
	"github.com/customeros/customeros/packages/server/scraper/internal/database"
	"github.com/customeros/customeros/packages/server/scraper/internal/server"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run <command>")
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
	scraperDB, err := database.InitDatabase(&database.DatabaseConfig{
		DBName:          cfg.ScraperDatabaseConfig.DBName,
		Host:            cfg.ScraperDatabaseConfig.Host,
		Port:            cfg.ScraperDatabaseConfig.Port,
		User:            cfg.ScraperDatabaseConfig.User,
		Password:        cfg.ScraperDatabaseConfig.Password,
		MaxConn:         cfg.ScraperDatabaseConfig.MaxConn,
		MaxIdleConn:     cfg.ScraperDatabaseConfig.MaxIdleConn,
		ConnMaxLifetime: cfg.ScraperDatabaseConfig.ConnMaxLifetime,
		LogLevel:        cfg.ScraperDatabaseConfig.LogLevel,
	})
	if err != nil {
		log.Fatalf("Scraper database initialization failed: %v", err)
	}

	warehouseDB, err := database.InitDatabase(&database.DatabaseConfig{
		DBName:          cfg.DataWarehouseConfig.DBName,
		Host:            cfg.DataWarehouseConfig.Host,
		Port:            cfg.DataWarehouseConfig.Port,
		User:            cfg.DataWarehouseConfig.User,
		Password:        cfg.DataWarehouseConfig.Password,
		MaxConn:         cfg.DataWarehouseConfig.MaxConn,
		MaxIdleConn:     cfg.DataWarehouseConfig.MaxIdleConn,
		ConnMaxLifetime: cfg.DataWarehouseConfig.ConnMaxLifetime,
		LogLevel:        cfg.DataWarehouseConfig.LogLevel,
	})
	if err != nil {
		log.Fatalf("Openline database initialization failed: %v", err)
	}

	switch os.Args[1] {
	case "migrate":
		// Run Mailstack database migrations
		// err := repository.MigrateInboxDB(cfg.InboxDatabaseConfig, inboxDB)
		// if err != nil {
		// 	log.Fatalf("Inbox database migration failed: %v", err)
		// }
		// log.Println("Inbox database migration completed successfully")
		//
		// // Run DataWarehouse migrations
		// err = repository.MigrateDataWarehouse(cfg.DataWarehouseConfig, warehouseDB)
		// if err != nil {
		// 	log.Fatalf("Warehouse database migration failed: %v", err)
		// }
		// log.Println("Warehouse database migration completed successfully")

	case "server":
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println("Scraper starting up...")

		srv, err := server.NewServer(cfg, inboxDB, warehouseDB)
		if err != nil {
			log.Fatalf("Server setup failed: %v", err)
		}

		// Start the server
		err = srv.Run()
		if err != nil {
			log.Fatalf("Server startup failed: %v", err)
		}

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Usage: go run <command>")
		fmt.Println("Commands:")
		fmt.Println("  migrate   Run database migrations")
		fmt.Println("  server    Start the application server")
		os.Exit(1)
	}
}
