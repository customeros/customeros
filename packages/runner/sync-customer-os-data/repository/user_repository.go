package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"time"
)

type UserRepository interface {
	GetMatchedUserId(ctx context.Context, tenant string, user entity.UserData) (string, error)
	MergeUser(ctx context.Context, tenant string, syncDate time.Time, user entity.UserData) error
	MergeEmail(ctx context.Context, tenant string, user entity.UserData) error
	MergePhoneNumber(ctx context.Context, tenant string, user entity.UserData, phoneNumber entity.PhoneNumber) error
	GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error)
	GetUserIdById(ctx context.Context, tenant, id string) (string, error)
	GetUserIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	GetUserIdByExternalOwnerId(ctx context.Context, tenant, externalOwnerId, externalSystemId string) (string, error)
}

type userRepository struct {
	driver *neo4j.DriverWithContext
}

func NewUserRepository(driver *neo4j.DriverWithContext) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) GetMatchedUserId(ctx context.Context, tenant string, user entity.UserData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetMatchedUserId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u1:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e)
				OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u2:User)-[:HAS]->(m:Email)
					WHERE (m.rawEmail=$email OR m.email=$email) AND $email <> '' 
				with coalesce(u1, u2) as user
				where user is not null
				return user.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": user.ExternalSystem,
				"userExternalId": user.ExternalId,
				"email":          user.Email,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	userIDs := dbRecords.([]*db.Record)
	if len(userIDs) == 1 {
		return userIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *userRepository) MergeUser(ctx context.Context, tenant string, syncDate time.Time, user entity.UserData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.MergeUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	// Create new User if it does not exist
	// If User exists, and sourceOfTruth is acceptable then update it.
	//   otherwise create/update AlternateUser for incoming source, with a new relationship 'ALTERNATE'
	// Link User with Tenant
	query := "MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem}) " +
		" MERGE (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET u.createdAt=$createdAt, " +
		"				u.updatedAt=$updatedAt, " +
		"               u.name=$name, " +
		"               u.firstName=$firstName, " +
		"				u.lastName=$lastName, " +
		"				u.timezone=$timezone,  " +
		"				u.profilePhotoUrl=$profilePhotoUrl, " +
		"               u.source=$source, " +
		"				u.sourceOfTruth=$sourceOfTruth, " +
		"				u.appSource=$appSource, " +
		"               u:%s" +
		" ON MATCH SET " +
		"				u.name = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR u.name is null OR u.name = '' THEN $name ELSE u.name END, " +
		"				u.firstName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR u.firstName is null OR u.firstName = '' THEN $firstName ELSE u.firstName END, " +
		"				u.lastName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR u.lastName is null OR u.lastName = '' THEN $lastName ELSE u.lastName END, " +
		"				u.timezone = CASE WHEN u.sourceOfTruth=$sourceOfTruth  OR u.timezone is null OR u.timezone = '' THEN $timezone ELSE u.timezone END, " +
		"				u.profilePhotoUrl = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR u.profilePhotoUrl is null OR u.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE u.profilePhotoUrl END, " +
		"				u.updatedAt=$now " +
		" WITH u, ext " +
		" MERGE (u)-[r:IS_LINKED_WITH {externalId:$externalId}]->(ext) " +
		" ON CREATE SET r.externalIdSecond = $externalOwnerId, " +
		"				r.syncDate=$syncDate, " +
		"				r.externalSource=$externalSource " +
		" ON MATCH SET r.syncDate=$syncDate, r.externalSource=$externalSource " +
		" WITH u " +
		" FOREACH (x in CASE WHEN u.sourceOfTruth <> $sourceOfTruth THEN [u] ELSE [] END | " +
		"  MERGE (x)-[:ALTERNATE]->(alt:AlternateUser {source:$source, id:x.id}) " +
		"    SET alt.updatedAt=$now, alt.appSource=$appSource, alt.firstName=$firstName, alt.lastName=$lastName, alt.name=$name, alt.profilePhotoUrl=$profilePhotoUrl " +
		") RETURN u.id"

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "User_"+tenant),
			map[string]interface{}{
				"tenant":          tenant,
				"userId":          user.Id,
				"externalSystem":  user.ExternalSystem,
				"externalId":      user.ExternalId,
				"externalSource":  user.ExternalSourceTable,
				"externalOwnerId": user.ExternalOwnerId,
				"syncDate":        syncDate,
				"name":            user.Name,
				"firstName":       user.FirstName,
				"lastName":        user.LastName,
				"timezone":        user.Timezone,
				"profilePhotoUrl": user.ProfilePhotoUrl,
				"createdAt":       utils.TimePtrFirstNonNilNillableAsAny(user.CreatedAt),
				"updatedAt":       utils.TimePtrFirstNonNilNillableAsAny(user.UpdatedAt),
				"source":          user.ExternalSystem,
				"sourceOfTruth":   user.ExternalSystem,
				"appSource":       constants.AppSourceSyncCustomerOsData,
				"now":             utils.Now(),
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

func (r *userRepository) MergeEmail(ctx context.Context, tenant string, user entity.UserData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.MergeEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (e:Email {rawEmail: $email})-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				e.id=randomUUID(), " +
		"				e.createdAt=$now, " +
		"				e.updatedAt=$now, " +
		"				e.source=$source, " +
		"				e.sourceOfTruth=$sourceOfTruth, " +
		"				e.appSource=$appSource, " +
		"				e:%s " +
		" WITH DISTINCT u, e " +
		" MERGE (u)-[rel:HAS]->(e) " +
		" ON CREATE SET rel.primary=true, " +
		"				rel.label=$label " +
		" RETURN e.id"
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Email_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"userId":        user.Id,
				"email":         user.Email,
				"label":         "WORK",
				"source":        user.ExternalSystem,
				"sourceOfTruth": user.ExternalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"now":           utils.Now(),
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

func (r *userRepository) MergePhoneNumber(ctx context.Context, tenant string, user entity.UserData, phoneNumber entity.PhoneNumber) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.MergePhoneNumber")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant}) " +
		" MERGE (p:PhoneNumber {rawPhoneNumber: $phoneNumber})-[:PHONE_NUMBER_BELONGS_TO_TENANT]->(t) " +
		" ON CREATE SET " +
		"				p.id=randomUUID(), " +
		"				p.createdAt=$now, " +
		"				p.updatedAt=$now, " +
		"				p.source=$source, " +
		"				p.sourceOfTruth=$sourceOfTruth, " +
		"				p.appSource=$appSource, " +
		"				p:%s " +
		" WITH DISTINCT u, p " +
		" MERGE (u)-[rel:HAS]->(p) " +
		" ON CREATE SET rel.primary=$primary, " +
		"				rel.label=$label " +
		" RETURN p.id "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "PhoneNumber_"+tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"userId":        user.Id,
				"phoneNumber":   phoneNumber.Number,
				"label":         phoneNumber.Label,
				"primary":       phoneNumber.Primary,
				"source":        user.ExternalSystem,
				"sourceOfTruth": user.ExternalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"now":           utils.Now(),
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

func (r *userRepository) GetAllCrossTenantsNotSynced(ctx context.Context, size int) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetAllCrossTenantsNotSynced")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (u:User)--(t:Tenant)
 			WHERE (u.syncedWithEventStore is null or u.syncedWithEventStore=false)
			RETURN u, t.name limit $size`,
			map[string]any{
				"size": size,
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

func (r *userRepository) GetUserIdById(ctx context.Context, tenant, id string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUserIdById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
				return u.id order by u.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant": tenant,
			"userId": id,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *userRepository) GetUserIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUserIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return u.id order by u.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalId":       externalId,
			"externalSystemId": externalSystemId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *userRepository) GetUserIdByExternalOwnerId(ctx context.Context, tenant, externalOwnerId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUserIdByExternalOwnerId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User)-[:IS_LINKED_WITH {externalIdSecond:$externalOwnerId}]->(e)
				return u.id order by u.createdAt`

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalOwnerId":  externalOwnerId,
			"externalSystemId": externalSystemId,
		})
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}
