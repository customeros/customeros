package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type JobRoleReadRepository interface {
	GetAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*dbtype.Node, error)
	GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) ([]*dbtype.Node, error)
	GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForUsers(ctx context.Context, tenant string, userIds []string) ([]*utils.DbNodeAndId, error)
	ExistsForContactAndOrganization(ctx context.Context, tenant, contactId, organizationId string) (bool, error)
	GetAllForContactWithOrganizationId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string) ([]*utils.DbNodeAndId, error)
	GetByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
	GetById(ctx context.Context, tenant string, id string) (*dbtype.Node, error)
	GetJobRoleForContactAndOrganization(ctx context.Context, tenant, contactId, organizationId string) (*dbtype.Node, error)
	GetJobRoleForContactWithoutOrganization(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
}

type jobRoleReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *jobRoleReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func NewJobRoleReadRepository(driver *neo4j.DriverWithContext, database string) JobRoleReadRepository {
	return &jobRoleReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *jobRoleReadRepository) GetAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*dbtype.Node, error) {
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

func (r *jobRoleReadRepository) GetAllForContacts(ctx context.Context, tenant string, contactIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForContacts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

func (r *jobRoleReadRepository) GetAllForUsers(ctx context.Context, tenant string, userIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForUsers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

func (r *jobRoleReadRepository) GetAllForOrganization(ctx context.Context, session neo4j.SessionWithContext, tenant, organizationId string) ([]*dbtype.Node, error) {
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

func (r *jobRoleReadRepository) GetAllForOrganizations(ctx context.Context, tenant string, organizationIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
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

func (r *jobRoleReadRepository) ExistsForContactAndOrganization(ctx context.Context, tenant, contactId, organizationId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.ExistsForContactAndOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  			(o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  			(c)-[:WORKS_AS]->(r:JobRole)-[:ROLE_IN]->(o) 
				RETURN r`
	params := map[string]interface{}{
		"contactId":      contactId,
		"organizationId": organizationId,
		"tenant":         tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return false, err
	}
	span.LogFields(log.Bool("result.found", len(records.([]*neo4j.Record)) > 0))
	return len(records.([]*neo4j.Record)) > 0, err
}

func (r *jobRoleReadRepository) GetAllForContactWithOrganizationId(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, contactId string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetAllForContactWithOrganizationId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  			(c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(o:Organization)
				RETURN j, o.id`
	params := map[string]interface{}{
		"contactId": contactId,
		"tenant":    tenant,
	}

	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	return result.([]*utils.DbNodeAndId), err
}

func (r *jobRoleReadRepository) GetByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetByIds")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (job:JobRole_%s) WHERE job.id IN $ids RETURN job`, tenant)
	params := map[string]interface{}{
		"ids": ids,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		span.LogFields(log.Int("result.count", 0))
		return nil, err
	}
	nodes := result.([]*dbtype.Node)
	span.LogFields(log.Int("result.count", len(nodes)))
	return nodes, err
}

func (r *jobRoleReadRepository) GetJobRoleForContactAndOrganization(ctx context.Context, tenant, contactId, organizationId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetJobRoleForContactAndOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  	(o:Organization {id:$organizationId})-[:ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  	(c)-[:WORKS_AS]->(j:JobRole)-[:ROLE_IN]->(o) 
				RETURN j LIMIT 1`
	params := map[string]interface{}{
		"contactId":      contactId,
		"organizationId": organizationId,
		"tenant":         tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*dbtype.Node), nil
}

func (r *jobRoleReadRepository) GetJobRoleForContactWithoutOrganization(ctx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetJobRoleForContactWithoutOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  	(c)-[:WORKS_AS]->(j:JobRole)
				WHERE NOT (j)--(:Organization)
				RETURN j LIMIT 1`
	params := map[string]interface{}{
		"contactId": contactId,
		"tenant":    tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(*dbtype.Node), nil
}

func (r *jobRoleReadRepository) GetById(ctx context.Context, tenant string, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "JobRoleRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (job:JobRole_%s) WHERE job.id = $id RETURN job`, tenant)
	params := map[string]interface{}{
		"id": id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result.(*dbtype.Node), err
}
