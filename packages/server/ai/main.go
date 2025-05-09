package main

import (
	"log"

	"github.com/customeros/customeros/packages/server/ai/internal/config"
	"github.com/customeros/customeros/packages/server/ai/internal/database"
	"github.com/customeros/customeros/packages/server/ai/internal/server"
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
		DBName:          cfg.DataWarehouse.DBName,
		Host:            cfg.DataWarehouse.Host,
		ReadPort:        cfg.DataWarehouse.ReadPort,
		WritePort:       cfg.DataWarehouse.WritePort,
		User:            cfg.DataWarehouse.User,
		Password:        cfg.DataWarehouse.Password,
		MaxConn:         cfg.DataWarehouse.MaxConn,
		MaxIdleConn:     cfg.DataWarehouse.MaxIdleConn,
		ConnMaxLifetime: cfg.DataWarehouse.ConnMaxLifetime,
		LogLevel:        cfg.DataWarehouse.LogLevel,
	})
	if err != nil {
		log.Fatalf("Warehouse database initialization failed: %v", err)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("AI service starting up...")

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
