package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type OrganizationTypeRepository interface {
	Create(tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error)
	Update(tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error)
	Delete(tenant string, id string) error
	FindAll(tenant string) ([]*dbtype.Node, error)
	FindForOrganization(tenant, organizationId string) (*dbtype.Node, error)
}

type organizationTypeRepository struct {
	driver *neo4j.Driver
}

func NewOrganizationTypeRepository(driver *neo4j.Driver) OrganizationTypeRepository {
	return &organizationTypeRepository{
		driver: driver,
	}
}

func (r *organizationTypeRepository) Create(tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (t:Tenant {name:$tenant})" +
		" MERGE (t)<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:randomUUID()})" +
		" ON CREATE SET ot.name=$name, ot:%s" +
		" RETURN ot"

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, "OrganizationType_"+tenant),
			map[string]any{
				"tenant": tenant,
				"name":   organizationType.Name,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *organizationTypeRepository) Update(tenant string, organizationType *entity.OrganizationTypeEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType {id:$organizationId})
			SET ot.name=$name
			RETURN ot`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationType.Id,
				"name":           organizationType.Name,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *organizationTypeRepository) Delete(tenant string, organizationId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
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

func (r *organizationTypeRepository) FindAll(tenant string) ([]*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_TYPE_BELONGS_TO_TENANT]-(ot:OrganizationType)
			RETURN ot ORDER BY ot.name`,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
		}
	})
	organizationTypeDbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		organizationTypeDbNodes = append(organizationTypeDbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
	}
	return organizationTypeDbNodes, err
}

func (r *organizationTypeRepository) FindForOrganization(tenant, organizationId string) (*dbtype.Node, error) {
	session := utils.NewNeo4jReadSession(*r.driver)
	defer session.Close()

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		if queryResult, err := tx.Run(`
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})-[:IS_OF_TYPE]->(ot:OrganizationType)
			RETURN ot`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect()
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
