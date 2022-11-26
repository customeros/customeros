package config

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/gen"
	"log"
)

func NewPostgresClient(cfg *Config) (*gen.Client, error) {
	var connUrl = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s search_path=%s sslmode=disable",
		cfg.PostgresDb.Host, cfg.PostgresDb.Port, cfg.PostgresDb.User, cfg.PostgresDb.Name, cfg.PostgresDb.Pwd, cfg.PostgresDb.Schema)
	log.Printf("Connecting to postgres database: host=%s port=%d user=%s dbname=%s schema=%s",
		cfg.PostgresDb.Host, cfg.PostgresDb.Port, cfg.PostgresDb.User, cfg.PostgresDb.Name, cfg.PostgresDb.Schema)
	client, err := gen.Open("postgres", connUrl)
	return client, err
}
