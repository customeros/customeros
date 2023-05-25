package utils

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/sirupsen/logrus"
	"time"
)

type Cypher string

func CypherPtr(cypher Cypher) *Cypher {
	return &cypher
}

type PairDbNodesWithTotalCount struct {
	Pairs []*Pair[*dbtype.Node, *dbtype.Node]
	Count int64
}

type DbNodesWithTotalCount struct {
	Nodes []*dbtype.Node
	Count int64
}

type DbNodeWithRelationAndId struct {
	Node         *dbtype.Node
	Relationship *dbtype.Relationship
	LinkedNodeId string
}

type DbNodeWithRelationIdAndTenant struct {
	Node         *dbtype.Node
	Relationship *dbtype.Relationship
	LinkedNodeId string
	Tenant       string
}

type DbNodeAndRelation struct {
	Node         *dbtype.Node
	Relationship *dbtype.Relationship
}

type DbNodeAndId struct {
	Node         *dbtype.Node
	LinkedNodeId string
}

func NewNeo4jReadSession(ctx context.Context, driver neo4j.DriverWithContext) neo4j.SessionWithContext {
	return newNeo4jSession(ctx, driver, neo4j.AccessModeRead)
}

func NewNeo4jWriteSession(ctx context.Context, driver neo4j.DriverWithContext) neo4j.SessionWithContext {
	return newNeo4jSession(ctx, driver, neo4j.AccessModeWrite)
}

func newNeo4jSession(ctx context.Context, driver neo4j.DriverWithContext, accessMode neo4j.AccessMode) neo4j.SessionWithContext {
	err := driver.VerifyConnectivity(ctx)
	if err != nil {
		logrus.Fatalf("(VerifyConnectivity) Error connecting to Neo4j: %v", err)
	}
	logrus.Infof("(newNeo4jSession) Creating new session with access mode: %v", accessMode)
	session := driver.NewSession(
		ctx,
		neo4j.SessionConfig{
			AccessMode: accessMode,
			BoltLogger: neo4j.ConsoleBoltLogger(),
		},
	)
	return session
}

func ExtractSingleRecordFirstValue(ctx context.Context, result neo4j.ResultWithContext, err error) (any, error) {
	if err != nil {
		return nil, err
	}
	if record, err := result.Single(ctx); err != nil {
		return nil, err
	} else {
		return record.Values[0], nil
	}
}

func ExtractAllRecordsFirstValueAsNodePtrs(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*dbtype.Node, error) {
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range records {
		dbNodes = append(dbNodes, NodePtr(v.Values[0].(dbtype.Node)))
	}
	return dbNodes, nil
}

func ExtractAllRecordsAsDbNodeWithRelationAndId(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*DbNodeWithRelationAndId, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*DbNodeWithRelationAndId, 0)
	for _, v := range records {
		element := new(DbNodeWithRelationAndId)
		element.Node = NodePtr(v.Values[0].(neo4j.Node))
		element.Relationship = RelationshipPtr(v.Values[1].(neo4j.Relationship))
		element.LinkedNodeId = v.Values[2].(string)
		output = append(output, element)
	}
	return output, nil
}

func ExtractAllRecordsAsDbNodeWithRelationIdAndTenant(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*DbNodeWithRelationIdAndTenant, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*DbNodeWithRelationIdAndTenant, 0)
	for _, v := range records {
		element := new(DbNodeWithRelationIdAndTenant)
		element.Node = NodePtr(v.Values[0].(neo4j.Node))
		element.Relationship = RelationshipPtr(v.Values[1].(neo4j.Relationship))
		element.LinkedNodeId = v.Values[2].(string)
		element.Tenant = v.Values[3].(string)
		output = append(output, element)
	}
	return output, nil
}

func ExtractAllRecordsAsDbNodeAndId(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*DbNodeAndId, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*DbNodeAndId, 0)
	for _, v := range records {
		element := new(DbNodeAndId)
		element.Node = NodePtr(v.Values[0].(neo4j.Node))
		element.LinkedNodeId = v.Values[1].(string)
		output = append(output, element)
	}
	return output, nil
}

func ExtractAllRecordsAsDbNodeAndRelation(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*DbNodeAndRelation, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*DbNodeAndRelation, 0)
	for _, v := range records {
		element := new(DbNodeAndRelation)
		element.Node = NodePtr(v.Values[0].(neo4j.Node))
		element.Relationship = RelationshipPtr(v.Values[1].(neo4j.Relationship))
		output = append(output, element)
	}
	return output, nil
}

func ExtractSingleRecordFirstValueAsNode(ctx context.Context, result neo4j.ResultWithContext, err error) (*dbtype.Node, error) {
	node, err := ExtractSingleRecordFirstValue(ctx, result, err)
	if err != nil {
		return nil, err
	}
	dbTypeNode := node.(dbtype.Node)
	return &dbTypeNode, err
}

func ExtractSingleRecordFirstValueAsString(ctx context.Context, result neo4j.ResultWithContext, err error) (string, error) {
	value, err := ExtractSingleRecordFirstValue(ctx, result, err)
	if err != nil {
		return "", err
	}
	return value.(string), err
}

func ExtractSingleRecordNodeAndRelationship(ctx context.Context, result neo4j.ResultWithContext, err error) (*dbtype.Node, *dbtype.Relationship, error) {
	if err != nil {
		return nil, nil, err
	}
	if record, err := result.Single(ctx); err != nil {
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

func GetListStringPropOrEmpty(props map[string]any, key string) []string {
	if props[key] != nil {
		// print the type of the value
		fmt.Printf("type: %T\n", props[key])
		s := props[key].([]any)
		s2 := make([]string, len(s))
		for i, v := range s {
			s2[i] = v.(string)
		}
		return s2
	}
	return []string{}
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

func GetTimePropOrEpochStart(props map[string]any, key string) time.Time {
	if props[key] != nil {
		return props[key].(time.Time)
	}
	return GetEpochStart()
}

func GetEpochStart() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
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
