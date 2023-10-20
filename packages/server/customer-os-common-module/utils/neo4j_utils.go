package utils

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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

type DbNodePairAndId struct {
	Pair         Pair[*dbtype.Node, *dbtype.Node]
	LinkedNodeId string
}

type SessionConfigurationOption func(config *neo4j.SessionConfig)

func WithDatabaseName(databaseName string) SessionConfigurationOption {
	return func(config *neo4j.SessionConfig) {
		config.DatabaseName = databaseName
	}
}

func WithBoltLogger(logger log.BoltLogger) SessionConfigurationOption {
	return func(config *neo4j.SessionConfig) {
		config.BoltLogger = logger
	}
}

func WithFetchSize(fetchSize int) SessionConfigurationOption {
	return func(config *neo4j.SessionConfig) {
		config.FetchSize = fetchSize
	}
}

func NewNeo4jReadSession(ctx context.Context, driver neo4j.DriverWithContext, options ...SessionConfigurationOption) neo4j.SessionWithContext {
	return newNeo4jSession(ctx, driver, neo4j.AccessModeRead, options...)
}

func NewNeo4jWriteSession(ctx context.Context, driver neo4j.DriverWithContext, options ...SessionConfigurationOption) neo4j.SessionWithContext {
	return newNeo4jSession(ctx, driver, neo4j.AccessModeWrite, options...)
}

func newNeo4jSession(ctx context.Context, driver neo4j.DriverWithContext, accessMode neo4j.AccessMode, options ...SessionConfigurationOption) neo4j.SessionWithContext {
	accessModeStr := "read"
	if accessMode == neo4j.AccessModeWrite {
		accessModeStr = "write"
	}

	if err := ctx.Err(); errors.Is(err, context.Canceled) {
		zap.L().With(
			zap.String("accessMode", accessModeStr),
			zap.String("ctxErr", err.Error()),
		).Sugar().Errorf("(VerifyConnectivity) Context is cancelled by calling the cancel function")
	} else if errors.Is(err, context.DeadlineExceeded) {
		zap.L().With(
			zap.String("accessMode", accessModeStr),
			zap.String("ctxErr", err.Error()),
		).Sugar().Errorf("(VerifyConnectivity) Context is cancelled by deadline exceeded")
	} else if err != nil {
		zap.L().With(
			zap.String("accessMode", accessModeStr),
			zap.String("ctxErr", err.Error()),
		).Sugar().Errorf("(VerifyConnectivity) Context is cancelled by another error")
	}

	err := driver.VerifyConnectivity(ctx)
	if err != nil {
		zap.L().With(
			zap.String("accessMode", accessModeStr),
		).Sugar().Fatalf("(VerifyConnectivity) Error connecting to Neo4j: %s", err.Error())
	}

	zap.L().With(zap.String("accessMode", accessModeStr)).Sugar().Info("(newNeo4jSession) Creating new session")
	sessionConfig := neo4j.SessionConfig{
		AccessMode: accessMode,
		BoltLogger: neo4j.ConsoleBoltLogger(),
	}
	for _, option := range options {
		option(&sessionConfig)
	}
	return driver.NewSession(
		ctx,
		sessionConfig,
	)
}

func ExtractFirstRecordFirstValueAsDbNodePtr(ctx context.Context, result neo4j.ResultWithContext, err error) (*dbtype.Node, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return NodePtr(records[0].Values[0].(dbtype.Node)), nil
}

