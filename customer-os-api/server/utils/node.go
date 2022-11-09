package utils

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"time"
)

func GetPropsFromNode(node dbtype.Node) map[string]any {
	return node.Props
}

func GetPropsFromRelationship(rel dbtype.Relationship) map[string]any {
	return rel.Props
}

func GetStringPropOrEmpty(props map[string]any, key string) string {
	if props[key] != nil {
		return props[key].(string)
	}
	return ""
}

func GetStringPropOrNil(props map[string]any, key string) *string {
	if props[key] != nil {
		s := props[key].(string)
		return &s
	}
	return nil
}

func GetIntPropOrDefault(props map[string]any, key string) int64 {
	if props[key] != nil {
		return props[key].(int64)
	}
	return -1
}

func GetBoolPropOrFalse(props map[string]any, key string) bool {
	if props[key] != nil {
		return props[key].(bool)
	}
	return false
}

func GetTimePropOrNow(props map[string]any, key string) time.Time {
	if props[key] != nil {
		return props[key].(time.Time)
	}
	return time.Now().UTC()
}
