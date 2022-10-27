package utils

import "github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

func GetPropsFromNode(node dbtype.Node) map[string]interface{} {
	return node.Props
}
