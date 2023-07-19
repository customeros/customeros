package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"time"
)

type JobRoleRepository interface {
	MergeJobRole(ctx context.Context, tenant, contactId, jobTitle, organizationExternalId, externalSystemId string, contactCreatedAt time.Time) error
	RemoveOutdatedJobRoles(ctx context.Context, tenant, contactId, externalSystemId, organizationExternal string) error
}

type jobRoleRepository struct {
	driver *neo4j.DriverWithContext
}

func NewJobRoleRepository(driver *neo4j.DriverWithContext) JobRoleRepository {
	return &jobRoleRepository{
		driver: driver,
	}
}

func (r *jobRoleRepository) MergeJobRole(ctx context.Context, tenant, contactId, jobTitle, organizationExternalId, externalSystemId string, contactCreatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.MergeJobRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MATCH (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$organizationExternalId}]-(org:Organization) " +
		" MERGE (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(org) " +
		" ON CREATE SET j.id=randomUUID(), " +
		"				j.primary=true, " +
		"				j.source=$source, " +
		"				j.sourceOfTruth=$sourceOfTruth, " +
		"				j.appSource=$appSource, " +
		"				j.jobTitle=$jobTitle, " +
		"				j.createdAt=$now, " +
		"				j.updatedAt=$now, " +
		"				j.contactCreatedAt=$contactCreatedAt, " +
		"				j:%s " +
		" ON MATCH SET 	" +
		"				j.jobTitle = CASE WHEN j.sourceOfTruth=$sourceOfTruth THEN $jobTitle ELSE j.jobTitle END, " +
		"				j.updatedAt = CASE WHEN j.sourceOfTruth=$sourceOfTruth THEN $now ELSE j.updatedAt END " +
		" RETURN j"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "JobRole_"+tenant),
			map[string]interface{}{
				"tenant":                 tenant,
				"contactId":              contactId,
				"externalSystemId":       externalSystemId,
				"organizationExternalId": organizationExternalId,
				"source":                 externalSystemId,
				"sourceOfTruth":          externalSystemId,
				"appSource":              externalSystemId,
				"jobTitle":               jobTitle,
				"now":                    utils.Now(),
				"contactCreatedAt":       contactCreatedAt,
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *jobRoleRepository) RemoveOutdatedJobRoles(ctx context.Context, tenant, contactId, externalSystemId, organizationExternalId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.RemoveOutdatedJobRoles")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			MATCH (c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(:Organization)-[r:IS_LINKED_WITH]->(e:ExternalSystem {id:$externalSystemId})
			WHERE j.sourceOfTruth=$sourceOfTruth AND r.externalId<>$organizationExternalId
			DETACH DELETE j`,
			map[string]interface{}{
				"tenant":                 tenant,
				"contactId":              contactId,
				"externalSystemId":       externalSystemId,
				"organizationExternalId": organizationExternalId,
				"sourceOfTruth":          externalSystemId,
			})
		return nil, err
	})
	return err
}
