package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type EmailCreateFields struct {
	RawEmail  string            `json:"rawEmail"`
	Source    entity.DataSource `json:"source"`
	CreatedAt time.Time         `json:"createdAt"`
}

type EmailValidatedFields struct {
	EmailAddress    string    `json:"emailAddress"`
	Domain          string    `json:"domain"`
	IsCatchAll      bool      `json:"isCatchAll"`
	Deliverable     string    `json:"deliverable"`
	IsValidSyntax   bool      `json:"isValidSyntax"`
	Username        string    `json:"username"`
	ValidatedAt     time.Time `json:"validatedAt"`
	IsRoleAccount   bool      `json:"isRoleAccount"`
	IsRisky         bool      `json:"isRisky"`
	IsFirewalled    bool      `json:"isFirewalled"`
	Provider        string    `json:"provider"`
	Firewall        string    `json:"firewall"`
	IsMailboxFull   bool      `json:"isMailboxFull"`
	IsFreeAccount   bool      `json:"isFreeAccount"`
	SmtpSuccess     bool      `json:"smtpSuccess"`
	ResponseCode    string    `json:"responseCode"`
	ErrorCode       string    `json:"errorCode"`
	Description     string    `json:"description"`
	IsPrimaryDomain bool      `json:"isPrimaryDomain"`
	PrimaryDomain   string    `json:"primaryDomain"`
	AlternateEmail  string    `json:"alternateEmail"`
	RetryValidation bool      `json:"retryValidation"`
}

type EmailWriteRepository interface {
	CreateEmail(ctx context.Context, tenant, emailId string, data EmailCreateFields) error
	EmailValidated(ctx context.Context, tenant, emailId string, data EmailValidatedFields) error
	CleanEmailValidation(ctx context.Context, tenant, emailId string) error
	LinkWithContact(ctx context.Context, tenant, contactId, emailId string, primary bool) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId string, primary bool) error
	LinkWithUser(ctx context.Context, tenant, userId, emailId string, primary bool) error
	UnlinkFromUser(ctx context.Context, tenant, usedId, email string) error
	UnlinkFromContact(ctx context.Context, tenant, contactId, email string) error
	UnlinkFromOrganization(ctx context.Context, tenant, organizationId, email string) error
	SetPrimaryForEntity(ctx context.Context, tenant, entityId, email string, entityType model.EntityType) error
	DeleteEmail(ctx context.Context, tenant, emailId string) error
}

type emailWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewEmailWriteRepository(driver *neo4j.DriverWithContext, database string) EmailWriteRepository {
	return &emailWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *emailWriteRepository) CreateEmail(ctx context.Context, tenant, emailId string, data EmailCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.CreateEmail")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
              MERGE (e:Email:Email_%s {id:$id})
				 SET e.rawEmail = $rawEmail, 
					e.source = $source,
					e.createdAt = $createdAt,
					e.updatedAt = datetime() 
		 MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e)`, tenant)
	params := map[string]any{
		"id":        emailId,
		"rawEmail":  data.RawEmail,
		"tenant":    tenant,
		"source":    utils.StringFirstNonEmpty(data.Source.String(), constants.SourceOpenline),
		"createdAt": utils.NowIfZero(data.CreatedAt),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) EmailValidated(ctx context.Context, tenant, emailId string, data EmailValidatedFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.EmailValidated")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id})
		 		SET e.email = CASE WHEN $email <> '' THEN $email ELSE e.email END,
					e.isCatchAll = $isCatchAll,
					e.deliverable = $deliverable,
					e.isValidSyntax = $isValidSyntax,
					e.username = $username,
					e.updatedAt = datetime(),
					e.isRoleAccount = $isRoleAccount,
					e.techValidatedAt = $validatedAt,
					e.isRisky = $isRisky,
					e.isFirewalled = $isFirewalled,
					e.provider = $provider,
					e.firewall = $firewall,
					e.isMailboxFull = $isMailboxFull,
					e.isFreeAccount = $isFreeAccount,
					e.smtpSuccess = $smtpSuccess,
					e.verifyResponseCode = $verifyResponseCode,
					e.verifyErrorCode = $verifyErrorCode,
					e.verifyDescription = $verifyDescription,
					e.isPrimaryDomain = $isPrimaryDomain,
					e.primaryDomain = $primaryDomain,
					e.alternateEmail = $alternateEmail,
					e.retryValidation = $retryValidation,
					e.work = CASE WHEN e.work IS NULL THEN NOT $isFreeAccount ELSE e.work END
				WITH e, CASE WHEN $domain <> '' THEN true ELSE false END AS shouldMergeDomain
				WHERE shouldMergeDomain
				MERGE (d:Domain {domain:$domain})
				ON CREATE SET 	d.id=randomUUID(), 
								d.createdAt=$now, 
								d.updatedAt=datetime(),
								d.source=$source
				WITH d, e
				MERGE (e)-[:HAS_DOMAIN]->(d)`, tenant)
	params := map[string]any{
		"id":                 emailId,
		"tenant":             tenant,
		"email":              data.EmailAddress,
		"domain":             strings.ToLower(data.Domain),
		"isCatchAll":         data.IsCatchAll,
		"deliverable":        data.Deliverable,
		"isValidSyntax":      data.IsValidSyntax,
		"username":           data.Username,
		"validatedAt":        data.ValidatedAt,
		"isRoleAccount":      data.IsRoleAccount,
		"isRisky":            data.IsRisky,
		"isFirewalled":       data.IsFirewalled,
		"provider":           data.Provider,
		"firewall":           data.Firewall,
		"isMailboxFull":      data.IsMailboxFull,
		"isFreeAccount":      data.IsFreeAccount,
		"smtpSuccess":        data.SmtpSuccess,
		"verifyResponseCode": data.ResponseCode,
		"verifyErrorCode":    data.ErrorCode,
		"verifyDescription":  data.Description,
		"isPrimaryDomain":    data.IsPrimaryDomain,
		"primaryDomain":      data.PrimaryDomain,
		"alternateEmail":     data.AlternateEmail,
		"retryValidation":    data.RetryValidation,
		"now":                utils.Now(),
		"source":             constants.SourceOpenline,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) LinkWithContact(ctx context.Context, tenant, contactId, emailId string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.LinkWithContact")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (c)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			c.updatedAt = datetime()`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"emailId":   emailId,
		"primary":   primary,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (org)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			org.updatedAt = datetime()`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"emailId":        emailId,
		"primary":        primary,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) LinkWithUser(ctx context.Context, tenant, userId, emailId string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.LinkWithUser")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (u)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			u.updatedAt = datetime()`
	params := map[string]any{
		"tenant":  tenant,
		"userId":  userId,
		"emailId": emailId,
		"primary": primary,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) CleanEmailValidation(ctx context.Context, tenant, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.CleanEmailValidation")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$id})
				WHERE e:Email_%s
		 		SET e.email = "",
					e.isCatchAll = null,
					e.deliverable = null,
					e.isValidSyntax = null,
					e.username = null,
					e.isRoleAccount = null,
					e.techValidatedAt = null,
					e.isRisky = null,
					e.isFirewalled = null,
					e.provider = null,
					e.firewall = null,
					e.isMailboxFull = null,
					e.isFreeAccount = null,
					e.smtpSuccess = null,
					e.verifyResponseCode = null,
					e.verifyErrorCode = null,
					e.verifyDescription = null,
					e.isPrimaryDomain = null,
					e.primaryDomain = null,
					e.alternateEmail = null,
					e.retryValidation = null,
					e.work = null,
					e.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":     emailId,
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) UnlinkFromUser(ctx context.Context, tenant, usedId, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.UnlinkFromUser")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, usedId)
	span.LogKV("email", email)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})-[rel:HAS]->(e:Email)
				WHERE e.email = $email OR e.rawEmail = $email
				DELETE rel`
	params := map[string]any{
		"tenant": tenant,
		"userId": usedId,
		"email":  email,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) UnlinkFromContact(ctx context.Context, tenant, contactId, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.UnlinkFromContact")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, contactId)
	span.LogKV("email", email)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId})-[rel:HAS]->(e:Email)
				WHERE e.email = $email OR e.rawEmail = $email
				DELETE rel`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"email":     email,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) UnlinkFromOrganization(ctx context.Context, tenant, organizationId, email string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.UnlinkFromOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, organizationId)
	span.LogKV("email", email)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId})-[rel:HAS]->(e:Email)
				WHERE e.email = $email OR e.rawEmail = $email
				SET o.updatedAt = datetime()
				DELETE rel`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"email":          email,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) DeleteEmail(ctx context.Context, tenant, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.DeleteEmail")
	defer span.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[r:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$id})
				OPTIONAL MATCH (e)-[r2]-(d:Domain)
				DELETE r, r2, e`
	params := map[string]any{
		"tenant": tenant,
		"id":     emailId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) SetPrimaryForEntity(ctx context.Context, tenant, entityId, email string, entityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.SetPrimaryForEntity")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.LogKV("entityId", entityId, "email", email, "entityType", entityType.String())

	cypher := fmt.Sprintf(`MATCH (entity:%s {id:$entityId})-[rel:HAS]->(e:Email)
				WHERE e.email = $email OR e.rawEmail = $email
				SET rel.primary = true,
					entity.updatedAt = datetime()
				WITH entity
				MATCH (entity)-[r:HAS]->(oe:Email)
				WHERE oe.email <> $email AND oe.rawEmail <> $email
				SET r.primary = false`, entityType.Neo4jLabel()+"_"+tenant)

	params := map[string]any{
		"entityId": entityId,
		"email":    email,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
