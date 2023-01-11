package utils

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"time"
)

type Cypher string

func CypherPtr(cypher Cypher) *Cypher {
	return &cypher
}

type DbNodesWithTotalCount struct {
	Nodes []*dbtype.Node
	Count int64
}

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

func ExtractSingleRecordFirstValue(result neo4j.Result, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	if record, err := result.Single(); err != nil {
		return nil, err
	} else {
		return record.Values[0], nil
	}
}

func ExtractSingleRecordFirstValueAsNode(result neo4j.Result, err error) (*dbtype.Node, error) {
	node, err := ExtractSingleRecordFirstValue(result, err)
	if err != nil {
		return nil, err
	}
	dbTypeNode := node.(dbtype.Node)
	return &dbTypeNode, err
}

func ExtractSingleRecordNodeAndRelationship(result neo4j.Result, err error) (*dbtype.Node, *dbtype.Relationship, error) {
	if err != nil {
		return nil, nil, err
	}
	if record, err := result.Single(); err != nil {
		return nil, nil, err
	} else {
		return NodePtr(record.Values[0].(dbtype.Node)), RelationshipPtr(record.Values[1].(dbtype.Relationship)), nil
	}
}

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

func GetIntPropOrMinusOne(props map[string]any, key string) int64 {
	if props[key] != nil {
		return props[key].(int64)
	}
	return -1
}

func GetInt64PropOrZero(props map[string]any, key string) int64 {
	if props[key] != nil {
		return props[key].(int64)
	}
	return 0
}

func GetIntPropOrNil(props map[string]any, key string) *int64 {
	if props[key] != nil {
		i := props[key].(int64)
		return &i
	}
	return nil
}

func GetBoolPropOrFalse(props map[string]any, key string) bool {
	if props[key] != nil {
		return props[key].(bool)
	}
	return false
}

func GetBoolPropOrNil(props map[string]any, key string) *bool {
	if props[key] != nil {
		b := props[key].(bool)
		return &b
	}
	return nil
}

func GetFloatPropOrNil(props map[string]any, key string) *float64 {
	if props[key] != nil {
		f := props[key].(float64)
		return &f
	}
	return nil
}

func GetTimePropOrNow(props map[string]any, key string) time.Time {
	if props[key] != nil {
		return props[key].(time.Time)
	}
	return time.Now().UTC()
}

func GetTimePropOrNil(props map[string]any, key string) *time.Time {
	if props[key] != nil {
		t := props[key].(time.Time)
		return &t
	}
	return nil
}
