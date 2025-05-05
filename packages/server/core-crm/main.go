package main

import (
	"log"

	"github.com/customeros/customeros/packages/server/core-crm/internal/config"
	"github.com/customeros/customeros/packages/server/core-crm/internal/database"
	"github.com/customeros/customeros/packages/server/core-crm/internal/server"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Config initialization failed: %v", err)
	}
	if cfg == nil {
		log.Fatalf("config is empty")
	}

	// Setup the data warehouse
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
		log.Fatalf("Warehouse database initialization failed: %v", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Core CRM starting up...")

	srv, err := server.NewServer(cfg, warehouseDB)
	if err != nil {
		log.Fatalf("Server setup failed: %v", err)
	}

	// Start the server
	err = srv.Run()
	if err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
