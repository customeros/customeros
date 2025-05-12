package neo4j_repository

import (
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"golang.org/x/net/context"
)

type UserUpdateFields struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Source          string `json:"source"`
	ProfilePhotoUrl string `json:"profilePhotoUrl"`
	Timezone        string `json:"timezone"`
}

type UserWriteRepository interface {
	CreateUserInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, userId string, data data_fields.UserFields) error
	UpdateUserInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, userId string, data data_fields.UserFields) error

	AddRole(ctx context.Context, userId, role string) error
	AddRoleInTx(ctx context.Context, tx neo4j.ManagedTransaction, userId, role string) error
	RemoveRole(ctx context.Context, tenant, userId, role string) error
	RegisterLogin(ctx context.Context, tenant, userId string) error
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

func (r *userWriteRepository) CreateUserInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, userId string, data data_fields.UserFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserWriteRepository.CreateUserInTx")
	defer spans.Finish()

	spans.TagEntity(userId)
	spans.LogObjectAsJson("data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:USER_BELONGS_TO_TENANT]-(u:User:User_%s {id:$id}) 
		 ON CREATE SET 	u.firstName = $firstName,
						u.lastName = $lastName,
						u.source = $source,
						u.appSource = $appSource,
						u.createdAt = $createdAt,
						u.updatedAt = datetime(),
						u.internal = $internal,
						u.test = $test,
						u.roles = $roles,
						u.bot = $bot,
						u.profilePhotoUrl = $profilePhotoUrl,
						u.timezone = $timezone`, tenant)
	roles := []string{}
	if data.Roles != nil {
		roles = *data.Roles
	}
	params := map[string]any{
		"tenant":          tenant,
		"id":              userId,
		"createdAt":       utils.IfNotNilTimeWithDefault(data.CreatedAt, utils.Now()),
		"firstName":       utils.IfNotNilString(data.FirstName),
		"lastName":        utils.IfNotNilString(data.LastName),
		"internal":        utils.IfNotNilBool(data.Internal),
		"test":            utils.IfNotNilBool(data.Test),
		"bot":             utils.IfNotNilBool(data.Bot),
		"roles":           roles,
		"profilePhotoUrl": utils.IfNotNilString(data.ProfilePhotoUrl),
		"profilePhotoKey": utils.IfNotNilString(data.ProfilePhotoKey),
		"timezone":        utils.IfNotNilString(data.Timezone),
		"source":          utils.IfNotNilString(data.Source),
		"appSource":       utils.IfNotNilString(data.AppSource),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *userWriteRepository) UpdateUserInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, userId string, data data_fields.UserFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserWriteRepository.UpdateUser")
	defer spans.Finish()

	spans.TagEntity(userId)
	spans.LogObjectAsJson("data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$id}) SET u.updatedAt=datetime()`
	params := map[string]any{
		"id":     userId,
		"tenant": tenant,
	}

	if data.FirstName != nil {
		params["firstName"] = *data.FirstName
		cypher += ", u.firstName=$firstName"
	}
	if data.LastName != nil {
		params["lastName"] = *data.LastName
		cypher += ", u.lastName=$lastName"
	}
	if data.Timezone != nil {
		params["timezone"] = *data.Timezone
		cypher += ", u.timezone=$timezone"
	}
	if data.ProfilePhotoUrl != nil {
		params["profilePhotoUrl"] = *data.ProfilePhotoUrl
		cypher += ", u.profilePhotoUrl=$profilePhotoUrl"
	}
	if data.ProfilePhotoKey != nil {
		params["profilePhotoKey"] = *data.ProfilePhotoKey
		cypher += ", u.profilePhotoKey=$profilePhotoKey"
	}
	if data.Bot != nil {
		params["bot"] = *data.Bot
		cypher += ", u.bot=$bot"
	}
	if data.Internal != nil {
		params["internal"] = *data.Internal
		cypher += ", u.internal=$internal"
	}
	if data.Test != nil {
		params["test"] = *data.Test
		cypher += ", u.test=$test"
	}
	if data.ShowOnboardingPage != nil {
		params["showOnboardingPage"] = *data.ShowOnboardingPage
		cypher += ", u.showOnboardingPage=$showOnboardingPage"
	}
	if data.OnboardingInboundStepCompleted != nil {
		params["onboardingInboundStepCompleted"] = *data.OnboardingInboundStepCompleted
		cypher += ", u.onboardingInboundStepCompleted=$onboardingInboundStepCompleted"
	}
	if data.OnboardingOutboundStepCompleted != nil {
		params["onboardingOutboundStepCompleted"] = *data.OnboardingOutboundStepCompleted
		cypher += ", u.onboardingOutboundStepCompleted=$onboardingOutboundStepCompleted"
	}
	if data.OnboardingCrmStepCompleted != nil {
		params["onboardingCrmStepCompleted"] = *data.OnboardingCrmStepCompleted
		cypher += ", u.onboardingCrmStepCompleted=$onboardingCrmStepCompleted"
	}
	if data.OnboardingMailstackStepCompleted != nil {
		params["onboardingMailstackStepCompleted"] = *data.OnboardingMailstackStepCompleted
		cypher += ", u.onboardingMailstackStepCompleted=$onboardingMailstackStepCompleted"
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *userWriteRepository) AddRole(c context.Context, userId, role string) error {
	spans, ctx := telemetry.StartNeo4jSpan(c, "UserWriteRepository.AddRole")
	defer spans.Finish()

	spans.TagEntity(userId)
	spans.LogKV("role", role)

	session := r.prepareWriteSession(ctx)
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	defer tx.Close(ctx)

	err = r.AddRoleInTx(ctx, tx, userId, role)
	if err != nil {
		spans.TraceError(err)
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		spans.TraceError(err)
		return err
	}

	return err
}

func (r *userWriteRepository) AddRoleInTx(c context.Context, tx neo4j.ManagedTransaction, userId, role string) error {
	spans, ctx := telemetry.StartNeo4jSpan(c, "UserWriteRepository.AddRoleInTx")
	defer spans.Finish()

	spans.TagEntity(userId)
	spans.LogKV("role", role)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	if err := utils.ExecuteQueryInTx(ctx, tx, cypher, params); err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (r *userWriteRepository) RemoveRole(c context.Context, tenant, userId, role string) error {
	spans, ctx := telemetry.StartNeo4jSpan(c, "UserWriteRepository.RemoveRole")
	defer spans.Finish()

	spans.TagEntity(userId)
	spans.LogKV("role", role)

	cypher := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) 
		 	SET u.roles = [item IN u.roles WHERE item <> $role],
				u.updatedAt=datetime()`
	params := map[string]any{
		"tenant": tenant,
		"role":   role,
		"userId": userId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *userWriteRepository) RegisterLogin(ctx context.Context, tenant, userId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "UserWriteRepository.RegisterLogin")
	defer spans.Finish()

	spans.TagEntity(userId)

	cypher := `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			SET u.lastLogin = $now,
				u.firstLogin = CASE WHEN u.firstLogin IS NULL THEN $now ELSE u.firstLogin END`
	params := map[string]any{
		"userId": userId,
		"tenant": tenant,
		"now":    utils.Now(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}
