package integration_tests

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"log"
)

func ExecuteWriteQuery(driver *neo4j.Driver, query string, params map[string]interface{}) {
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Fatalf("Error executing query %s", query)
	}
}

func ExecuteReadQueryWithSingleReturn(driver *neo4j.Driver, query string, params map[string]any) any {
	session := (*driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		record, err := tx.Run(query, params)
		if err != nil {
			log.Fatalf("Error executing query %s", query)
		}
		return record.Single()
	})
	if err != nil {
		log.Fatalf("Error executing query %s", query)
	}
	return queryResult
}
