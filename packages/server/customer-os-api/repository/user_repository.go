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
	"github.com/opentracing/opentracing-go/log"
)

type UserRepository interface {
	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entity entity.UserEntity) (*dbtype.Node, error)
	Update(ctx context.Context, session neo4j.SessionWithContext, tenant string, entity entity.UserEntity) (*dbtype.Node, error)
	FindUserByEmail(ctx context.Context, session neo4j.SessionWithContext, tenant string, email string) (*dbtype.Node, error)
	IsOwner(parentCtx context.Context, tx neo4j.ManagedTransaction, tenant, userId string) (*bool, error)
	GetOwnerForContact(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string) (*dbtype.Node, error)
	GetCreatorForNote(ctx context.Context, tx neo4j.ManagedTransaction, tenant, noteId string) (*dbtype.Node, error)
	GetPaginatedNonInternalUsers(ctx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetById(ctx context.Context, session neo4j.SessionWithContext, tenant, userId string) (*dbtype.Node, error)
	GetAllForEmails(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error)
	GetAllOwnersForOrganizations(ctx context.Context, tenant string, organizationIDs []string) ([]*utils.DbNodeAndId, error)

	GetAllUserPhoneNumberRelationships(ctx context.Context, size int) ([]*neo4j.Record, error)
	GetAllUserEmailRelationships(ctx context.Context, size int) ([]*neo4j.Record, error)

	AddRole(ctx context.Context, session neo4j.SessionWithContext, tenant string, userId string, role string) (*dbtype.Node, error)
	DeleteRole(ctx context.Context, session neo4j.SessionWithContext, tenant string, userId string, role string) (*dbtype.Node, error)
	GetDistinctOrganizationOwners(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetUsers(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
}

type userRepository struct {
	driver *neo4j.DriverWithContext
}

func NewUserRepository(driver *neo4j.DriverWithContext) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) Create(parentCtx context.Context, tx neo4j.ManagedTransaction, tenant string, entity entity.UserEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
			MERGE (u:User {id: randomUUID()})-[:USER_BELONGS_TO_TENANT]->(t) 
		 	ON CREATE SET u.firstName=$firstName, 
						u.lastName=$lastName, 
						u.createdAt=$now, 
						u.updatedAt=$now, 
		 				u.source=$source, 
						u.sourceOfTruth=$sourceOfTruth, 
						u.appSource=$appSource, 
						u.roles = $roles, 
						u.timezone = $timezone, 
						u:User_%s
	 		RETURN u`, tenant)
	span.LogFields(log.String("query", query))

	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"tenant":        tenant,
			"firstName":     entity.FirstName,
			"lastName":      entity.LastName,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"appSource":     entity.AppSource,
			"roles":         entity.Roles,
			"timezone":      entity.Timezone,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *userRepository) Update(parentCtx context.Context, session neo4j.SessionWithContext, tenant string, entity entity.UserEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.Update")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) 
		 	SET u.firstName=$firstName, 
				u.lastName=$lastName, 
				u.timezone=$timezone,
				u.updatedAt=datetime({timezone: 'UTC'}), 
				u.sourceOfTruth=$sourceOfTruth 
		 	RETURN u`
	span.LogFields(log.String("query", query))

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"userId":        entity.Id,
				"tenant":        tenant,
				"firstName":     entity.FirstName,
				"lastName":      entity.LastName,
				"sourceOfTruth": entity.SourceOfTruth,
				"timezone":      entity.Timezone,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *userRepository) AddRole(parentCtx context.Context, session neo4j.SessionWithContext, tenant string, userId string, role string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.AddRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" SET 	u.roles = CASE " +
		" 		WHEN NOT $role IN u.roles THEN u.roles + $role " +
		" 		ELSE u.roles " +
		" 		END, " +
		"		u.updatedAt=datetime({timezone: 'UTC'}) " +
		" RETURN u"

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"role":   role,
				"userId": userId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *userRepository) DeleteRole(parentCtx context.Context, session neo4j.SessionWithContext, tenant string, userId string, role string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.DeleteRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" SET 	u.roles = [item IN u.roles WHERE item <> $role], " +
		"		u.updatedAt=datetime({timezone: 'UTC'}) " +
		" RETURN u"

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"role":   role,
				"userId": userId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *userRepository) FindUserByEmail(parentCtx context.Context, session neo4j.SessionWithContext, tenant string, email string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.FindUserByEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email) 
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN DISTINCT u limit 1`,
			map[string]any{
				"tenant": tenant,
				"email":  email,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *userRepository) IsOwner(parentCtx context.Context, tx neo4j.ManagedTransaction, tenant, userId string) (*bool, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.IsOwner")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)<-[:OWNS]-(u:User{id:$userId})
			RETURN count(o)`,
		map[string]any{
			"tenant": tenant,
			"userId": userId,
		}); err != nil {
		return nil, err
	} else {
		count, err := queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		isOwner := count.Values[0].(int64) > 0
		return &isOwner, nil
	}
}

