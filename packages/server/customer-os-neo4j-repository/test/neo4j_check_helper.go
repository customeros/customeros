package test

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"reflect"
	"sort"
	"testing"
)

func AssertNeo4jNodeCount(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, nodes map[string]int) {
	for name, expectedCount := range nodes {
		actualCount := GetCountOfNodes(ctx, driver, name)
		require.Equal(t, expectedCount, actualCount, "Unexpected count for node: "+name)
	}
}

func AssertNeo4jRelationCount(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, relations map[string]int) {
	for name, expectedCount := range relations {
		actualCount := GetCountOfRelationships(ctx, driver, name)
		require.Equal(t, expectedCount, actualCount, "Unexpected count for relationship: "+name)
	}
}

func AssertRelationship(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, fromNodeId, relationshipType, toNodeId string) {
	rel, err := GetRelationship(ctx, driver, fromNodeId, toNodeId)
	require.Nil(t, err)
	require.NotNil(t, rel)
	require.Equal(t, relationshipType, rel.Type)
}

func AssertRelationships(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, fromNodeId string, relationshipTypes []string, toNodeId string) {
	rels, err := GetRelationships(ctx, driver, fromNodeId, toNodeId)
	require.Nil(t, err)
	require.NotNil(t, rels)
	require.Equal(t, len(relationshipTypes), len(rels))
	for _, rel := range rels {
		require.Contains(t, relationshipTypes, rel.Type)
	}
}

func AssertRelationshipWithProperties(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, fromNodeId, relationshipType, toNodeId string, expectedProperties map[string]any) {
	rel, err := GetRelationship(ctx, driver, fromNodeId, toNodeId)
	require.Nil(t, err)
	require.NotNil(t, rel)
	require.Equal(t, relationshipType, rel.Type)
	for k, v := range expectedProperties {
		require.Equal(t, v, rel.Props[k])
	}
}

func AssertNeo4jLabels(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, expectedLabels []string) {
	actualLabels := GetAllLabels(ctx, driver)
	sort.Strings(expectedLabels)
	sort.Strings(actualLabels)
	if !reflect.DeepEqual(actualLabels, expectedLabels) {
		t.Errorf("Expected labels: %v, \nActual labels: %v", expectedLabels, actualLabels)
	}
}

func GetCountOfNodes(ctx context.Context, driver *neo4j.DriverWithContext, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetNodeById(ctx context.Context, driver *neo4j.DriverWithContext, label, id string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *driver)
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (n:%s {id:$id}) RETURN n`, label),
			map[string]interface{}{
				"id": id,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	node := queryResult.(dbtype.Node)
	return &node, nil
}

func GetFirstNodeByLabel(ctx context.Context, driver *neo4j.DriverWithContext, label string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *driver)
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (n:%s) RETURN n limit 1`, label),
			map[string]interface{}{})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	node := queryResult.(dbtype.Node)
	return &node, nil
}

func GetRelationship(ctx context.Context, driver *neo4j.DriverWithContext, fromNodeId, toNodeId string) (*dbtype.Relationship, error) {
	session := utils.NewNeo4jReadSession(ctx, *driver)
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, `MATCH (n {id:$fromNodeId})-[rel]->(m {id:$toNodeId}) RETURN rel limit 1`,
			map[string]interface{}{
				"fromNodeId": fromNodeId,
				"toNodeId":   toNodeId,
			})
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	node := queryResult.(dbtype.Relationship)
	return &node, nil
}

func GetRelationships(ctx context.Context, driver *neo4j.DriverWithContext, fromNodeId, toNodeId string) ([]dbtype.Relationship, error) {
	session := utils.NewNeo4jReadSession(ctx, *driver)
	defer session.Close(ctx)

	queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, `MATCH (n {id:$fromNodeId})-[rel]->(m {id:$toNodeId}) RETURN rel`,
			map[string]interface{}{
				"fromNodeId": fromNodeId,
				"toNodeId":   toNodeId,
			})
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	var relationships []dbtype.Relationship
	for _, record := range queryResult.([]*db.Record) {
		relationships = append(relationships, record.Values[0].(dbtype.Relationship))
	}
	return relationships, nil
}

func GetCountOfRelationships(ctx context.Context, driver *neo4j.DriverWithContext, relationship string) int {
	query := fmt.Sprintf(`MATCH (a)-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetTotalCountOfNodes(ctx context.Context, driver *neo4j.DriverWithContext) int {
	query := `MATCH (n) RETURN count(n)`
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetAllLabels(ctx context.Context, driver *neo4j.DriverWithContext) []string {
	query := `MATCH (n) RETURN DISTINCT labels(n)`
	dbRecords := ExecuteReadQueryWithCollectionReturn(ctx, driver, query, map[string]any{})
	labels := []string{}
	for _, v := range dbRecords {
		for _, nodeLabels := range v.Values {
			for _, label := range nodeLabels.([]interface{}) {
				if !contains(labels, label.(string)) {
					labels = append(labels, label.(string))
				}
			}
		}
	}
	return labels
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
