package neo4j_repository

import (
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type AuthenticatedUserInTenant struct {
	Tenant              string   `json:"tenant"`
	AuthenticatedUserId string   `json:"authenticatedUserId"`
	UserId              string   `json:"userId"`
	UserFirstname       string   `json:"userFirstname"`
	UserLastname        string   `json:"userLastname"`
	UserPrimaryEmail    string   `json:"userPrimaryEmail"`
	Roles               []string `json:"roles"`
}

type UserReadRepository interface {
	GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
	GetUserById(ctx context.Context, tenant, userId string) (*dbtype.Node, error)
	FindAllUsersWithRolesByEmail(ctx context.Context, email string) ([]*AuthenticatedUserInTenant, error)
	GetCurrentTenantByUserEmail(ctx context.Context, email string) (string, error)
	FindPlatformOwners(ctx context.Context) ([]*AuthenticatedUserInTenant, error)
	FindFirstUserWithRolesByEmail(ctx context.Context, tenant, email string) (*AuthenticatedUserInTenant, error)
	FindTestUser(ctx context.Context) (*dbtype.Node, error)
	GetAuthenticatedUserInTenant(ctx context.Context, authUserId, email string) (*dbtype.Node, error)
	GetFirstUserByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error)
	GetAllOwnersForOrganizations(ctx context.Context, tenant string, organizationIDs []string) ([]*utils.DbNodeAndId, error)
	GetOwnerForOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error)
	GetOwnerForContact(ctx context.Context, tenant, contactId string) (*dbtype.Node, error)
	GetCreatorForNote(ctx context.Context, tenant, noteId string) (*dbtype.Node, error)
	GetPaginatedCustomerUsers(ctx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error)
	GetUsersByEmailIds(ctx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error)
	GetUsersByEmailAddresses(ctx context.Context, tenant string, emailAddresses []string) ([]*utils.DbNodeAndId, error)
	GetAllForPhoneNumbers(ctx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error)
	GetAllOwnersForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error)
	GetAllCreatorsForOpportunities(ctx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error)
	GetAllCreatorsForTasks(ctx context.Context, tenant string, taskIds []string) ([]*utils.DbNodeAndId, error)
	GetAllAssigneesForTasks(ctx context.Context, tenant string, taskIds []string) ([]*utils.DbNodeAndId, error)
	GetAllCreatorsForServiceLineItems(ctx context.Context, tenant string, serviceLineItemIds []string) ([]*utils.DbNodeAndId, error)
	GetAllCreatorsForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error)
	GetAllAuthorsForLogEntries(ctx context.Context, tenant string, logEntryIDs []string) ([]*utils.DbNodeAndId, error)
	GetAllAuthorsForComments(ctx context.Context, tenant string, commentIds []string) ([]*utils.DbNodeAndId, error)
	GetAllSendersForFlowSenders(ctx context.Context, tenant string, flowSenderIds []string) ([]*utils.DbNodeAndId, error)
	GetUsersConnectedForContacts(ctx context.Context, tenant string, contactsIds []string) ([]*utils.DbNodeAndId, error)
	GetDistinctOrganizationOwners(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetUsers(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
	GetOwnerForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetOwnerForReminder(ctx context.Context, tenant, reminderId string) (*dbtype.Node, error)
}

type userReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewUserReadRepository(driver *neo4j.DriverWithContext, database string) UserReadRepository {
	return &userReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *userReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *userReadRepository) GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetAllForTenant")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) RETURN u `
	params := map[string]any{
		"tenant": tenant,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbNodes := make([]*dbtype.Node, 0)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	spans.LogKV("result.count", len(dbNodes))
	return dbNodes, nil
}

func (r *userReadRepository) GetByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetByIds")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) where u.id in $ids RETURN u`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbNodes := make([]*dbtype.Node, 0)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodes = append(dbNodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	spans.LogKV("result.count", len(dbNodes))
	return dbNodes, nil
}

