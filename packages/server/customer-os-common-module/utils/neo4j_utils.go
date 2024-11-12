package utils

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/log"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
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

type RecordsWithTotalCount struct {
	Records []*db.Record
	Count   int64
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

type DbNodeAndTenant struct {
	Node   *dbtype.Node
	Tenant string
}

type DbPropsAndId struct {
	Props        map[string]any
	LinkedNodeId string
}

type DbNodePairAndId struct {
	Pair         Pair[*dbtype.Node, *dbtype.Node]
	LinkedNodeId string
}

type SessionConfigurationOption func(config *neo4j.SessionConfig)

func WithDatabaseName(databaseName string) SessionConfigurationOption {
	return func(config *neo4j.SessionConfig) {
		if databaseName != "" {
			config.DatabaseName = databaseName
		}
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

	var err error
	maxRetries := 3
	for retry := 0; retry < maxRetries; retry++ {
		sugarLogger := zap.L().With(
			zap.String("accessMode", accessModeStr),
			zap.String("retry", fmt.Sprintf("%d/%d", retry+1, maxRetries)),
		).Sugar()
		if err := ctx.Err(); errors.Is(err, context.Canceled) {
			sugarLogger.Errorf("(newNeo4jSession) Context is cancelled by calling the cancel function: %s", err.Error())
			panic(err)
		} else if errors.Is(err, context.DeadlineExceeded) {
			sugarLogger.Errorf("(newNeo4jSession) Context is cancelled by deadline exceeded: %s", err.Error())
			panic(err)
		} else if err != nil {
			sugarLogger.Errorf("(newNeo4jSession) Context is cancelled by another error: %s", err.Error())
			panic(err)
		}

		// Verify connectivity
		err = driver.VerifyConnectivity(ctx)
		if err == nil {
			sugarLogger.Info("(VerifyConnectivity) Successfully verified connectivity, creating new session")
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
		} else {
			sugarLogger.Errorf("(VerifyConnectivity) Failed to verify connectivity: %s", err.Error())
		}
	}
	if err != nil {
		zap.L().Sugar().Fatalf("(VerifyConnectivity) Failed to verify connectivity after all attempts: %s", err.Error())
	}
	return nil
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

func ExtractAllRecordsAsDbNodeAndTenant(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*DbNodeAndTenant, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*DbNodeAndTenant, 0)
	for _, v := range records {
		element := new(DbNodeAndTenant)
		element.Node = NodePtr(v.Values[0].(neo4j.Node))
		element.Tenant = v.Values[1].(string)
		output = append(output, element)
	}
	return output, nil
}

func ExtractAllRecordsAsDbPropsAndId(ctx context.Context, result neo4j.ResultWithContext, err error) ([]*DbPropsAndId, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*DbPropsAndId, 0)
	for _, v := range records {
		element := new(DbPropsAndId)
		element.Props = v.Values[0].(map[string]any)
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

func GetInt64PropOrDefault(props map[string]any, key string, defaultVal int64) int64 {
	if props[key] != nil {
		return props[key].(int64)
	}
	return defaultVal
}

func GetInt64PropOrNil(props map[string]any, key string) *int64 {
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

func GetBoolPropOrTrue(props map[string]any, key string) bool {
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
		switch v := props[key].(type) {
		case float64:
			return &v
		case int64:
			f := float64(v)
			return &f
		}
	}
	return nil
}

func GetFloatPropOrZero(props map[string]any, key string) float64 {
	if props[key] != nil {
		switch v := props[key].(type) {
		case float64:
			return v
		case int64:
			return float64(v)
		}
	}
	return float64(0)
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
		case neo4j.Date:
			t := v.Time()
			return &t
		case string:
			t, _ := UnmarshalDateTime(v)
			if t != nil {
				return t
			}
		}
	}
	return nil
}

func GetTimePropFromNeo4jOrZeroTime(valueToExtract any) time.Time {
	if valueToExtract != nil {
		valueOrNil := GetTimePropOrNil(map[string]any{
			"key": valueToExtract,
		}, "key")
		if valueOrNil != nil {
			return *valueOrNil
		} else {
			return ZeroTime()
		}
	}

	return ZeroTime()
}

func ExecuteWriteQueryOnDb(ctx context.Context, driver neo4j.DriverWithContext, database, cypher string, params map[string]any) error {
	session := NewNeo4jWriteSession(ctx, driver, WithDatabaseName(database))
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	return err
}

func ExecuteReadInTransaction(
	ctx context.Context,
	driver *neo4j.DriverWithContext,
	database string,
	tx *neo4j.ManagedTransaction,
	action func(tx neo4j.ManagedTransaction) (any, error),
) (any, error) {
	if tx != nil {
		return action(*tx)
	}

	session := NewNeo4jReadSession(ctx, *driver, WithDatabaseName(database))
	defer session.Close(ctx)

	// Otherwise, create a new transaction using ExecuteWrite
	r, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Pass the new transaction to the action
		result, err := action(tx)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	return r, err
}

func ExecuteWriteInTransaction(
	ctx context.Context,
	driver *neo4j.DriverWithContext,
	database string,
	tx *neo4j.ManagedTransaction,
	action func(tx neo4j.ManagedTransaction) (any, error),
) (any, error) {
	if tx != nil {
		return action(*tx)
	}

	// Otherwise, create a new transaction using ExecuteWrite
	session := NewNeo4jWriteSession(ctx, *driver, WithDatabaseName(database))
	defer session.Close(ctx)

	r, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Pass the new transaction to the action
		result, err := action(tx)
		if err != nil {
			return nil, err
		}
		return result, nil
	})

	return r, err
}

