package repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type PhoneNumberRepository interface {
	// Deprecated
	Exists(ctx context.Context, tenant, phoneNumber string) (bool, error)
	// Deprecated
	GetById(ctx context.Context, phoneNumberId string) (*dbtype.Node, error)
	// Deprecated
	GetByPhoneNumber(ctx context.Context, tenant, phoneNumber string) (*dbtype.Node, error)
}

type phoneNumberRepository struct {
	driver *neo4j.DriverWithContext
}

func NewPhoneNumberRepository(driver *neo4j.DriverWithContext) PhoneNumberRepository {
	return &phoneNumberRepository{
		driver: driver,
	}
}

func (r *phoneNumberRepository) Exists(ctx context.Context, tenant string, phoneNumber string) (bool, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberRepository.Exists")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "phoneNumber", phoneNumber)

	query := "MATCH (p:PhoneNumber_%s) WHERE p.rawPhoneNumber = $phoneNumber OR p.e164 = $phoneNumber RETURN p LIMIT 1"
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"phoneNumber": phoneNumber,
			}); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil

		}
	})
	if err != nil {
		spans.TraceError(err)
		return false, err
	}
	return result.(bool), err
}

func (r *phoneNumberRepository) GetByPhoneNumber(ctx context.Context, tenant, phoneNumber string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberRepository.GetByPhoneNumber")
	defer spans.Finish()
	spans.LogKV("tenant", tenant, "phoneNumber", phoneNumber)

	query := fmt.Sprintf("MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber_%s) "+
		"WHERE p.rawPhoneNumber = $phoneNumber OR p.e164 = $phoneNumber RETURN p ORDER BY p.createdAt LIMIT 1", tenant)
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":      tenant,
				"phoneNumber": phoneNumber,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *phoneNumberRepository) GetById(ctx context.Context, phoneNumberId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "PhoneNumberRepository.GetById")
	defer spans.Finish()
	spans.LogKV("phoneNumberId", phoneNumberId)

	query := "MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId}) RETURN p"
	spans.LogKV("query", query)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"phoneNumberId": phoneNumberId,
				"tenant":        common.GetTenantFromContext(ctx),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
