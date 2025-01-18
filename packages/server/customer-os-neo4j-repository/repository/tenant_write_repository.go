package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantWriteRepository interface {
	CreateTenantIfNotExistAndReturn(ctx context.Context, tenant neo4jentity.TenantEntity) (*dbtype.Node, error)

	CreateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error
	UpdateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error

	UpdateTenantSettings(ctx context.Context, tenant string, data data_fields.TenantSettingsFields) error

	HardDeleteTenant(ctx context.Context, tenant string) error

	LinkWithWorkspace(ctx context.Context, tenant string, workspace neo4jentity.WorkspaceEntity) (bool, error)
}

type tenantWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewTenantWriteRepository(driver *neo4j.DriverWithContext, database string) TenantWriteRepository {
	return &tenantWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *tenantWriteRepository) CreateTenantIfNotExistAndReturn(ctx context.Context, tenant neo4jentity.TenantEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.CreateTenantIfNotExistAndReturn")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant.Name)
	tracing.LogObjectAsJson(span, "inputTenantEntity", tenant)

	cypher := `MERGE (t:Tenant {name:$name}) 
		 ON CREATE SET 
		  t.id=randomUUID(), 
		  t.createdAt=datetime(), 
		  t.updatedAt=datetime(), 
		  t.source=$source, 
		  t.appSource=$appSource,
		  t.active=true
		WITH t
		MERGE (t)-[:HAS_SETTINGS]->(ts:TenantSettings {tenant:$name})
		ON CREATE SET
			ts.id=randomUUID(),
		  	ts.createdAt=datetime(),	
			ts.updatedAt=datetime(),
			ts.invoicingEnabled=$invoicingEnabled,
			ts.invoicingPostpaid=$invoicingPostpaid,
			ts.enrichContacts=$enrichContacts,
			ts.baseCurrency=$currency
		 RETURN t`
	params := map[string]any{
		"name":              tenant.Name,
		"source":            tenant.Source,
		"appSource":         tenant.AppSource,
		"invoicingEnabled":  false,
		"invoicingPostpaid": false,
		"enrichContacts":    true,
		"currency":          enum.CurrencyUSD.String(),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *tenantWriteRepository) CreateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.CreateTenantBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile {id:$billingProfileId}) 
							ON CREATE SET 
								tbp:TenantBillingProfile_%s,
								tbp.createdAt=datetime(),
								tbp.updatedAt=datetime(),
								tbp.source=$source,
								tbp.appSource=$appSource,
								tbp.phone=$phone,
								tbp.legalName=$legalName,	
								tbp.addressLine1=$addressLine1,	
								tbp.addressLine2=$addressLine2,
								tbp.addressLine3=$addressLine3,
								tbp.locality=$locality,
								tbp.country=$country,
								tbp.region=$region,
								tbp.zip=$zip,
								tbp.vatNumber=$vatNumber,	
								tbp.sendInvoicesFrom=$sendInvoicesFrom,
								tbp.sendInvoicesBcc=$sendInvoicesBcc,
								tbp.canPayWithPigeon=$canPayWithPigeon,
								tbp.canPayWithBankTransfer=$canPayWithBankTransfer,
								tbp.check=$check
							`, tenant)
	params := map[string]any{
		"tenant":                 tenant,
		"billingProfileId":       data.ID,
		"source":                 utils.IfNotNilString(data.Source),
		"appSource":              utils.IfNotNilString(data.AppSource),
		"phone":                  utils.IfNotNilString(data.Phone),
		"legalName":              utils.IfNotNilString(data.LegalName),
		"addressLine1":           utils.IfNotNilString(data.AddressLine1),
		"addressLine2":           utils.IfNotNilString(data.AddressLine2),
		"addressLine3":           utils.IfNotNilString(data.AddressLine3),
		"locality":               utils.IfNotNilString(data.Locality),
		"country":                utils.IfNotNilString(data.Country),
		"region":                 utils.IfNotNilString(data.Region),
		"zip":                    utils.IfNotNilString(data.Zip),
		"vatNumber":              utils.IfNotNilString(data.VatNumber),
		"sendInvoicesFrom":       utils.IfNotNilString(data.SendInvoicesFrom),
		"sendInvoicesBcc":        utils.IfNotNilString(data.SendInvoicesBcc),
		"canPayWithPigeon":       utils.IfNotNilBool(data.CanPayWithPigeon),
		"canPayWithBankTransfer": utils.IfNotNilBool(data.CanPayWithBankTransfer),
		"check":                  utils.IfNotNilBool(data.Check),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.UpdateTenantBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (:Tenant {name:$tenant})-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile {id:$billingProfileId}) 
							SET tbp.updatedAt=datetime()
							
							`
	params := map[string]any{
		"tenant":           tenant,
		"billingProfileId": data.ID,
	}
	if data.Phone != nil {
		cypher += `,tbp.phone=$phone`
		params["phone"] = *data.Phone
	}
	if data.LegalName != nil {
		cypher += `,tbp.legalName=$legalName`
		params["legalName"] = *data.LegalName
	}
	if data.AddressLine1 != nil {
		cypher += `,tbp.addressLine1=$addressLine1`
		params["addressLine1"] = *data.AddressLine1
	}
	if data.AddressLine2 != nil {
		cypher += `,tbp.addressLine2=$addressLine2`
		params["addressLine2"] = *data.AddressLine2
	}
	if data.AddressLine3 != nil {
		cypher += `,tbp.addressLine3=$addressLine3`
		params["addressLine3"] = *data.AddressLine3
	}
	if data.Locality != nil {
		cypher += `,tbp.locality=$locality`
		params["locality"] = *data.Locality
	}
	if data.Country != nil {
		cypher += `,tbp.country=$country`
		params["country"] = *data.Country
	}
	if data.Region != nil {
		cypher += `,tbp.region=$region`
		params["region"] = *data.Region
	}
	if data.Zip != nil {
		cypher += `,tbp.zip=$zip`
		params["zip"] = *data.Zip
	}
	if data.VatNumber != nil {
		cypher += `,tbp.vatNumber=$vatNumber`
		params["vatNumber"] = *data.VatNumber
	}
	if data.SendInvoicesFrom != nil {
		cypher += `,tbp.sendInvoicesFrom=$sendInvoicesFrom`
		params["sendInvoicesFrom"] = *data.SendInvoicesFrom
	}
	if data.SendInvoicesBcc != nil {
		cypher += `,tbp.sendInvoicesBcc=$sendInvoicesBcc`
		params["sendInvoicesBcc"] = *data.SendInvoicesBcc
	}
	if data.CanPayWithPigeon != nil {
		cypher += `,tbp.canPayWithPigeon=$canPayWithPigeon`
		params["canPayWithPigeon"] = *data.CanPayWithPigeon
	}
	if data.CanPayWithBankTransfer != nil {
		cypher += `,tbp.canPayWithBankTransfer=$canPayWithBankTransfer`
		params["canPayWithBankTransfer"] = *data.CanPayWithBankTransfer
	}
	if data.Check != nil {
		cypher += `,tbp.check=$check`
		params["check"] = *data.Check
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantSettings(ctx context.Context, tenant string, data data_fields.TenantSettingsFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.UpdateTenantSettings")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})
				MERGE (t)-[:HAS_SETTINGS]->(ts:TenantSettings {tenant:$tenant})
				ON CREATE SET
					ts.id=randomUUID(),
					ts.createdAt=datetime()
				SET
					ts.updatedAt=datetime()`
	params := map[string]any{
		"tenant": tenant,
	}
	if data.InvoicingEnabled != nil {
		cypher += ", ts.invoicingEnabled=$invoicingEnabled"
		params["invoicingEnabled"] = *data.InvoicingEnabled
	}
	if data.InvoicingPostpaid != nil {
		cypher += ", ts.invoicingPostpaid=$invoicingPostpaid"
		params["invoicingPostpaid"] = *data.InvoicingPostpaid
	}
	if data.BaseCurrency != nil {
		cypher += ", ts.baseCurrency=$baseCurrency"
		params["baseCurrency"] = *data.BaseCurrency
	}
	if data.LogoRepositoryFileId != nil {
		cypher += ", ts.logoRepositoryFileId=$logoRepositoryFileId"
		params["logoRepositoryFileId"] = *data.LogoRepositoryFileId
	}
	if data.WorkspaceLogo != nil {
		cypher += ", ts.workspaceLogo=$workspaceLogo"
		params["workspaceLogo"] = *data.WorkspaceLogo
	}
	if data.WorkspaceName != nil {
		cypher += ", ts.workspaceName=$workspaceName"
		params["workspaceName"] = *data.WorkspaceName
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tenantWriteRepository) HardDeleteTenant(ctx context.Context, tenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.HardDelete")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "tenant", tenant)

	nodeWithTenantSuffix := []string{
		commonmodel.NodeLabelTenantBillingProfile,
		commonmodel.NodeLabelBankAccount,
		commonmodel.NodeLabelTimelineEvent,
		commonmodel.NodeLabelContact,
		commonmodel.NodeLabelCustomField,
		commonmodel.NodeLabelJobRole,
		commonmodel.NodeLabelEmail,
		commonmodel.NodeLabelLocation,
		commonmodel.NodeLabelInteractionEvent,
		commonmodel.NodeLabelInteractionSession,
		commonmodel.NodeLabelNote,
		commonmodel.NodeLabelLogEntry,
		commonmodel.NodeLabelOrganization,
		commonmodel.NodeLabelBillingProfile,
		commonmodel.NodeLabelAction,
		commonmodel.NodeLabelPageView,
		commonmodel.NodeLabelPhoneNumber,
		commonmodel.NodeLabelTag,
		commonmodel.NodeLabelIssue,
		commonmodel.NodeLabelUser,
		commonmodel.NodeLabelAttachment,
		commonmodel.NodeLabelMeeting,
		commonmodel.NodeLabelSocial,
		commonmodel.NodeLabelActionItem,
		commonmodel.NodeLabelComment,
		commonmodel.NodeLabelContract,
		commonmodel.NodeLabelDeletedContract,
		commonmodel.NodeLabelServiceLineItem,
		commonmodel.NodeLabelOpportunity,
		commonmodel.NodeLabelInvoicingCycle,
		commonmodel.NodeLabelExternalSystem,
		commonmodel.NodeLabelInvoice,
		commonmodel.NodeLabelInvoiceLine,
		commonmodel.NodeLabelReminder,
		commonmodel.NodeLabelFlow,
		commonmodel.NodeLabelFlowParticipant,
		commonmodel.NodeLabelFlowSender,
		commonmodel.NodeLabelFlowAction,
		commonmodel.NodeLabelFlowActionExecution,
		commonmodel.NodeLabelFlowExecutionSettings,
	}

	//drop nodes with NodeLabel_Tenant
	for _, nodeLabel := range nodeWithTenantSuffix {
		err := utils.ExecuteWriteQuery(ctx, *r.driver, fmt.Sprintf(`MATCH (n:%s_%s) DETACH DELETE n`, nodeLabel, tenant), nil)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	//drop TenantSettings
	err := utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (t:TenantSettings{tenant: $tenant}) DETACH DELETE t`, map[string]any{"tenant": tenant})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//drop TenantMetadata
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (t:TenantMetadata{tenantName: $tenant}) DETACH DELETE t`, map[string]any{"tenant": tenant})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//drop External systems
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (e:ExternalSystem)-[r:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t:Tenant{name: $tenant}) DELETE r, e`, map[string]any{"tenant": tenant})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//drop workspaces
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (w:Workspace)<-[r:HAS_WORKSPACE]-(t:Tenant{name: $tenant}) DELETE r, w`, map[string]any{"tenant": tenant})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//drop tenant
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (t:Tenant{name: $tenant}) DELETE t`, map[string]any{"tenant": tenant})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	//clear Player nodes not linked to a user in the system
	err = utils.ExecuteWriteQuery(ctx, *r.driver,
		`match (p:Player)
					optional match (p)-[r]-(u:User)
					with p, r, u
					where u is null
					delete p`, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return err
}

func (r *tenantWriteRepository) LinkWithWorkspace(ctx context.Context, tenant string, workspace neo4jentity.WorkspaceEntity) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.LinkWithWorkspace")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	query := `
			MATCH (t:Tenant {name:$tenant})
			MATCH (w:Workspace {name:$name, provider:$provider})
			WHERE NOT ()-[:HAS_WORKSPACE]->(w)
			CREATE (t)-[:HAS_WORKSPACE]->(w)
			RETURN t`
	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":   tenant,
				"name":     workspace.Name,
				"provider": workspace.Provider,
			})
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		return false, err
	}
	convertedResult, isOk := result.([]*dbtype.Node)
	if !isOk {
		return false, errors.New("unexpected result type")
	}
	if len(convertedResult) == 0 {
		return false, nil
	}
	return true, nil
}
