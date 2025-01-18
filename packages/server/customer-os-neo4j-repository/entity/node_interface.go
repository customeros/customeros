package neo4j_entity

type Neo4jNode interface {
	Labels(tenant string) []string
}
