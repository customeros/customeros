package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type OrganizationRepository interface {
	GetOrganizationWithDomain(ctx context.Context, tx neo4j.ManagedTransaction, tenant, domainId string) (*dbtype.Node, error)
	CreateOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, name, source, sourceOfTruth, appSource string, date time.Time, hide bool) (*dbtype.Node, error)
	LinkDomainToOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, domainName, organizationId string) error
	GetOrganizationsLinkedToEmailsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, emailIdList []string) ([]string, error)
	UpdateLastTouchpointInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error
}

type organizationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewOrganizationRepository(driver *neo4j.DriverWithContext) OrganizationRepository {
	return &organizationRepository{
		driver: driver,
	}
}

func (r *organizationRepository) GetOrganizationWithDomain(ctx context.Context, tx neo4j.ManagedTransaction, tenant, domainName string) (*dbtype.Node, error) {
	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain{domain:$domainName}) RETURN o`

	queryResult, err := tx.Run(ctx, query, map[string]any{
		"tenant":     tenant,
		"domainName": domainName,
	})
	if err != nil {
		return nil, err
	}
	result, err := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	} else if err != nil {
		return nil, nil
	} else {
		return result, nil
	}
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, name, source, sourceOfTruth, appSource string, date time.Time, hide bool) (*dbtype.Node, error) {
	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:randomUUID()}) " +
		" ON CREATE SET org.createdAt=$now, " +
		"				org.updatedAt=$now, " +
		"               org.name=$name, " +
		"				org.source=$source, " +
		"				org.sourceOfTruth=$sourceOfTruth, " +
		"				org.appSource=$appSource, " +
		"				org.hide=$hide, " +
		"				org:%s " +
		" RETURN org"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Organization_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"name":          name,
			"source":        source,
			"sourceOfTruth": sourceOfTruth,
			"appSource":     appSource,
			"now":           date,
			"hide":          hide,
		})
	if err != nil {
		return nil, err
	}
	result, err := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (r *organizationRepository) LinkDomainToOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, domainName, organizationId string) error {
	query := "MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId}) " +
		" MATCH (d:Domain {domain:$domainName}) " +
		" MERGE (o)-[:HAS_DOMAIN]->(d)"

	_, err := tx.Run(ctx, query, map[string]interface{}{
		"tenant":         tenant,
		"domainName":     domainName,
		"organizationId": organizationId,
	})
	return err
}

func (r *organizationRepository) GetOrganizationsLinkedToEmailsInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, emailIdList []string) ([]string, error) {
	query := "match (e:Email_%s)-[:HAS]-(c:Contact)-[:WORKS_AS]-(j:JobRole)-[:ROLE_IN]-(o:Organization) " +
		" where e.id in $emailIdList " +
		" return distinct(o.id)"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
		map[string]interface{}{
			"tenant":      tenant,
			"emailIdList": emailIdList,
		})
	if err != nil {
		return nil, err
	}
	result, err := utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (r *organizationRepository) UpdateLastTouchpointInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, organizationId string, touchpointAt time.Time, touchpointId string) error {
	query := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})
		 SET org.lastTouchpointAt=$touchpointAt, org.lastTouchpointId=$touchpointId`

	_, err := tx.Run(ctx, query,
		map[string]any{
			"tenant":         tenant,
			"organizationId": organizationId,
			"touchpointAt":   touchpointAt,
			"touchpointId":   touchpointId,
		})

	return err
}
