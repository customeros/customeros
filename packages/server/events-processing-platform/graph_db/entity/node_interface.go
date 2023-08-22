package entity

type Neo4jNode interface {
	Labels(tenant string) []string
}
