package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type UserUpdateFields struct {
	Name            string `json:"name"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Source          string `json:"source"`
	ProfilePhotoUrl string `json:"profilePhotoUrl"`
	Timezone        string `json:"timezone"`
}

type UserWriteRepository interface {
	CreateUser(ctx context.Context, input neo4jentity.UserEntity) error
	CreateUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, input neo4jentity.UserEntity) error

	UpdateUser(ctx context.Context, tenant, userId string, data UserUpdateFields) error

	AddRole(ctx context.Context, userId, role string) error
	AddRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, userId, role string) error
	RemoveRole(ctx context.Context, tenant, userId, role string) error
}

type userWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewUserWriteRepository(driver *neo4j.DriverWithContext, database string) UserWriteRepository {
	return &userWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *userWriteRepository) prepareWriteSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *userWriteRepository) CreateUser(c context.Context, input neo4jentity.UserEntity) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserWriteRepository.CreateUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, input.Id)

	tracing.LogObjectAsJson(span, "input", input)

	tenant := common.GetTenantFromContext(ctx)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateUserInTx(ctx, tx, tenant, input)
	})
	return err
}

func (r *userWriteRepository) CreateUserInTx(c context.Context, tx neo4j.ManagedTransaction, tenant string, input neo4jentity.UserEntity) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserWriteRepository.CreateUserInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, input.Id)

	tracing.LogObjectAsJson(span, "input", input)

	if input.Source == "" {
		input.Source = constants.SourceOpenline
	}
	if input.SourceOfTruth == "" {
		input.SourceOfTruth = constants.SourceOpenline
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User:User_%s {id:$id}) 
		 ON CREATE SET 	u.name = $name,
						u.firstName = $firstName,
						u.lastName = $lastName,
						u.source = $source,
						u.sourceOfTruth = $sourceOfTruth,
						u.appSource = $appSource,
						u.createdAt = $createdAt,
						u.updatedAt = datetime(),
						u.internal = $internal,
						u.test = $test,
						u.roles = $roles,
						u.bot = $bot,
						u.profilePhotoUrl = $profilePhotoUrl,
						u.timezone = $timezone
		 ON MATCH SET 	u.name = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.name is null OR u.name = '' THEN $name ELSE u.name END,
						u.firstName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.firstName is null OR u.firstName = '' THEN $firstName ELSE u.firstName END,
						u.lastName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.lastName is null OR u.lastName = '' THEN $lastName ELSE u.lastName END,
						u.timezone = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.timezone is null OR u.timezone = '' THEN $timezone ELSE u.timezone END,
						u.roles = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.roles is null THEN $roles ELSE u.roles END,
						u.profilePhotoUrl = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.profilePhotoUrl is null OR u.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE u.profilePhotoUrl END,
						u.internal = $internal,
						u.bot = $bot,
						u.updatedAt = datetime(),
						u.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE u.sourceOfTruth END`, tenant)
	params := map[string]any{
		"tenant":          tenant,
		"id":              input.Id,
		"name":            input.Name,
		"firstName":       input.FirstName,
		"lastName":        input.LastName,
		"internal":        input.Internal,
		"test":            input.Test,
		"roles":           input.Roles,
		"bot":             input.Bot,
		"profilePhotoUrl": input.ProfilePhotoUrl,
		"timezone":        input.Timezone,
		"source":          input.Source,
		"sourceOfTruth":   input.SourceOfTruth,
		"appSource":       input.AppSource,
		"createdAt":       utils.NowIfZero(input.CreatedAt),
		"overwrite":       input.SourceOfTruth == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if err := utils.ExecuteQueryInTx(ctx, tx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (r *userWriteRepository) UpdateUser(c context.Context, tenant, userId string, data UserUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserWriteRepository.UpdateUser")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, userId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User:User_%s {id:$id})
		 SET	u.name = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.name is null OR u.name = '' THEN $name ELSE u.name END,
				u.firstName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.firstName is null OR u.firstName = '' THEN $firstName ELSE u.firstName END,
				u.lastName = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.lastName is null OR u.lastName = '' THEN $lastName ELSE u.lastName END,
				u.timezone = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.timezone is null OR u.timezone = '' THEN $timezone ELSE u.timezone END,
				u.profilePhotoUrl = CASE WHEN u.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR u.profilePhotoUrl is null OR u.profilePhotoUrl = '' THEN $profilePhotoUrl ELSE u.profilePhotoUrl END,
				u.updatedAt = datetime(),
				u.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE u.sourceOfTruth END`, tenant)
	params := map[string]any{
		"id":              userId,
		"tenant":          tenant,
		"name":            data.Name,
		"firstName":       data.FirstName,
		"lastName":        data.LastName,
		"sourceOfTruth":   data.Source,
		"profilePhotoUrl": data.ProfilePhotoUrl,
		"timezone":        data.Timezone,
		"overwrite":       data.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *userWriteRepository) AddRole(c context.Context, userId, role string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserWriteRepository.AddRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, userId)
	span.LogFields(log.String("role", role))

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	defer tx.Close(ctx)

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = r.AddRoleInTx(ctx, tx, userId, role)
	if err != nil {
		tracing.TraceErr(span, err)
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return err
}

func (r *userWriteRepository) AddRoleInTx(c context.Context, tx neo4j.ManagedTransaction, userId, role string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserWriteRepository.AddRoleInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, userId)
	span.LogFields(log.String("role", role))

	tenant := common.GetTenantFromContext(ctx)

	cypher := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
		 	SET u.roles = CASE
					WHEN u.roles IS NULL THEN [$role]
					ELSE CASE
		 				WHEN NOT $role IN u.roles THEN u.roles + $role 
		 				ELSE u.roles 
		 				END
					END, 
				u.updatedAt=datetime()`
	params := map[string]any{
		"tenant": tenant,
		"role":   role,
		"userId": userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if err := utils.ExecuteQueryInTx(ctx, tx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (r *userWriteRepository) RemoveRole(c context.Context, tenant, userId, role string) error {
	span, ctx := opentracing.StartSpanFromContext(c, "UserWriteRepository.RemoveRole")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, userId)
	span.LogFields(log.String("role", role))

	cypher := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
		 	SET u.roles = [item IN u.roles WHERE item <> $role],
				u.updatedAt=datetime()`
	params := map[string]any{
		"tenant": tenant,
		"role":   role,
		"userId": userId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
