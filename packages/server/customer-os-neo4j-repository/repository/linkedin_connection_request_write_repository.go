package neo4j_repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type LinkedinConnectionRequestWriteRepository interface {
	Save(ctx context.Context, tx *neo4j.ManagedTransaction, input *neo4j_entity.LinkedinConnectionRequest) error
}

type linkedinConnectionRequestWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewLinkedinConnectionRequestWriteRepository(driver *neo4j.DriverWithContext, database string) LinkedinConnectionRequestWriteRepository {
	return &linkedinConnectionRequestWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *linkedinConnectionRequestWriteRepository) Save(ctx context.Context, tx *neo4j.ManagedTransaction, input *neo4j_entity.LinkedinConnectionRequest) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "LinkedinConnectionRequestWriteRepository.Save")
	defer spans.Finish()

	spans.LogObjectAsJson("input", input)

	tenant := common.GetTenantFromContext(ctx)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)<-[:BELONGS_TO_TENANT]-(l:LinkedinConnectionRequest:LinkedinConnectionRequest_%s {id:$id})
							ON CREATE SET 
								l.createdAt=$createdAt,
								l.updatedAt=$createdAt,
								l.producerId=$producerId,
								l.producerType=$producerType,
								l.socialUrl=$socialUrl,
								l.userId=$userId,
								l.scheduledAt=$scheduledAt,
								l.status=$status
							ON MATCH SET
								l.status=$status
							`, tenant)
		params := map[string]any{
			"tenant":       tenant,
			"id":           input.Id,
			"createdAt":    utils.Now(),
			"producerId":   input.ProducerId,
			"producerType": input.ProducerType,
			"socialUrl":    input.SocialUrl,
			"userId":       input.UserId,
			"scheduledAt":  input.ScheduledAt,
			"status":       input.Status,
		}
		spans.LogKV("cypher", cypher)
		spans.LogObjectAsJson("params", params)

		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

		return nil, nil
	})

	return err
}
