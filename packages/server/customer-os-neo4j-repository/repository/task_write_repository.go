package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"
)

type TaskWriteRepository interface {
	Create(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, taskId string, data data_fields.TaskFields) error
	Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, taskId string, data data_fields.TaskFields) error
}

type taskWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTaskWriteRepository(driver *neo4j.DriverWithContext, database string) TaskWriteRepository {
	return &taskWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *taskWriteRepository) Create(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, taskId string, data data_fields.TaskFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskWriteRepository.Create")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, taskId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task {id:$taskId}) 
							ON CREATE SET 
								tsk:Task_%s,
								tsk.createdAt=datetime(),
								tsk.updatedAt=datetime(),
								tsk.source=$source,
								tsk.appSource=$appSource,
								tsk.subject=$subject,
								tsk.description=$description,	
								tsk.status=$status,	
								tsk.dueAt=$dueAt
							WITH tsk, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$createdByUserId}) 
							WHERE $createdByUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (tsk)-[:CREATED_BY]->(u))
							`, tenant)
	params := map[string]any{
		"tenant":          tenant,
		"taskId":          taskId,
		"source":          utils.IfNotNilString(data.Source),
		"appSource":       utils.IfNotNilString(data.AppSource),
		"subject":         utils.IfNotNilString(data.Subject),
		"description":     utils.IfNotNilString(data.Description),
		"status":          utils.IfNotNilString(data.Status),
		"createdByUserId": utils.IfNotNilString(data.CreatedByUserId),
		"dueAt":           utils.TimePtrAsAny(data.DueAt),
	}

	return LogAndExecuteWriteQueryInTx(ctx, tx, r.driver, r.database, cypher, params, span)
}

func (r *taskWriteRepository) Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, taskId string, data data_fields.TaskFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TaskWriteRepository.Update")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, taskId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:TASK_BELONGS_TO_TENANT]-(tsk:Task {id:$taskId})
		 	SET	tsk.updatedAt = datetime() `
	params := map[string]any{
		"tenant": tenant,
		"taskId": taskId,
	}

	if data.Subject != nil {
		cypher += `, tsk.subject = $subject`
		params["subject"] = *data.Subject
	}
	if data.Description != nil {
		cypher += `, tsk.description = $description`
		params["description"] = *data.Description
	}
	if data.Status != nil {
		cypher += `, i.status = $status`
		params["status"] = *data.Status
	}
	if data.DueAt != nil {
		cypher += `, i.dueAt = $dueAt`
		params["dueAt"] = utils.TimePtrAsAny(data.DueAt)
	}

	return LogAndExecuteWriteQueryInTx(ctx, tx, r.driver, r.database, cypher, params, span)
}
