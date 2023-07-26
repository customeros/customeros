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

type JobRoleRepository interface {
	GetAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*dbtype.Node, error)
	GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) ([]*dbtype.Node, error)
	GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	DeleteJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string) error
	SetOtherJobRolesForContactNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, skipRoleId string) error
	CreateJobRole(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, input entity.JobRoleEntity) (*dbtype.Node, error)
	UpdateJobRoleDetails(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string, input entity.JobRoleEntity) (*dbtype.Node, error)
	LinkWithOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant, roleId, organizationId string) error
	GetAllForUsers(ctx context.Context, tenant string, userIds []string) ([]*utils.DbNodeAndId, error)
}

type jobRoleRepository struct {
	driver *neo4j.DriverWithContext
}

func NewJobRoleRepository(driver *neo4j.DriverWithContext) JobRoleRepository {
	return &jobRoleRepository{
		driver: driver,
	}
}

func (r *jobRoleRepository) GetAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[:WORKS_AS]->(r:JobRole) 
				RETURN r ORDER BY r.jobTitle`,
			map[string]interface{}{
				"contactId": contactId,
				"tenant":    tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}

func (r *jobRoleRepository) GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForContacts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact)-[:WORKS_AS]->(job:JobRole)
			WHERE c.id IN $contactIds
			RETURN job, c.id as contactId ORDER BY job.jobTitle`,
			map[string]any{
				"tenant":     tenant,
				"contactIds": contactIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *jobRoleRepository) GetAllForUsers(ctx context.Context, tenant string, userIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForUsers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:WORKS_AS]->(job:JobRole)
			WHERE u.id IN $userIds
			RETURN job, u.id as userId ORDER BY job.jobTitle`,
			map[string]any{
				"tenant":  tenant,
				"userIds": userIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *jobRoleRepository) GetAllForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(org)<-[:ROLE_IN]-(r:JobRole) 
				RETURN r ORDER BY r.jobTitle`,
			map[string]interface{}{
				"organizationId": organizationId,
				"tenant":         tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	dbNodes := []*dbtype.Node{}
	for _, v := range records.([]*neo4j.Record) {
		if v.Values[0] != nil {
			dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(dbtype.Node)))
		}
	}
	return dbNodes, err
}

func (r *jobRoleRepository) GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization)<-[:ROLE_IN]-(job:JobRole)
			WHERE org.id IN $organizationIds
			RETURN job, org.id as organizationId ORDER BY job.jobTitle`,
			map[string]any{
				"tenant":          tenant,
				"organizationIds": organizationIds,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *jobRoleRepository) CreateJobRole(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, contactId string, input entity.JobRoleEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.CreateJobRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (c)-[:WORKS_AS]->(r:JobRole {id:randomUUID()}) " +
		" ON CREATE SET r.jobTitle=$jobTitle, " +
		"				r.primary=$primary, " +
		"				r.description=$description, " +
		"				r.company=$company, " +
		"				r.responsibilityLevel=$responsibilityLevel, " +
		"				r.source=$source, " +
		"				r.sourceOfTruth=$sourceOfTruth, " +
		"				r.appSource=$appSource, " +
		"				r.createdAt=$now, " +
		"				r.updatedAt=$now, " +
		"				r.startedAt=$startedAt, " +
		"				r.endedAt=$endedAt, " +
		"				r:%s " +
		" RETURN r"

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "JobRole_"+tenant),
		map[string]interface{}{
			"tenant":              tenant,
			"contactId":           contactId,
			"jobTitle":            input.JobTitle,
			"description":         input.Description,
			"company":             input.Company,
			"primary":             input.Primary,
			"responsibilityLevel": input.ResponsibilityLevel,
			"source":              input.Source,
			"sourceOfTruth":       input.SourceOfTruth,
			"appSource":           input.AppSource,
			"startedAt":           utils.TimePtrFirstNonNilNillableAsAny(input.StartedAt, utils.TimePtr(utils.Now())),
			"endedAt":             utils.TimePtrFirstNonNilNillableAsAny(input.EndedAt),
			"now":                 utils.Now(),
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *jobRoleRepository) UpdateJobRoleDetails(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string, input entity.JobRoleEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.UpdateJobRoleDetails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[:WORKS_AS]->(r:JobRole {id:$roleId})
			SET r.jobTitle=$jobTitle, 
				r.primary=$primary,
				r.description=$description,
				r.company=$company,
				r.responsibilityLevel=$responsibilityLevel,
				r.sourceOfTruth=$sourceOfTruth,
				r.startedAt=$startedAt,
				r.endedAt=$endedAt,
				r.updatedAt=$now
			RETURN r`,
		map[string]interface{}{
			"tenant":              tenant,
			"contactId":           contactId,
			"roleId":              roleId,
			"jobTitle":            input.JobTitle,
			"description":         input.Description,
			"company":             input.Company,
			"primary":             input.Primary,
			"responsibilityLevel": input.ResponsibilityLevel,
			"sourceOfTruth":       input.SourceOfTruth,
			"now":                 utils.Now(),
			"startedAt":           utils.TimePtrFirstNonNilNillableAsAny(input.StartedAt),
			"endedAt":             utils.TimePtrFirstNonNilNillableAsAny(input.EndedAt),
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *jobRoleRepository) DeleteJobRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, roleId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.DeleteJobRoleInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}),
					(c)-[:WORKS_AS]->(r:JobRole {id:$roleId})
			DETACH DELETE r`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"roleId":    roleId,
		})
	return err
}

func (r *jobRoleRepository) SetOtherJobRolesForContactNonPrimaryInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, skipRoleId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.SetOtherJobRolesForContactNonPrimaryInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				 (c)-[:WORKS_AS]->(r:JobRole)
			WHERE r.id <> $skipRoleId
            SET r.primary=false,
				r.updatedAt=datetime({timezone: 'UTC'})`,
		map[string]interface{}{
			"tenant":     tenant,
			"contactId":  contactId,
			"skipRoleId": skipRoleId,
		})
	return err
}

func (r *jobRoleRepository) LinkWithOrganization(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, roleId string, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := tx.Run(ctx, `
			MATCH (org:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(r:JobRole {id:$roleId})<-[:WORKS_AS]-(c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			OPTIONAL MATCH (r)-[rel:ROLE_IN]->(org2:Organization)
				WHERE org2.id <> org.id
			DELETE rel
			WITH r, org
			MERGE (r)-[:ROLE_IN]->(org)
			`,
		map[string]interface{}{
			"tenant":         tenant,
			"roleId":         roleId,
			"organizationId": organizationId,
		})
	return err
}