func (r *userReadRepository) GetUserById(ctx context.Context, tenant, userId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetUserById")
	defer spans.Finish()

	spans.LogKV("userId", userId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$id}) RETURN u`
	params := map[string]any{
		"tenant": tenant,
		"id":     userId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (u *userReadRepository) FindPlatformOwners(ctx context.Context) ([]*AuthenticatedUserInTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.FindPlatformOwners")
	defer spans.Finish()

	cypher := `MATCH (e:Email)-[:HAS{primary:true}]-(u:User)-[:AUTHENTICATED_BY]->(au:AuthenticationUser)-[:HAS_WORKSPACE]->(t:Tenant{name:"customerosai"}), (au)--(a:Authentication{provider:"google"})
			WHERE 'PLATFORM_OWNER' in u.roles
			RETURN t.name, au.id, u.id, u.roles, e.rawEmail, u.firstName, u.lastName ORDER BY u.createdAt ASC`
	params := map[string]any{}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *u.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if len(records.([]*neo4j.Record)) > 0 {

		var result []*AuthenticatedUserInTenant

		for _, record := range records.([]*neo4j.Record) {

			tenant := record.Values[0].(string)
			authenticatedUserId := record.Values[1].(string)
			userId := record.Values[2].(string)
			roleList, ok := record.Values[3].([]interface{})
			var roles []string
			if !ok {
				roles = []string{}
			} else {
				roles = u.toStringList(roleList)
			}
			rawEmail := record.Values[4].(string)
			userFirstname := record.Values[5].(string)
			userLastname := record.Values[6].(string)

			result = append(result, &AuthenticatedUserInTenant{
				Tenant:              tenant,
				AuthenticatedUserId: authenticatedUserId,
				UserId:              userId,
				UserPrimaryEmail:    rawEmail,
				UserFirstname:       userFirstname,
				UserLastname:        userLastname,
				Roles:               roles,
			})
		}

		return result, nil
	} else {
		return nil, nil
	}
}

func (u *userReadRepository) FindAllUsersWithRolesByEmail(ctx context.Context, email string) ([]*AuthenticatedUserInTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.FindAllUsersWithRoles")
	defer spans.Finish()

	cypher := `MATCH (e:Email)<-[:HAS]-(u:User)-[:AUTHENTICATED_BY]->(au:AuthenticationUser)-[:HAS_WORKSPACE]->(t:Tenant)
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN t.name, au.id, u.id, u.roles ORDER BY u.createdAt ASC`
	params := map[string]any{
		"email": email,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *u.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if len(records.([]*neo4j.Record)) > 0 {

		var result []*AuthenticatedUserInTenant

		for _, record := range records.([]*neo4j.Record) {

			tenant := record.Values[0].(string)
			authenticatedUserId := record.Values[1].(string)
			userId := record.Values[2].(string)
			roleList, ok := record.Values[3].([]interface{})
			var roles []string
			if !ok {
				roles = []string{}
			} else {
				roles = u.toStringList(roleList)
			}

			result = append(result, &AuthenticatedUserInTenant{
				Tenant:              tenant,
				AuthenticatedUserId: authenticatedUserId,
				UserId:              userId,
				Roles:               roles,
			})
		}

		return result, nil
	} else {
		return nil, nil
	}
}

func (u *userReadRepository) GetCurrentTenantByUserEmail(ctx context.Context, email string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetCurrentTenantByUserEmail")
	defer spans.Finish()

	spans.LogKV("email", email)

	cypher := `MATCH (e:Email)<-[:HAS]-(u:User)-[:AUTHENTICATED_BY]->(au:AuthenticationUser)-[:HAS_WORKSPACE]->(t:Tenant)
				WHERE toLower(e.email)=$email OR toLower(e.rawEmail)=$email
				WITH COALESCE(au.currentTenant, au.defaultTenant) as tenant
				WHERE tenant IS NOT NULL 
				RETURN tenant`
	params := map[string]interface{}{
		"email": strings.ToLower(email),
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *u.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	if len(records.([]string)) == 0 {
		spans.LogKV("result.found", false)
		return "", nil
	}
	spans.LogKV("result.found", true)
	spans.LogKV("result.tenant", records.([]string)[0])
	return records.([]string)[0], nil
}

func (u *userReadRepository) FindFirstUserWithRolesByEmail(ctx context.Context, tenant, email string) (*AuthenticatedUserInTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.FindFirstUserWithRolesByEmail")
	defer spans.Finish()

	spans.LogKV("email", email)
	spans.LogKV("tenant", tenant)

	cypher := fmt.Sprintf(`
			MATCH (e:Email_%s)<-[:HAS]-(u:User_%s)-[:AUTHENTICATED_BY]->(au:AuthenticationUser)-[:HAS_WORKSPACE]->(t:Tenant {name: $tenant})
			WHERE toLower(e.email)=$email OR toLower(e.rawEmail)=$email
			RETURN t.name, au.id, u.id, u.roles ORDER BY u.createdAt ASC LIMIT 1`, tenant, tenant)
	params := map[string]interface{}{
		"email":  strings.ToLower(email),
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *u.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(records.([]*neo4j.Record)))
	if len(records.([]*neo4j.Record)) > 0 {
		authenticatedUserId := records.([]*neo4j.Record)[0].Values[1].(string)
		userId := records.([]*neo4j.Record)[0].Values[2].(string)
		roleList, ok := records.([]*neo4j.Record)[0].Values[3].([]interface{})
		var roles []string
		if !ok {
			roles = []string{}
		} else {
			roles = u.toStringList(roleList)
		}
		output := AuthenticatedUserInTenant{
			Tenant:              tenant,
			AuthenticatedUserId: authenticatedUserId,
			UserId:              userId,
			Roles:               roles,
		}
		spans.LogObjectAsJson("result", output)
		return &output, nil
	} else {
		return nil, nil
	}
}

func (u *userReadRepository) toStringList(values []interface{}) []string {
	var result []string
	for _, value := range values {
		result = append(result, value.(string))
	}
	return result
}

func (r *userReadRepository) GetAuthenticatedUserInTenant(ctx context.Context, authUserId, email string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetAuthenticatedUserInTenant")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (e:Email_%s)<-[:HAS]-(u:User_%s)-[:%s]->(a:AuthenticationUser {id:$authUserId})-[:%s]->(t:Tenant {name:$tenant})
		WHERE toLower(e.email)=$email OR toLower(e.rawEmail)=$email
		RETURN u`, tenant, tenant, model.AUTHENTICATED_BY.String(), model.HAS_WORKSPACE.String())
	params := map[string]any{
		"tenant":     tenant,
		"email":      strings.ToLower(email),
		"authUserId": authUserId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) GetFirstUserByEmail(ctx context.Context, tenant, email string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetFirstUserByEmail")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email) 
			WHERE e.email=$email OR e.rawEmail=$email
			RETURN DISTINCT(u) ORDER by u.createdAt ASC limit 1`
	params := map[string]any{
		"tenant": tenant,
		"email":  email,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) FindTestUser(ctx context.Context) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.FindTestUser")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)
			WHERE u.test = true and u.firstName = "Test" and u.lastName = "Sender"
			RETURN u limit 1`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) GetAllOwnersForOrganizations(ctx context.Context, tenant string, organizationIDs []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetAllOwnersForOrganizations")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)<-[:OWNS]-(u:User)-[:USER_BELONGS_TO_TENANT]->(t)
			WHERE o.id IN $organizationIds
			RETURN u, o.id as orgId`
	params := map[string]any{
		"tenant":          tenant,
		"organizationIds": organizationIDs,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetOwnerForOrganization(ctx context.Context, tenant, organizationId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetOwnerForOrganization")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId})<-[:OWNS]-(u:User) RETURN u`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})

	if err != nil {
		return nil, err
	} else if len(result.([]*neo4j.Record)) == 0 {
		return nil, nil
	} else {
		return utils.NodePtr(result.([]*neo4j.Record)[0].Values[0].(dbtype.Node)), nil
	}
}

