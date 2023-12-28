package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/constant"
	"github.com/openline-ai/customer-os-neo4j-repository/model"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"time"
)

type UserCreateData struct {
	Name            string       `json:"name"`
	FirstName       string       `json:"firstName"`
	LastName        string       `json:"lastName"`
	SourceFields    model.Source `json:"sourceFields"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	Internal        bool         `json:"internal"`
	Bot             bool         `json:"bot"`
	ProfilePhotoUrl string       `json:"profilePhotoUrl"`
	Timezone        string       `json:"timezone"`
}

type UserWriteRepository interface {
	CreateUser(ctx context.Context, tenant, userId string, data UserCreateData) error
	CreateUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, userId string, data UserCreateData) error
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

func (r *userWriteRepository) CreateUser(ctx context.Context, tenant, userId string, data UserCreateData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserWriteRepository.CreateUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, userId)
	tracing.LogObjectAsJson(span, "data", data)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return nil, r.CreateUserInTx(ctx, tx, tenant, userId, data)
	})
	return err
}

func (r *userWriteRepository) CreateUserInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, userId string, data UserCreateData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "UserWriteRepository.CreateUserInTx")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, userId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
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
						u.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"tenant":          tenant,
		"id":              userId,
		"name":            data.Name,
		"firstName":       data.FirstName,
		"lastName":        data.LastName,
		"internal":        data.Internal,
		"bot":             data.Bot,
		"profilePhotoUrl": data.ProfilePhotoUrl,
		"timezone":        data.Timezone,
		"source":          data.SourceFields.Source,
		"sourceOfTruth":   data.SourceFields.SourceOfTruth,
		"appSource":       data.SourceFields.AppSource,
		"createdAt":       data.CreatedAt,
		"updatedAt":       data.UpdatedAt,
		"overwrite":       data.SourceFields.SourceOfTruth == constant.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if err := utils.ExecuteQueryInTx(ctx, tx, cypher, params); err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}