func (r *userRepository) GetOwnerForContact(parentCtx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetOwnerForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})<-[:OWNS]-(u:User)
			RETURN u`,
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
		}); err != nil {
		return nil, err
	} else {
		dbRecords, err := queryResult.Collect(ctx)
		if err != nil {
			return nil, err
		} else if len(dbRecords) == 0 {
			return nil, nil
		} else {
			return utils.NodePtr(dbRecords[0].Values[0].(dbtype.Node)), nil
		}
	}
}

func (r *userRepository) GetCreatorForNote(parentCtx context.Context, tx neo4j.ManagedTransaction, tenant, noteId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetCreatorForNote")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:CREATED]->(n:Note {id:$noteId})
			RETURN u`,
		map[string]any{
			"tenant": tenant,
			"noteId": noteId,
		}); err != nil {
		return nil, err
	} else {
		dbRecords, err := queryResult.Collect(ctx)
		if err != nil {
			return nil, err
		} else if len(dbRecords) == 0 {
			return nil, nil
		} else {
			return utils.NodePtr(dbRecords[0].Values[0].(dbtype.Node)), nil
		}
	}
}

func (r *userRepository) GetPaginatedNonInternalUsers(parentCtx context.Context, session neo4j.SessionWithContext, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetPaginatedNonInternalUsers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("u")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(ctx, fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) 
											WHERE u.internal=false or u.internal is null
											WITH u
											%s RETURN count(u) as count`, filterCypherStr),
			countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		params := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}
		utils.MergeMapToMap(filterParams, params)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
			`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) 
					WHERE u.internal=false or u.internal is null
					WITH u
					%s
					RETURN u 
					%s 
					SKIP $skip LIMIT $limit`, filterCypherStr, sort.SortingCypherFragment("u")),
			params)
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *userRepository) GetById(parentCtx context.Context, session neo4j.SessionWithContext, tenant, userId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
			RETURN u`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *userRepository) GetAllForEmails(parentCtx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetAllForEmails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
			WHERE e.id IN $emailIds
			RETURN u, e.id as emailId ORDER BY u.firstName, u.lastName`,
			map[string]any{
				"tenant":   tenant,
				"emailIds": emailIds,
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

func (r *userRepository) GetAllForPhoneNumbers(parentCtx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetAllForPhoneNumbers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(p:PhoneNumber)-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t)
			WHERE p.id IN $phoneNumberIds
			RETURN u, p.id as phoneNumberId ORDER BY u.firstName, u.lastName`,
			map[string]any{
				"tenant":         tenant,
				"phoneNumberIds": phoneNumberIds,
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

func (r *userRepository) GetAllOwnersForOrganizations(parentCtx context.Context, tenant string, organizationIDs []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetAllOwnersForOrganizations")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)<-[:OWNS]-(u:User)-[:USER_BELONGS_TO_TENANT]->(t)
			WHERE o.id IN $organizationIds
			RETURN u, o.id as orgId`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":          tenant,
				"organizationIds": organizationIDs,
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

func (r *userRepository) GetAllUserPhoneNumberRelationships(parentCtx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetAllUserPhoneNumberRelationships")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(p:PhoneNumber)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false)
			RETURN rel, u.id, p.id, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*neo4j.Record), err
}

func (r *userRepository) GetAllUserEmailRelationships(parentCtx context.Context, size int) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetAllUserEmailRelationships")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[rel:HAS]->(e:Email)
 			WHERE (rel.syncedWithEventStore is null or rel.syncedWithEventStore=false)
			RETURN rel, u.id, e.id, t.name limit $size`,
			map[string]any{
				"size": size,
			}); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*neo4j.Record), err
}

func (r *userRepository) GetDistinctOrganizationOwners(parentCtx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "UserRepository.GetDistinctOrganizationOwners")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:OWNS]->(:Organization)
			RETURN distinct(u) order by u.firstName, u.lastName, u.name`

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*dbtype.Node), err
}

func (r *userRepository) GetUsers(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUsers")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)
			WHERE u.id IN $ids
			RETURN u`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return dbRecords.([]*dbtype.Node), err
}
