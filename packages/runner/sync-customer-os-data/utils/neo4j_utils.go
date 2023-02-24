package utils

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewNeo4jReadSession(ctx context.Context, driver neo4j.DriverWithContext) neo4j.SessionWithContext {
	return newNeo4jSession(ctx, driver, neo4j.AccessModeRead)
}

func NewNeo4jWriteSession(ctx context.Context, driver neo4j.DriverWithContext) neo4j.SessionWithContext {
	return newNeo4jSession(ctx, driver, neo4j.AccessModeWrite)
}

func newNeo4jSession(ctx context.Context, driver neo4j.DriverWithContext, accessMode neo4j.AccessMode) neo4j.SessionWithContext {
	return driver.NewSession(ctx,
		neo4j.SessionConfig{
			AccessMode: accessMode,
			BoltLogger: neo4j.ConsoleBoltLogger(),
		},
	)
}
