package main

import (
	"flag"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/config"
	server "github.com/openline-ai/openline-customer-os/platform/events-processing-platform/event-processor-server"
	"github.com/openline-ai/openline-customer-os/platform/events-processing-platform/logger"
	"log"
)

func main() {

	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName(server.GetMicroserviceName(cfg))
	appLogger.Fatal(server.NewServer(cfg, appLogger).Run())
}
