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
		DBName:          cfg.CommonConfig.Infrastructure.DataWarehouseConfig.DBName,
		Host:            cfg.CommonConfig.Infrastructure.DataWarehouseConfig.Host,
		ReadPort:        cfg.CommonConfig.Infrastructure.DataWarehouseConfig.ReadPort,
		WritePort:       cfg.CommonConfig.Infrastructure.DataWarehouseConfig.WritePort,
		User:            cfg.CommonConfig.Infrastructure.DataWarehouseConfig.User,
		Password:        cfg.CommonConfig.Infrastructure.DataWarehouseConfig.Password,
		MaxConn:         cfg.CommonConfig.Infrastructure.DataWarehouseConfig.MaxConn,
		MaxIdleConn:     cfg.CommonConfig.Infrastructure.DataWarehouseConfig.MaxIdleConn,
		ConnMaxLifetime: cfg.CommonConfig.Infrastructure.DataWarehouseConfig.ConnMaxLifetime,
		LogLevel:        cfg.CommonConfig.Infrastructure.DataWarehouseConfig.LogLevel,
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