type TxWithPostCommit struct {
	Tx       *neo4j.ManagedTransaction
	onCommit []func(ctx context.Context) error
}

func NewTxWithPostCommit(tx *neo4j.ManagedTransaction) *TxWithPostCommit {
	return &TxWithPostCommit{
		Tx:       tx,
		onCommit: []func(ctx context.Context) error{},
	}
}

func ExecuteWriteInTransactionWithPostCommitActions(
	ctx context.Context,
	driver *neo4j.DriverWithContext,
	database string,
	txWithPostCommit *TxWithPostCommit,
	action func(tx *TxWithPostCommit) (any, error), // Updated to use *TxWithPostCommit
) (any, error) {
	// If a transaction is already provided, execute the action in it
	if txWithPostCommit != nil {
		return action(txWithPostCommit)
	}

	// Otherwise, create a new transaction using ExecuteWrite
	session := NewNeo4jWriteSession(ctx, *driver, WithDatabaseName(database))
	defer session.Close(ctx)

	var result any
	var err error

	// Execute the action within a managed transaction
	r, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Wrap the new transaction with post-commit functionality
		txWithPostCommit = NewTxWithPostCommit(&tx)

		// Execute the action with the wrapped transaction
		result, err = action(txWithPostCommit)
		if err != nil {
			return nil, err
		}

		return result, nil
	})

	// Only run post-commit actions if ExecuteWrite completes successfully
	if err == nil && txWithPostCommit != nil {
		if err := txWithPostCommit.runPostCommitActions(ctx); err != nil {
			return nil, err
		}
	}

	return r, err
}

// AddPostCommitAction adds a post-commit action (e.g., publish message to RabbitMQ)
func (t *TxWithPostCommit) AddPostCommitAction(action func(ctx context.Context) error) {
	t.onCommit = append(t.onCommit, action)
}

// RunPostCommitActions executes all post-commit actions after the transaction successfully commits.
func (t *TxWithPostCommit) runPostCommitActions(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Neo4jUtils.RunPostCommitActions")
	defer span.Finish()

	for _, action := range t.onCommit {
		err := action(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}
	return nil
}

func ExecuteWriteQuery(ctx context.Context, driver neo4j.DriverWithContext, cypher string, params map[string]any) error {
	session := NewNeo4jWriteSession(ctx, driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})

	return err
}

func ExecuteQueryInTx(ctx context.Context, tx neo4j.ManagedTransaction, query string, params map[string]any) error {
	_, err := tx.Run(ctx, query, params)
	return err
}

func ExecuteQuery(ctx context.Context, driver neo4j.DriverWithContext, database, cypher string, params map[string]any, traceError func(err error)) (*neo4j.EagerResult, error) {
	sugarLogger := zap.L().With(zap.String("database", database)).Sugar()
	sugarLogger.Infof("(ExecuteQuery): %s", cypher)

	if err := ctx.Err(); errors.Is(err, context.Canceled) {
		traceError(err)
		sugarLogger.Errorf("(newNeo4 - ExecuteQuery) Context is cancelled by calling the cancel function: %s", err.Error())
		panic(err)
	} else if errors.Is(err, context.DeadlineExceeded) {
		traceError(err)
		sugarLogger.Errorf("(newNeo4 - ExecuteQuery) Context is cancelled by deadline exceeded: %s", err.Error())
		panic(err)
	} else if err != nil {
		traceError(err)
		sugarLogger.Errorf("(newNeo4 - ExecuteQuery) Context is cancelled by another error: %s", err.Error())
		panic(err)
	}

	return neo4j.ExecuteQuery(ctx,
		driver,
		cypher,
		params,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(database),
		neo4j.ExecuteQueryWithBoltLogger(neo4j.ConsoleBoltLogger()))
}

func ToNeo4jDateAsAny(t *time.Time) any {
	if t == nil {
		return nil
	}
	return neo4j.DateOf(*t)
}