func ExtractAllRecordsFirstValueAsDbNodePtrs(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*dbtype.Node, error) {
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

func ExtractAllRecordsAsDbNodeAndIdFromEagerResult(result *neo4j.EagerResult) []*DbNodeAndId {
	output := make([]*DbNodeAndId, 0)
	for _, v := range result.Records {
		element := new(DbNodeAndId)
		element.Node = NodePtr(v.Values[0].(neo4j.Node))
		element.LinkedNodeId = v.Values[1].(string)
		output = append(output, element)
	}
	return output
}

func ExtractAllRecordsAsDbNodePairAndId(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*DbNodePairAndId, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*DbNodePairAndId, 0)
	for _, v := range records {
		element := new(DbNodePairAndId)
		pair := Pair[*dbtype.Node, *dbtype.Node]{}
		if v.Values[0] == nil {
			pair.First = nil
		} else {
			pair.First = NodePtr(v.Values[0].(neo4j.Node))
		}
		if v.Values[1] == nil {
			pair.Second = nil
		} else {
			pair.Second = NodePtr(v.Values[1].(neo4j.Node))
		}
		element.Pair = pair
		element.LinkedNodeId = v.Values[2].(string)
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

func ExtractAllRecordsAsString(ctx context.Context, result neo4j.ResultWithContext, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]string, 0)
	for _, v := range records {
		output = append(output, v.Values[0].(string))
	}
	return output, nil
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

func ExtractSingleRecordFirstValueAsNode(ctx context.Context, result neo4j.ResultWithContext, err error) (*dbtype.Node, error) {
	node, err := ExtractSingleRecordFirstValue(ctx, result, err)
	if err != nil {
		return nil, err
	}
	dbTypeNode := node.(dbtype.Node)
	return &dbTypeNode, err
}

func ExtractSingleRecordAsNodeFromEagerResult(result *neo4j.EagerResult) (*dbtype.Node, error) {
	if len(result.Records) == 0 {
		return nil, errors.New("no records found")
	}
	if len(result.Records) > 1 {
		return nil, errors.New("more than one record found")
	}
	node := result.Records[0].Values[0].(dbtype.Node)
	return &node, nil
}

func ExtractSingleRecordFirstValueAsString(ctx context.Context, result neo4j.ResultWithContext, err error) (string, error) {
	value, err := ExtractSingleRecordFirstValue(ctx, result, err)
	if err != nil {
		return "", err
	}
	return value.(string), err
}

func ExtractSingleRecordFirstValueAsType[T any](ctx context.Context, result neo4j.ResultWithContext, err error) (T, error) {
	value, err := ExtractSingleRecordFirstValue(ctx, result, err)
	if err != nil {
		return *new(T), err
	}

	converted, ok := value.(T)
	if !ok {
		return *new(T), errors.New("invalid type")
	}

	return converted, nil
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
	if val, ok := props[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
		return ""
	}
	return ""
}

func GetStringPropOrNil(props map[string]any, key string) *string {
	if val, ok := props[key]; ok {
		if strVal, ok := val.(string); ok {
			return &strVal
		}
		return nil
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
	timePtr := GetTimePropOrNil(props, key)
	if timePtr != nil {
		return *timePtr
	}
	return GetEpochStart()
}

func GetTimePropOrZeroTime(props map[string]any, key string) time.Time {
	timePtr := GetTimePropOrNil(props, key)
	if timePtr != nil {
		return *timePtr
	}
	return ZeroTime()
}

func GetEpochStart() time.Time {
	return time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
}

func GetTimePropOrNow(props map[string]any, key string) time.Time {
	timePtr := GetTimePropOrNil(props, key)
	if timePtr != nil {
		return *timePtr
	}
	return time.Now().UTC()
}

func GetTimePropOrNil(props map[string]any, key string) *time.Time {
	if props[key] != nil {
		switch v := props[key].(type) {
		case time.Time:
			return &v
		case string:
			t, _ := UnmarshalDateTime(v)
			if t != nil {
				return t
			}
		}
	}
	return nil
}

func ExecuteWriteQuery(ctx context.Context, driver neo4j.DriverWithContext, query string, params map[string]any) error {
	session := NewNeo4jWriteSession(ctx, driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

func ExecuteQueryInTx(ctx context.Context, tx neo4j.ManagedTransaction, query string, params map[string]any) error {
	_, err := tx.Run(ctx, query, params)
	return err
}

func ExecuteQuery(ctx context.Context, driver neo4j.DriverWithContext, database, cypher string, params map[string]any) (*neo4j.EagerResult, error) {
	return neo4j.ExecuteQuery(ctx,
		driver,
		cypher,
		params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(database),
		neo4j.ExecuteQueryWithBoltLogger(neo4j.ConsoleBoltLogger()))
}
