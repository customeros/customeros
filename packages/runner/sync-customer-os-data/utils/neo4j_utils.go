package utils

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func NewNeo4jReadSession(driver neo4j.Driver) neo4j.Session {
	return newNeo4jSession(driver, neo4j.AccessModeRead)
}

func NewNeo4jWriteSession(driver neo4j.Driver) neo4j.Session {
	return newNeo4jSession(driver, neo4j.AccessModeWrite)
}

func newNeo4jSession(driver neo4j.Driver, accessMode neo4j.AccessMode) neo4j.Session {
	return driver.NewSession(
		neo4j.SessionConfig{
			AccessMode: accessMode,
			BoltLogger: neo4j.ConsoleBoltLogger(),
		},
	)
}