func (r *userReadRepository) GetOwnerForContact(parentCtx context.Context, tenant, contactId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetOwnerForContact")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})<-[:OWNS]-(u:User)
			RETURN u`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) GetCreatorForNote(parentCtx context.Context, tenant, noteId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetCreatorForNote")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:CREATED]->(n:Note {id:$noteId})
			RETURN u`
	params := map[string]any{
		"tenant": tenant,
		"noteId": noteId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) GetPaginatedCustomerUsers(parentCtx context.Context, tenant string, skip, limit int, filter *utils.CypherFilter, sort *utils.CypherSort) (*utils.DbNodesWithTotalCount, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetPaginatedCustomerUsers")
	defer spans.Finish()

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		filterCypherStr, filterParams := filter.CypherFilterFragment("u")
		countParams := map[string]any{
			"tenant": tenant,
		}
		utils.MergeMapToMap(filterParams, countParams)
		queryResult, err := tx.Run(ctx, fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) 
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
					WITH u
					%s
					RETURN u 
					%s 
					SKIP $skip LIMIT $limit`, filterCypherStr, sort.SortingCypherFragment("u")),
			params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	spans.LogKV("result.count", len(dbNodesWithTotalCount.Nodes))
	return dbNodesWithTotalCount, nil
}

func (r *userReadRepository) GetUsersByEmailIds(parentCtx context.Context, tenant string, emailIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetUsersByEmailIds")
	defer spans.Finish()

	spans.LogObjectAsJson("emailIds", emailIds)
	cypher := ` MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
			WHERE e.id IN $emailIds
			RETURN u, e.id as emailId ORDER BY u.firstName, u.lastName`
	params := map[string]any{
		"tenant":   tenant,
		"emailIds": emailIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetUsersByEmailAddresses(ctx context.Context, tenant string, emailAddresses []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetUsersByEmailAddresses")
	defer spans.Finish()
	spans.LogObjectAsJson("emailAddresses", emailAddresses)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:HAS]->(e:Email)-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t)
			WHERE toLower(e.email)IN $emailAddresses OR toLower(e.rawEmail) IN $emailAddresses 
			RETURN u, COALESCE(CASE WHEN e.email IS NOT NULL AND e.email <> '' THEN e.email ELSE e.rawEmail END, e.email) as emailAddress`
	params := map[string]any{
		"tenant":         tenant,
		"emailAddresses": emailAddresses,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *userReadRepository) GetAllForPhoneNumbers(parentCtx context.Context, tenant string, phoneNumberIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetAllForPhoneNumbers")
	defer spans.Finish()

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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

func (r *userReadRepository) GetAllOwnersForOpportunities(parentCtx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetAllOwnersForOpportunities")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:OWNS]->(op:Opportunity)
			WHERE op.id IN $opportunityIds
			RETURN u, op.id as opId`
	params := map[string]any{
		"tenant":         tenant,
		"opportunityIds": opportunityIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetAllCreatorsForOpportunities(parentCtx context.Context, tenant string, opportunityIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetAllCreatorsForOpportunities")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(op:Opportunity)
			WHERE op.id IN $opportunityIds
			RETURN u, op.id as opId`
	params := map[string]any{
		"tenant":         tenant,
		"opportunityIds": opportunityIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetAllCreatorsForServiceLineItems(parentCtx context.Context, tenant string, serviceLineItemIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetAllCreatorsForServiceLineItems")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(sli:ServiceLineItem)
			WHERE sli.id IN serviceLineItemIds
			RETURN u, sli.id as sliId`
	params := map[string]any{
		"tenant":             tenant,
		"serviceLineItemIds": serviceLineItemIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetAllCreatorsForContracts(parentCtx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetAllCreatorsForContracts")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(c:Contract)
			WHERE c.id IN $contractIds
			RETURN u, c.id as cId`
	params := map[string]any{
		"tenant":      tenant,
		"contractIds": contractIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetAllCreatorsForTasks(ctx context.Context, tenant string, taskIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetAllCreatorsForTasks")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(tsk:Task)
			WHERE tsk.id IN $taskIds
			RETURN u, tsk.id`
	params := map[string]any{
		"tenant":  tenant,
		"taskIds": taskIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetAllAssigneesForTasks(ctx context.Context, tenant string, taskIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetAllAssigneesForTasks")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:ASSIGNED_TO]-(tsk:Task)
			WHERE tsk.id IN $taskIds
			RETURN u, tsk.id`
	params := map[string]any{
		"tenant":  tenant,
		"taskIds": taskIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetAllAuthorsForLogEntries(parentCtx context.Context, tenant string, logEntryIDs []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetAllAuthorsForLogEntries")
	defer spans.Finish()

	spans.LogObjectAsJson("logEntryIDs", logEntryIDs)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(l:LogEntry_%s)
			WHERE l.id IN $logEntryIDs
			RETURN u, l.id as logEntryId`, tenant)
	spans.LogKV("cypher", cypher)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher,
			map[string]any{
				"tenant":      tenant,
				"logEntryIDs": logEntryIDs,
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

func (r *userReadRepository) GetAllAuthorsForComments(parentCtx context.Context, tenant string, commentIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetAllAuthorsForComments")
	defer spans.Finish()

	spans.LogObjectAsJson("commentIds", commentIds)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CREATED_BY]-(c:Comment_%s)
			WHERE c.id IN $commentIds
			RETURN u, c.id`, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"commentIds": commentIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetAllSendersForFlowSenders(ctx context.Context, tenant string, flowSenderIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetAllSendersForFlowSenders")
	defer spans.Finish()

	spans.LogObjectAsJson("flowSenderIds", flowSenderIds)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:HAS]-(fs:FlowSender_%s)
			WHERE fs.id IN $flowSenderIds
			RETURN u, fs.id`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"flowSenderIds": flowSenderIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetUsersConnectedForContacts(ctx context.Context, tenant string, contactsIds []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetUsersConnectedForContacts")
	defer spans.Finish()

	spans.LogObjectAsJson("contactsIds", contactsIds)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)<-[:CONNECTED_WITH]-(c:Contact_%s)
			WHERE c.id IN $contactsIds
			RETURN u, c.id`, tenant)
	params := map[string]any{
		"tenant":      tenant,
		"contactsIds": contactsIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
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

func (r *userReadRepository) GetDistinctOrganizationOwners(parentCtx context.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetDistinctOrganizationOwners")
	defer spans.Finish()

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:OWNS]->(:Organization)
			RETURN distinct(u) order by u.firstName, u.lastName`

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher,
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

func (r *userReadRepository) GetUsers(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserReadRepository.GetUsers")
	defer spans.Finish()

	spans.LogObjectAsJson("ids", ids)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User)
			WHERE u.id IN $ids
			RETURN u`
	spans.LogKV("cypher", cypher)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)
	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher,
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

func (r *userReadRepository) GetOwnerForContract(parentCtx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetOwnerForContract")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$contractId})<-[:OWNS]-(u:User)
			RETURN u`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}

func (r *userReadRepository) GetOwnerForReminder(parentCtx context.Context, tenant, reminderId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "UserReadRepository.GetOwnerForReminder")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$reminderId})-[:REMINDER_BELONGS_TO_USER]->(u:User)
			RETURN u`
	params := map[string]any{
		"tenant":     tenant,
		"reminderId": reminderId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil && err.Error() == "Result contains no more records" {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	spans.LogKV("result.found", true)
	return result.(*dbtype.Node), nil
}
