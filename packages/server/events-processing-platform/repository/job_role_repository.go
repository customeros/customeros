package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/events"
	"time"
)

type JobRoleRepository interface {
	LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string, updatedAt time.Time) error
	CreateJobRole(ctx context.Context, tenant, jobRoleId string, event events.JobRoleCreateEvent) error
}

type jobRoleRepository struct {
	driver *neo4j.DriverWithContext
}

func NewJobRoleRepository(driver *neo4j.DriverWithContext) JobRoleRepository {
	return &jobRoleRepository{
		driver: driver,
	}
}

func (r *jobRoleRepository) LinkWithUser(ctx context.Context, tenant, userId, jobRoleId string, updatedAt time.Time) error {
	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (u:User_%s {id: $userId})
              MERGE (jr:JobRole:JobRole_%s {id: $jobRoleId})
              ON CREATE SET jr.syncedWithEventStore = true
              MERGE (u)-[r:WORKS_AS]->(jr)
			  SET u.updatedAt = $updatedAt`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if _, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]any{
				"userId":    userId,
				"jobRoleId": jobRoleId,
				"updatedAt": updatedAt,
			}); err != nil {
			return nil, err
		} else {
			return nil, nil
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *jobRoleRepository) CreateJobRole(ctx context.Context, tenant, jobRoleId string, event events.JobRoleCreateEvent) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `  MATCH (t:Tenant {name:$tenant})
				MERGE (jr:JobRole:JobRole_%s {id:$id}) 
				SET 	jr.jobTitle = $jobTitle,
						jr.description = $description,
						jr.createdAt = $createdAt,
						jr.updatedAt = $updatedAt,
						jr.startedAt = $startedAt,
 						jr.endedAt = $endedAt,
						jr.sourceOfTruth = $sourceOfTruth,
						jr.source = $source,
						jr.appSource = $appSource,
						jr.syncedWithEventStore = true
`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"id":            jobRoleId,
				"jobTitle":      event.JobTitle,
				"description":   event.Description,
				"tenant":        event.Tenant,
				"startedAt":     event.StartedAt,
				"endedAt":       event.EndedAt,
				"sourceOfTruth": event.SourceOfTruth,
				"source":        event.Source,
				"appSource":     event.AppSource,
				"createdAt":     event.CreatedAt,
				"updatedAt":     event.UpdatedAt,
			})
		return nil, err
	})
	if err != nil {
		return err
	}
	return nil
}
