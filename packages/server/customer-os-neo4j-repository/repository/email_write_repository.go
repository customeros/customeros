package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type EmailCreateFields struct {
	RawEmail     string       `json:"rawEmail"`
	SourceFields model.Source `json:"sourceFields"`
	CreatedAt    time.Time    `json:"createdAt"`
}

type EmailValidatedFields struct {
	ValidationError string    `json:"validationError"`
	EmailAddress    string    `json:"emailAddress"`
	Domain          string    `json:"domain"`
	AcceptsMail     bool      `json:"acceptsMail"`
	CanConnectSmtp  bool      `json:"canConnectSmtp"`
	HasFullInbox    bool      `json:"hasFullInbox"`
	IsCatchAll      bool      `json:"isCatchAll"`
	IsDeliverable   bool      `json:"isDeliverable"`
	IsDisabled      bool      `json:"isDisabled"`
	IsValidSyntax   bool      `json:"isValidSyntax"`
	Username        string    `json:"username"`
	ValidatedAt     time.Time `json:"validatedAt"`
	IsReachable     string    `json:"isReachable"`
	IsDisposable    bool      `json:"isDisposable"`
	IsRoleAccount   bool      `json:"isRoleAccount"`
}

type EmailWriteRepository interface {
	CreateEmail(ctx context.Context, tenant, emailId string, data EmailCreateFields) error
	UpdateEmail(ctx context.Context, tenant, emailId, rawEmail, source string) error
	FailEmailValidation(ctx context.Context, tenant, emailId, validationError string) error
	EmailValidated(ctx context.Context, tenant, emailId string, data EmailValidatedFields) error
	CleanEmailValidation(ctx context.Context, tenant, emailId string) error
	LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool) error
	LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId, label string, primary bool) error
	LinkWithUser(ctx context.Context, tenant, userId, emailId, label string, primary bool) error
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
              MERGE (e:Email:Email_%s {id:$id})
				 SET e.rawEmail = $rawEmail, 
					e.validated = null,
					e.source = $source,
					e.sourceOfTruth = $sourceOfTruth,
					e.appSource = $appSource,
					e.createdAt = $createdAt,
					e.updatedAt = datetime(),
					e.syncedWithEventStore = true 
		 MERGE (t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e)`, tenant)
	params := map[string]any{
		"id":            emailId,
		"rawEmail":      data.RawEmail,
		"tenant":        tenant,
		"source":        utils.StringFirstNonEmpty(data.SourceFields.Source, constants.SourceOpenline),
		"sourceOfTruth": utils.StringFirstNonEmpty(data.SourceFields.SourceOfTruth, constants.SourceOpenline),
		"appSource":     data.SourceFields.AppSource,
		"createdAt":     utils.TimeOrNow(data.CreatedAt),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) UpdateEmail(ctx context.Context, tenant, emailId, rawEmail, source string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.UpdateEmail")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id})
		 SET 	e.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE e.sourceOfTruth END,
				e.updatedAt = datetime(),
				e.rawEmail = $rawEmail,
				e.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"id":            emailId,
		"tenant":        tenant,
		"sourceOfTruth": source,
		"rawEmail":      rawEmail,
		"overwrite":     source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) FailEmailValidation(ctx context.Context, tenant, emailId, validationError string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.FailEmailValidation")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$id})
				WHERE e:Email_%s
		 		SET e.validationError = $validationError,
		     		e.validated = false,
					e.updatedAt = datetime()`, tenant)
	params := map[string]any{
		"id":              emailId,
		"tenant":          tenant,
		"validationError": validationError,
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email:Email_%s {id:$id})
		 		SET e.validationError = $validationError,
					e.email = $email,
		     		e.validated = true,
					e.acceptsMail = $acceptsMail,
					e.canConnectSmtp = $canConnectSmtp,
					e.hasFullInbox = $hasFullInbox,
					e.isCatchAll = $isCatchAll,
					e.isDeliverable = $isDeliverable,
					e.isDisabled = $isDisabled,
					e.isValidSyntax = $isValidSyntax,
					e.username = $username,
					e.updatedAt = datetime(),
					e.isReachable = $isReachable,
					e.isDisposable = $isDisposable,
					e.isRoleAccount = $isRoleAccount
				WITH e, CASE WHEN $domain <> '' THEN true ELSE false END AS shouldMergeDomain
				WHERE shouldMergeDomain
				MERGE (d:Domain {domain:$domain})
				ON CREATE SET 	d.id=randomUUID(), 
								d.createdAt=$now, 
								d.updatedAt=datetime(),
								d.appSource=$source,
								d.source=$appSource
				WITH d, e
				MERGE (e)-[:HAS_DOMAIN]->(d)`, tenant)
	params := map[string]any{
		"id":              emailId,
		"tenant":          tenant,
		"validationError": data.ValidationError,
		"email":           data.EmailAddress,
		"domain":          strings.ToLower(data.Domain),
		"acceptsMail":     data.AcceptsMail,
		"canConnectSmtp":  data.CanConnectSmtp,
		"hasFullInbox":    data.HasFullInbox,
		"isCatchAll":      data.IsCatchAll,
		"isDeliverable":   data.IsDeliverable,
		"isDisabled":      data.IsDisabled,
		"isValidSyntax":   data.IsValidSyntax,
		"username":        data.Username,
		"validatedAt":     data.ValidatedAt,
		"isReachable":     data.IsReachable,
		"isDisposable":    data.IsDisposable,
		"isRoleAccount":   data.IsRoleAccount,
		"now":             utils.Now(),
		"source":          constants.SourceOpenline,
		"appSource":       constants.AppSourceEventProcessingPlatform,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *emailWriteRepository) LinkWithContact(ctx context.Context, tenant, contactId, emailId, label string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.LinkWithContact")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (c)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			c.updatedAt = datetime(),
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":    tenant,
		"contactId": contactId,
		"emailId":   emailId,
		"label":     label,
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

func (r *emailWriteRepository) LinkWithOrganization(ctx context.Context, tenant, organizationId, emailId, label string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.LinkWithOrganization")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$organizationId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (org)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			org.updatedAt = datetime(),
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"emailId":        emailId,
		"label":          label,
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

func (r *emailWriteRepository) LinkWithUser(ctx context.Context, tenant, userId, emailId, label string, primary bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailWriteRepository.LinkWithUser")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := `
		MATCH (t:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(t)<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$emailId})
		MERGE (u)-[rel:HAS]->(e)
		SET	rel.primary = $primary,
			rel.label = $label,	
			u.updatedAt = datetime(),
			rel.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":  tenant,
		"userId":  userId,
		"emailId": emailId,
		"label":   label,
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
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, emailId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EMAIL_ADDRESS_BELONGS_TO_TENANT]-(e:Email {id:$id})
				WHERE e:Email_%s
		 		SET e.validationError = null,
		     		e.validated = null,
					e.email = "",
					e.acceptsMail = null,
					e.canConnectSmtp = null,
					e.hasFullInbox = null,
					e.isCatchAll = null,
					e.isDeliverable = null,
					e.isDisabled = null,
					e.isValidSyntax = null,
					e.username = null,
					e.isReachable = null,
					e.isDisposable = null,
					e.isRoleAccount = null,
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
