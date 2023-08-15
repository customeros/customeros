package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userId string, event events.UserCreateEvent) error
	UpdateUser(ctx context.Context, userId string, event events.UserUpdateEvent) error
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
	span.LogFields(log.String("userId", userId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(p:User:User_%s {id:$id}) 
		 ON CREATE SET 	p.name = $name,
						p.firstName = $firstName,
						p.lastName = $lastName,
						p.source = $source,
						p.sourceOfTruth = $sourceOfTruth,
						p.appSource = $appSource,
						p.createdAt = $createdAt,
						p.updatedAt = $updatedAt,
						p.internal = $internal,
						p.profilePhotoUrl = $profilePhotoUrl,
						p.syncedWithEventStore = true 
		 ON MATCH SET 	p.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              userId,
				"name":            event.Name,
				"firstName":       event.FirstName,
				"lastName":        event.LastName,
				"tenant":          event.Tenant,
				"source":          event.Source,
				"sourceOfTruth":   event.SourceOfTruth,
				"appSource":       event.AppSource,
				"createdAt":       event.CreatedAt,
				"updatedAt":       event.UpdatedAt,
				"internal":        event.Internal,
				"profilePhotoUrl": event.ProfilePhotoUrl,
			})
		return nil, err
	})
	return err
}

func (r *userRepository) UpdateUser(ctx context.Context, userId string, event events.UserUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserRepository.UpdateUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, event.Tenant)
	span.LogFields(log.String("userId", userId))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(p:User:User_%s {id:$id})
		 SET	p.name = $name,
				p.firstName = $firstName,
				p.lastName = $lastName,
				p.sourceOfTruth = $sourceOfTruth,
				p.updatedAt = $updatedAt,
				p.internal = $internal,
				p.profilePhotoUrl = $profilePhotoUrl,
				p.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":              userId,
				"tenant":          event.Tenant,
				"name":            event.Name,
				"firstName":       event.FirstName,
				"lastName":        event.LastName,
				"sourceOfTruth":   event.SourceOfTruth,
				"updatedAt":       event.UpdatedAt,
				"internal":        event.Internal,
				"profilePhotoUrl": event.ProfilePhotoUrl,
			})
		return nil, err
	})
	return err
}
