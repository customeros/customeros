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
	GetOrganizationWithDomain(ctx context.Context, tx neo4j.ManagedTransaction, tenant, domainName string) (*dbtype.Node, error)
	CreateOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, name, relationship, stage, leadSource, source, sourceOfTruth, appSource string, date time.Time, hide bool) (*dbtype.Node, error)
	LinkDomainToOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, domainName, organizationId string) error
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
	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_DOMAIN]->(d:Domain{domain:$domainName}) RETURN o ORDER BY o.createdAt ASC LIMIT 1`

	queryResult, err := tx.Run(ctx, query, map[string]any{
		"tenant":     tenant,
		"domainName": domainName,
	})
	if err != nil {
		return nil, err
	}
	result, _ := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	return result, nil
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, name, relationship, stage, leadSource, source, sourceOfTruth, appSource string, date time.Time, hide bool) (*dbtype.Node, error) {
	query := "MATCH (t:Tenant {name:$tenant}) " +
		" MERGE (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:randomUUID()}) " +
		" ON CREATE SET org.createdAt=$now, " +
		"				org.updatedAt=$now, " +
		"               org.name=$name, " +
		"				org.source=$source, " +
		"				org.sourceOfTruth=$sourceOfTruth, " +
		"				org.appSource=$appSource, " +
		"				org.hide=$hide, " +
		"				org.relationship=$relationship, " +
		"				org.stage=$stage, " +
		"				org.leadSource=$leadSource, " +
		"				org.onboardingStatus=$onboardingStatus, " +
		"				org:%s " +
		" RETURN org"

	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Organization_"+tenant),
		map[string]interface{}{
			"tenant":           tenant,
			"name":             name,
			"source":           source,
			"sourceOfTruth":    sourceOfTruth,
			"appSource":        appSource,
			"now":              date,
			"hide":             hide,
			"relationship":     relationship,
			"stage":            stage,
			"leadSource":       leadSource,
			"onboardingStatus": "NOT_APPLICABLE",
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
