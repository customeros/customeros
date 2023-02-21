package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/net/context"
)

type OrganizationTypeRepository interface {
	Create(ctx context.Context, tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error)
	Update(ctx context.Context, tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, tenant string, id string) error
	FindAll(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	FindForOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
}

type organizationTypeRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationTypeRepository(driver *neo4j.DriverWithContext) OrganizationTypeRepository {
	return &organizationTypeRepository{
		driver: driver,
	}
}

func (r *organizationTypeRepository) Create(ctx context.Context, tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (t:Tenant {name:$tenant})" +
		" MERGE (t)<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:randomUUID()})" +
		" ON CREATE SET ot.name=$name, " +
		"				ot.createdAt=$now, " +
		"				ot.updatedAt=$now, " +
		"				ot:%s" +
		" RETURN ot"

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "OrganizationType_"+tenant),
			map[string]any{
				"tenant": tenant,
				"name":   organizationType.Name,
				"now":    utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *organizationTypeRepository) Update(ctx context.Context, tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:$organizationId})
			SET ot.name=$name, ot.updatedAt=$now
			RETURN ot`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationType.Id,
				"name":           organizationType.Name,
				"now":            utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *organizationTypeRepository) Delete(ctx context.Context, tenant string, organizationId string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[r:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:$organizationId})
			DELETE r, ot`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			})
		return nil, err
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func (r *organizationTypeRepository) FindAll(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType)
			RETURN ot ORDER BY ot.name`,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	organizationTypeDbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		organizationTypeDbNodes = append(organizationTypeDbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
	}
	return organizationTypeDbNodes, err
}

func (r *organizationTypeRepository) FindForOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})-[:IS_OF_TYPE]->(ot:OrganizationType)
			RETURN ot`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	} else if len(dbRecords.([]*neo4j.Record)) == 0 {
		return nil, nil
	} else {
		return utils.NodePtr(dbRecords.([]*neo4j.Record)[0].Values[0].(dbtype.Node)), nil
	}
}
