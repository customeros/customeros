package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userId string, event events.UserCreateEvent) error
	CreateUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, userId string, event events.UserCreateEvent) error
	UpdateUser(ctx context.Context, userId string, event events.UserUpdateEvent) error
	GetUser(ctx context.Context, tenant, userId string) (*dbtype.Node, error)
	AddRole(ctx context.Context, tenant, userId, role string, timestamp time.Time) error
	RemoveRole(ctx context.Context, tenant, userId, role string, timestamp time.Time) error
}

type userRepository struct {
	driver *neo4j.DriverWithContext
}

func NewUserRepository(driver *neo4j.DriverWithContext) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, userId string, event events.UserCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.CreateUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateUserInTx(ctx, tx, userId, event)
	})
	return err
}

func (r *userRepository) CreateUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, userId string, event events.UserCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.CreateUserInTx")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("userId", userId), log.Object("event", event))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User:User_%s {id:$id}) 
		 ON CREATE SET 	u.name = $name,
						u.firstName = $firstName,
						u.lastName = $lastName,
						u.source = $source,
						u.sourceOfTruth = $sourceOfTruth,
						u.appSource = $appSource,
						u.createdAt = $createdAt,
						u.updatedAt = $updatedAt,
						u.internal = $internal,
						u.bot = $bot,
						u.profilePhotoUrl = $profilePhotoUrl,
						u.timezone = $timezone,
						u.syncedWithEventStore = true 
		 ON MATCH SET 	u.name = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.name is null OR u.name = '' THEN $name ELSE u.name END,
						u.firstName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.firstName is null OR u.firstName = '' THEN $firstName ELSE u.firstName END,
						u.lastName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.lastName is null OR u.lastName = '' THEN $lastName ELSE u.lastName END,
						u.timezone = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.timezone is null OR u.timezone = '' THEN $timezone ELSE u.timezone END,
						u.profilePhotoUrl = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.profilePhotoUrl is null OR u.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE u.profilePhotoUrl END,
						u.internal = $internal,
						u.bot = $bot,
						u.updatedAt = $updatedAt,
						u.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE u.sourceOfTruth END,
						u.syncedWithEventStore = true`, event.Tenant)
	span.LogFields(log.String("query", query))

	return utils.ExecuteQueryInTx(ctx, tx, query, map[string]any{
		"tenant":          event.Tenant,
		"id":              userId,
		"name":            event.Name,
		"firstName":       event.FirstName,
		"lastName":        event.LastName,
		"internal":        event.Internal,
		"bot":             event.Bot,
		"profilePhotoUrl": event.ProfilePhotoUrl,
		"timezone":        event.Timezone,
		"source":          helper.GetSource(event.SourceFields.Source),
		"sourceOfTruth":   helper.GetSourceOfTruth(event.SourceFields.SourceOfTruth),
		"appSource":       helper.GetAppSource(event.SourceFields.AppSource),
		"createdAt":       event.CreatedAt,
		"updatedAt":       event.UpdatedAt,
		"overwrite":       helper.GetSourceOfTruth(event.SourceFields.SourceOfTruth) == constants.SourceOpenline,
	})
}

func (r *userRepository) UpdateUser(ctx context.Context, userId string, event events.UserUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.UpdateUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("userId", userId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User:User_%s {id:$id})
		 SET	u.name = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.name is null OR u.name = '' THEN $name ELSE u.name END,
				u.firstName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.firstName is null OR u.firstName = '' THEN $firstName ELSE u.firstName END,
				u.lastName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.lastName is null OR u.lastName = '' THEN $lastName ELSE u.lastName END,
				u.timezone = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.timezone is null OR u.timezone = '' THEN $timezone ELSE u.timezone END,
				u.profilePhotoUrl = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.profilePhotoUrl is null OR u.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE u.profilePhotoUrl END,
				u.updatedAt = $updatedAt,
				u.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE u.sourceOfTruth END,
				u.syncedWithEventStore = true`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              userId,
				"tenant":          event.Tenant,
				"name":            event.Name,
				"firstName":       event.FirstName,
				"lastName":        event.LastName,
				"sourceOfTruth":   helper.GetSource(event.Source),
				"updatedAt":       event.UpdatedAt,
				"internal":        event.Internal,
				"bot":             event.Bot,
				"profilePhotoUrl": event.ProfilePhotoUrl,
				"timezone":        event.Timezone,
				"overwrite":       helper.GetSource(event.Source) == constants.SourceOpenline,
			})
		return nil, err
	})
	return err
}

func (r *userRepository) GetUser(ctx context.Context, tenant, userId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.GetUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("userId", userId))

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$id}) RETURN u`
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"id":     userId,
			}); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r *userRepository) AddRole(ctx context.Context, tenant, userId, role string, timestamp time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.AddRole")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("userId", userId), log.String("role", role))

	query := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
		 	SET u.roles = CASE
					WHEN u.roles IS NULL THEN [$role]
					ELSE CASE
		 				WHEN NOT $role IN u.roles THEN u.roles + $role 
		 				ELSE u.roles 
		 				END
					END, 
				u.updatedAt=$updatedAt`
	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":    tenant,
		"role":      role,
		"userId":    userId,
		"updatedAt": timestamp,
	})
}

func (r *userRepository) RemoveRole(ctx context.Context, tenant, userId, role string, timestamp time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.RemoveRole")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("userId", userId), log.String("role", role))

	query := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
		 	SET u.roles = [item IN u.roles WHERE item <> $role],
				u.updatedAt=$updatedAt`
	span.LogFields(log.String("query", query))

	return r.executeQuery(ctx, query, map[string]any{
		"tenant":    tenant,
		"role":      role,
		"userId":    userId,
		"updatedAt": timestamp,
	})
}

func (r *userRepository) executeQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
