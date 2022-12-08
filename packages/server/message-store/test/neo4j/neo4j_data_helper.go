package neo4j

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

func CleanupAllData(driver *neo4j.Driver) {
	ExecuteWriteQuery(driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateConversation(driver *neo4j.Driver, conversationId string) {
	query := `MERGE (o:Conversation {id:$conversationId})`
	ExecuteWriteQuery(driver, query, map[string]any{
		"conversationId": conversationId,
	})
}

func GetCountOfNodes(driver *neo4j.Driver, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetCountOfRelationships(driver *neo4j.Driver, relationship string) int {
	query := fmt.Sprintf(`MATCH (a)-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := ExecuteReadQueryWithSingleReturn(driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}
