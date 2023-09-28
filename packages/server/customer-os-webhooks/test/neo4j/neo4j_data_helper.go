package neo4j

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/entity"
	"time"
)

func CleanupAllData(ctx context.Context, driver *neo4j.DriverWithContext) {
	ExecuteWriteQuery(ctx, driver, `MATCH (n) DETACH DELETE n`, map[string]any{})
}

func CreateFullTextBasicSearchIndexes(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := fmt.Sprintf("DROP INDEX basicSearchStandard_location_terms IF EXISTS")
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	query = fmt.Sprintf("CREATE FULLTEXT INDEX basicSearchStandard_location_terms IF NOT EXISTS FOR (n:State) ON EACH [n.name, n.code] " +
		"OPTIONS {  indexConfig: { `fulltext.analyzer`: 'standard', `fulltext.eventually_consistent`: true } }")
	ExecuteWriteQuery(ctx, driver, query, map[string]any{})

	ExecuteWriteQuery(ctx, driver, query, map[string]any{})
}

func CreateTenant(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MERGE (t:Tenant {name:$tenant})`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant": tenant,
	})
}

func CreateWorkspace(ctx context.Context, driver *neo4j.DriverWithContext, workspace string, provider string, tenant string) {
	query := `MATCH (t:Tenant {name: $tenant})
			  MERGE (t)-[:HAS_WORKSPACE]->(w:Workspace {name:$workspace, provider:$provider})`

	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"provider":  provider,
		"workspace": workspace,
	})
}

func CreateHubspotExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": "hubspot",
	})
}

func CreateSlackExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)
			SET e.externalSource=$externalSource`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": "slack",
		"externalSource":   "Slack",
	})
}

func CreateCalComExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) {
	query := `MATCH (t:Tenant {name:$tenant})
			MERGE (e:ExternalSystem {id:$externalSystemId})-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t)`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": "calcom",
	})
}

func LinkWithHubspotExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, entityId, externalId string, externalUrl, externalSource *string, syncDate time.Time) {
	LinkWithExternalSystem(ctx, driver, entityId, externalId, "hubspot", externalUrl, externalSource, syncDate)
}

func LinkWithSlackExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, entityId, externalId string, externalUrl, externalSource *string, syncDate time.Time) {
	LinkWithExternalSystem(ctx, driver, entityId, externalId, "slack", externalUrl, externalSource, syncDate)
}

func LinkWithExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, entityId, externalId, externalSystemId string, externalUrl, externalSource *string, syncDate time.Time) {
	query := `MATCH (e:ExternalSystem {id:$externalSystemId}), (n {id:$entityId})
			MERGE (n)-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e)
			ON CREATE SET rel.externalUrl=$externalUrl, rel.syncDate=$syncDate, rel.externalSource=$externalSource`
	ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"externalSystemId": externalSystemId,
		"entityId":         entityId,
		"externalId":       externalId,
		"externalUrl":      externalUrl,
		"syncDate":         syncDate,
		"externalSource":   externalSource,
	})
}

func CreateDefaultUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string) string {
	return CreateUser(ctx, driver, tenant, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateDefaultUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string) string {
	return CreateUserWithId(ctx, driver, tenant, userId, entity.UserEntity{
		FirstName:     "first",
		LastName:      "last",
		Source:        "openline",
		SourceOfTruth: "openline",
	})
}

func CreateUser(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, user entity.UserEntity) string {
	return CreateUserWithId(ctx, driver, tenant, "", user)
}

func CreateUserWithId(ctx context.Context, driver *neo4j.DriverWithContext, tenant, userId string, user entity.UserEntity) string {
	if userId == "" {
		userUuid, _ := uuid.NewRandom()
		userId = userUuid.String()
	}
	query := `
		MATCH (t:Tenant {name:$tenant})
			MERGE (u:User {
				  	id: $userId,
				  	firstName: $firstName,
				  	lastName: $lastName,
					profilePhotoUrl: $profilePhotoUrl,
					createdAt :datetime({timezone: 'UTC'}),
					source: $source,
					sourceOfTruth: $sourceOfTruth
				})-[:USER_BELONGS_TO_TENANT]->(t)
			SET u:User_%s, 
				u.roles=$roles,
				u.internal=$internal`
	ExecuteWriteQuery(ctx, driver, fmt.Sprintf(query, tenant), map[string]any{
		"tenant":          tenant,
		"userId":          userId,
		"firstName":       user.FirstName,
		"lastName":        user.LastName,
		"source":          user.Source,
		"sourceOfTruth":   user.SourceOfTruth,
		"roles":           user.Roles,
		"internal":        user.Internal,
		"profilePhotoUrl": user.ProfilePhotoUrl,
	})
	return userId
}

func GetCountOfNodes(ctx context.Context, driver *neo4j.DriverWithContext, nodeLabel string) int {
	query := fmt.Sprintf(`MATCH (n:%s) RETURN count(n)`, nodeLabel)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetCountOfRelationships(ctx context.Context, driver *neo4j.DriverWithContext, relationship string) int {
	query := fmt.Sprintf(`MATCH (a)-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetCountOfRelationshipsForNodeWithId(ctx context.Context, driver *neo4j.DriverWithContext, relationship, id string) int {
	query := fmt.Sprintf(`MATCH (a {id:$id})-[r:%s]-(b) RETURN count(distinct r)`, relationship)
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{
		"id": id,
	})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetTotalCountOfNodes(ctx context.Context, driver *neo4j.DriverWithContext) int {
	query := `MATCH (n) RETURN count(n)`
	result := ExecuteReadQueryWithSingleReturn(ctx, driver, query, map[string]any{})
	return int(result.(*db.Record).Values[0].(int64))
}

func GetAllLabels(ctx context.Context, driver *neo4j.DriverWithContext) []string {
	query := `MATCH (n) RETURN DISTINCT labels(n)`
	dbRecords := ExecuteReadQueryWithCollectionReturn(ctx, driver, query, map[string]any{})
	labels := []string{}
	for _, v := range dbRecords {
		for _, nodeLabels := range v.Values {
			for _, label := range nodeLabels.([]interface{}) {
				if !utils.Contains(labels, label.(string)) {
					labels = append(labels, label.(string))
				}
			}
		}
	}
	return labels
}
