package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type PhoneNumberRepository interface {
	Exists(ctx context.Context, tenant, phoneNumber string) (bool, error)
	GetById(ctx context.Context, phoneNumberId string) (*dbtype.Node, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.Exists")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (p:PhoneNumber_%s) WHERE p.rawPhoneNumber = $phoneNumber OR p.e164 = $phoneNumber RETURN p LIMIT 1"
	span.LogFields(log.String("query", query))

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
		return false, err
	}
	return result.(bool), err
}

func (r *phoneNumberRepository) GetByPhoneNumber(ctx context.Context, tenant, phoneNumber string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf("MATCH (t:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber_%s) "+
		"WHERE p.rawPhoneNumber = $phoneNumber OR p.e164 = $phoneNumber RETURN p ORDER BY p.createdAt LIMIT 1", tenant)
	span.LogFields(log.String("query", query))

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
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *phoneNumberRepository) GetById(ctx context.Context, phoneNumberId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	query := "MATCH (:Tenant {name:$tenant})<-[:PHONE_NUMBER_BELONGS_TO_TENANT]-(p:PhoneNumber {id:$phoneNumberId}) RETURN p"
	span.LogFields(log.String("query", query))

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
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
