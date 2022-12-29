package neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"github.com/sirupsen/logrus"
)

func ExecuteWriteQuery(driver *neo4j.Driver, query string, params map[string]interface{}) {
	session := utils.NewNeo4jWriteSession(*driver)
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(query, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		logrus.Fatalf("Failed executing query: %s\n Error: %s", query, err)
	}
}

func ExecuteReadQueryWithSingleReturn(driver *neo4j.Driver, query string, params map[string]any) any {
	session := utils.NewNeo4jReadSession(*driver)
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		record, err := tx.Run(query, params)
		if err != nil {
			logrus.Fatalf("Error executing query %s", query)
		}
		return record.Single()
	})
	if err != nil {
		logrus.Fatalf("Error executing query %s", query)
	}
	return queryResult
}

func ExecuteReadQueryWithCollectionReturn(driver *neo4j.Driver, query string, params map[string]any) []*db.Record {
	session := utils.NewNeo4jReadSession(*driver)
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		records, err := tx.Run(query, params)
		if err != nil {
			logrus.Fatalf("Error executing query %s", query)
		}
		return records.Collect()
	})
	if err != nil {
		logrus.Fatalf("Error executing query %s", query)
	}
	return queryResult.([]*db.Record)
}
