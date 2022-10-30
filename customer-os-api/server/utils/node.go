package utils

import "github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"

func GetPropsFromNode(node dbtype.Node) map[string]interface{} {
	return node.Props
}

func GetStringPropOrEmpty(props map[string]interface{}, key string) string {
	if props[key] != nil {
		return props[key].(string)
	}
	return ""
}
