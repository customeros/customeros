package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type TenantWriteRepository interface {
	CreateTenantIfNotExistAndReturn(ctx context.Context, tx neo4j.ManagedTransaction, tenant neo4jentity.TenantEntity) (*dbtype.Node, error)

	CreateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error
	UpdateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error

	UpdateTenantSettings(ctx context.Context, tenant string, data data_fields.TenantSettingsFields) error

	HardDeleteTenant(ctx context.Context, tenant string) error

	MarkOnboardingChecked(ctx context.Context, tenant string) error
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

func (r *tenantWriteRepository) CreateTenantIfNotExistAndReturn(ctx context.Context, tx neo4j.ManagedTransaction, tenant neo4jentity.TenantEntity) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TenantWriteRepository.CreateTenantIfNotExistAndReturn")
	defer spans.Finish()

	spans.LogObjectAsJson("inputTenantEntity", tenant)

	cypher := `MERGE (t:Tenant {name:$name}) 
		 ON CREATE SET 
		  t.id=randomUUID(), 
		  t.createdBy=$createdBy, 
		  t.createdAt=datetime(), 
		  t.updatedAt=datetime(),
		  t.active=true
		WITH t
		MERGE (t)-[:HAS_SETTINGS]->(ts:TenantSettings {tenant:$name})
		ON CREATE SET
			ts.id=randomUUID(),
		  	ts.createdAt=datetime(),	
			ts.updatedAt=datetime(),
			ts.invoicingPostpaid=$invoicingPostpaid,
			ts.enrichContacts=$enrichContacts,
			ts.baseCurrency=$currency
		 RETURN t`
	params := map[string]any{
		"name":              tenant.Name,
		"createdBy":         tenant.CreatedBy,
		"invoicingPostpaid": false,
		"enrichContacts":    true,
		"currency":          enum.CurrencyUSD.String(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	queryResult, err := tx.Run(ctx, cypher, params)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *tenantWriteRepository) CreateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TenantWriteRepository.CreateTenantBillingProfile")
	defer spans.Finish()

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantBillingProfile(ctx context.Context, tenant string, data data_fields.TenantBillingProfileFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TenantWriteRepository.UpdateTenantBillingProfile")
	defer spans.Finish()

	spans.LogObjectAsJson("data", data)

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

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantSettings(ctx context.Context, tenant string, data data_fields.TenantSettingsFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TenantWriteRepository.UpdateTenantSettings")
	defer spans.Finish()

	spans.LogObjectAsJson("data", data)

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
	if data.InvoicingPostpaid != nil {
		cypher += ", ts.invoicingPostpaid=$invoicingPostpaid"
		params["invoicingPostpaid"] = *data.InvoicingPostpaid
	}
	if data.BaseCurrency != nil {
		cypher += ", ts.baseCurrency=$baseCurrency"
		params["baseCurrency"] = *data.BaseCurrency
	}
	if data.WorkspaceLogo != nil {
		cypher += ", ts.workspaceLogo=$workspaceLogo"
		params["workspaceLogo"] = *data.WorkspaceLogo
	}
	if data.WorkspaceName != nil {
		cypher += ", ts.workspaceName=$workspaceName"
		params["workspaceName"] = *data.WorkspaceName
	}
	if data.WorkspaceLogoKey != nil {
		cypher += ", ts.workspaceLogoKey=$workspaceLogoKey"
		params["workspaceLogoKey"] = *data.WorkspaceLogoKey
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *tenantWriteRepository) HardDeleteTenant(ctx context.Context, tenant string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TenantWriteRepository.HardDelete")
	defer spans.Finish()

	spans.LogObjectAsJson("tenant", tenant)

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
			spans.TraceError(err)
			return err
		}
	}

	//drop TenantSettings
	err := utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (t:TenantSettings{tenant: $tenant}) DETACH DELETE t`, map[string]any{"tenant": tenant})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	//drop TenantMetadata
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (t:TenantMetadata{tenantName: $tenant}) DETACH DELETE t`, map[string]any{"tenant": tenant})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	//drop External systems
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (e:ExternalSystem)-[r:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]->(t:Tenant{name: $tenant}) DELETE r, e`, map[string]any{"tenant": tenant})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	//drop workspaces
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (w:Workspace)<-[r:HAS_WORKSPACE]-(t:Tenant{name: $tenant}) DELETE r, w`, map[string]any{"tenant": tenant})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	err = utils.ExecuteWriteQuery(ctx, *r.driver,
		`match (au:AuthenticationUser)
					optional match (au)-[r:HAS_WORKSPACE]-(t:Tenant{name: $tenant})
					with au, r, t
					delete r`, map[string]any{"tenant": tenant})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	//drop tenant
	err = utils.ExecuteWriteQuery(ctx, *r.driver, `MATCH (t:Tenant{name: $tenant}) DELETE t`, map[string]any{"tenant": tenant})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return err
}

func (r *tenantWriteRepository) MarkOnboardingChecked(ctx context.Context, tenant string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "TenantWriteRepository.MarkOnboardingChecked")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})
				SET t.techOnboardingCheckedAt=datetime()
				RETURN t`
	params := map[string]any{
		"tenant": tenant,
	}
	return LogAndExecuteWriteQuery(ctx, *r.driver, cypher, params, spans)
}
