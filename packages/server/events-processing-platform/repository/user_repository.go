package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/events"
)

type UserRepository interface {
	CreateUser(ctx context.Context, aggregateId string, event events.UserCreatedEvent) error
	UpdateUser(ctx context.Context, aggregateId string, event events.UserUpdatedEvent) error
}

type userRepository struct {
	driver *neo4j.DriverWithContext
}

func NewUserRepository(driver *neo4j.DriverWithContext) UserRepository {
	return &userRepository{
		driver: driver,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, aggregateId string, event events.UserCreatedEvent) error {
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
						p.syncedWithEventStore = true 
		 ON MATCH SET 	p.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            aggregateId,
				"name":          event.Name,
				"firstName":     event.FirstName,
				"lastName":      event.LastName,
				"tenant":        event.Tenant,
				"source":        event.Source,
				"sourceOfTruth": event.SourceOfTruth,
				"appSource":     event.AppSource,
				"createdAt":     event.CreatedAt,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}

func (r *userRepository) UpdateUser(ctx context.Context, aggregateId string, event events.UserUpdatedEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(p:User:User_%s {id:$id})
		 SET	p.name = $name,
				p.firstName = $firstName,
				p.lastName = $lastName,
				p.sourceOfTruth = $sourceOfTruth,
				p.updatedAt = $updatedAt,
				p.syncedWithEventStore = true`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, event.Tenant),
			map[string]any{
				"id":            aggregateId,
				"tenant":        event.Tenant,
				"name":          event.Name,
				"firstName":     event.FirstName,
				"lastName":      event.LastName,
				"sourceOfTruth": event.SourceOfTruth,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	return err
}
