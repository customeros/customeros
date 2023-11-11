package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type ContractRepository interface {
	CreateContract(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.ContractEntity) (*dbtype.Node, error)
	SetContractCreator(ctx context.Context, tenant, userId, contractId string) error
}

type contractRepository struct {
	driver *neo4j.DriverWithContext
}

func NewContractRepository(driver *neo4j.DriverWithContext) ContractRepository {
	return &contractRepository{
		driver: driver,
	}
}

func (r *contractRepository) CreateContract(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity *entity.ContractEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.CreateContract")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MERGE (c:Contract_%s {id:randomUUID()}) " +
		" ON CREATE SET c:Contract, " +
		" 				c:TimelineEvent, " +
		" 				c:TimelineEvent_%s, " +
		"				c.name=$name, " +
		"				c.createdAt=$createdAt, " +
		"				c.updatedAt=$updatedAt, " +
		"				c.serviceStartedAt=serviceStartedAt, " +
		"				c.appSource=$appSource, " +
		"				c.source=$source, " +
		"				c.contractStatus=contractStatus " +
		"				c.contractRenewalCycle=contractRenewalCycle " +
		" RETURN m"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
		map[string]any{
			"name":                 entity.Name,
			"createdAt":            entity.CreatedAt,
			"updatedAt":            entity.CreatedAt,
			"serviceStartedAt":     entity.ServiceStartedAt,
			"appSource":            entity.AppSource,
			"source":               entity.Source,
			"contractStatus":       entity.ContractStatus,
			"contractRenewalCycle": entity.ContractRenewalCycle,
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *contractRepository) SetContractCreator(ctx context.Context, tenant, userId, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractRepository.SetContractCreator")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
		(c:Contract {id:$contractId})
		MERGE (u)-[:CREATED]->(c)`

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":     tenant,
				"userId":     userId,
				"contractId": contractId,
			})
		return nil, err
	})
	return err
}
