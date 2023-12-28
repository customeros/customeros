package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	UpdateUser(ctx context.Context, userId string, event events.UserUpdateEvent) error
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
