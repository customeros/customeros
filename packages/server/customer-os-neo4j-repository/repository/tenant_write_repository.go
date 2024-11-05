package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type TenantBillingProfileCreateFields struct {
	Id                     string             `json:"id"`
	CreatedAt              time.Time          `json:"createdAt"`
	SourceFields           model.SourceFields `json:"sourceFields"`
	Phone                  string             `json:"phone"`
	LegalName              string             `json:"legalName"`
	AddressLine1           string             `json:"addressLine1"`
	AddressLine2           string             `json:"addressLine2"`
	AddressLine3           string             `json:"addressLine3"`
	Locality               string             `json:"locality"`
	Country                string             `json:"country"`
	Region                 string             `json:"region"`
	Zip                    string             `json:"zip"`
	VatNumber              string             `json:"vatNumber"`
	SendInvoicesFrom       string             `json:"sendInvoicesFrom"`
	SendInvoicesBcc        string             `json:"sendInvoicesBcc"`
	CanPayWithPigeon       bool               `json:"canPayWithPigeon"`
	CanPayWithBankTransfer bool               `json:"canPayWithBankTransfer"`
	Check                  bool               `json:"check"`
}

type TenantBillingProfileUpdateFields struct {
	Id                           string `json:"id"`
	Phone                        string `json:"phone"`
	LegalName                    string `json:"legalName"`
	AddressLine1                 string `json:"addressLine1"`
	AddressLine2                 string `json:"addressLine2"`
	AddressLine3                 string `json:"addressLine3"`
	Locality                     string `json:"locality"`
	Country                      string `json:"country"`
	Region                       string `json:"region"`
	Zip                          string `json:"zip"`
	VatNumber                    string `json:"vatNumber"`
	SendInvoicesFrom             string `json:"sendInvoicesFrom"`
	SendInvoicesBcc              string `json:"sendInvoicesBcc"`
	CanPayWithPigeon             bool   `json:"canPayWithPigeon"`
	CanPayWithBankTransfer       bool   `json:"canPayWithBankTransfer"`
	Check                        bool   `json:"check"`
	UpdatePhone                  bool   `json:"updatePhone"`
	UpdateLegalName              bool   `json:"updateLegalName"`
	UpdateAddressLine1           bool   `json:"updateAddressLine1"`
	UpdateAddressLine2           bool   `json:"updateAddressLine2"`
	UpdateAddressLine3           bool   `json:"updateAddressLine3"`
	UpdateLocality               bool   `json:"updateLocality"`
	UpdateCountry                bool   `json:"updateCountry"`
	UpdateRegion                 bool   `json:"updateRegion"`
	UpdateZip                    bool   `json:"updateZip"`
	UpdateVatNumber              bool   `json:"updateVatNumber"`
	UpdateSendInvoicesFrom       bool   `json:"updateSendInvoicesFrom"`
	UpdateSendInvoicesBcc        bool   `json:"updateSendInvoicesBcc"`
	UpdateCanPayWithPigeon       bool   `json:"updateCanPayWithPigeon"`
	UpdateCanPayWithBankTransfer bool   `json:"updateCanPayWithBankTransfer"`
	UpdateCheck                  bool   `json:"updateCheck"`
}

type TenantSettingsFields struct {
	LogoRepositoryFileId       string        `json:"logoRepositoryFileId"`
	BaseCurrency               enum.Currency `json:"baseCurrency"`
	InvoicingEnabled           bool          `json:"invoicingEnabled"`
	InvoicingPostpaid          bool          `json:"invoicingPostpaid"`
	WorkspaceLogo              string        `json:"workspaceLogo"`
	WorkspaceName              string        `json:"workspaceName"`
	UpdateLogoRepositoryFileId bool          `json:"updateLogoRepositoryFileId"`
	UpdateInvoicingEnabled     bool          `json:"updateInvoicingEnabled"`
	UpdateInvoicingPostpaid    bool          `json:"updateInvoicingPostpaid"`
	UpdateBaseCurrency         bool          `json:"updateBaseCurrency"`
	UpdateWorkspaceLogo        bool          `json:"updateWorkspaceLogo"`
	UpdateWorkspaceName        bool          `json:"updateWorkspaceName"`
}

type TenantWriteRepository interface {
	CreateTenantIfNotExistAndReturn(ctx context.Context, tenant neo4jentity.TenantEntity) (*dbtype.Node, error)

	CreateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileCreateFields) error
	UpdateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileUpdateFields) error
	UpdateTenantSettings(ctx context.Context, tenant string, data TenantSettingsFields) error

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

func (r *tenantWriteRepository) CreateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.CreateTenantBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)-[:HAS_BILLING_PROFILE]->(tbp:TenantBillingProfile {id:$billingProfileId}) 
							ON CREATE SET 
								tbp:TenantBillingProfile_%s,
								tbp.createdAt=$createdAt,
								tbp.updatedAt=datetime(),
								tbp.source=$source,
								tbp.sourceOfTruth=$sourceOfTruth,
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
		"billingProfileId":       data.Id,
		"createdAt":              data.CreatedAt,
		"source":                 data.SourceFields.Source,
		"sourceOfTruth":          data.SourceFields.Source,
		"appSource":              data.SourceFields.AppSource,
		"phone":                  data.Phone,
		"legalName":              data.LegalName,
		"addressLine1":           data.AddressLine1,
		"addressLine2":           data.AddressLine2,
		"addressLine3":           data.AddressLine3,
		"locality":               data.Locality,
		"country":                data.Country,
		"region":                 data.Region,
		"zip":                    data.Zip,
		"vatNumber":              data.VatNumber,
		"sendInvoicesFrom":       data.SendInvoicesFrom,
		"sendInvoicesBcc":        data.SendInvoicesBcc,
		"canPayWithPigeon":       data.CanPayWithPigeon,
		"canPayWithBankTransfer": data.CanPayWithBankTransfer,
		"check":                  data.Check,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantBillingProfile(ctx context.Context, tenant string, data TenantBillingProfileUpdateFields) error {
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
		"billingProfileId": data.Id,
	}
	if data.UpdatePhone {
		cypher += `,tbp.phone=$phone`
		params["phone"] = data.Phone
	}
	if data.UpdateLegalName {
		cypher += `,tbp.legalName=$legalName`
		params["legalName"] = data.LegalName
	}
	if data.UpdateAddressLine1 {
		cypher += `,tbp.addressLine1=$addressLine1`
		params["addressLine1"] = data.AddressLine1
	}
	if data.UpdateAddressLine2 {
		cypher += `,tbp.addressLine2=$addressLine2`
		params["addressLine2"] = data.AddressLine2
	}
	if data.UpdateAddressLine3 {
		cypher += `,tbp.addressLine3=$addressLine3`
		params["addressLine3"] = data.AddressLine3
	}
	if data.UpdateLocality {
		cypher += `,tbp.locality=$locality`
		params["locality"] = data.Locality
	}
	if data.UpdateCountry {
		cypher += `,tbp.country=$country`
		params["country"] = data.Country
	}
	if data.UpdateRegion {
		cypher += `,tbp.region=$region`
		params["region"] = data.Region
	}
	if data.UpdateZip {
		cypher += `,tbp.zip=$zip`
		params["zip"] = data.Zip
	}
	if data.UpdateVatNumber {
		cypher += `,tbp.vatNumber=$vatNumber`
		params["vatNumber"] = data.VatNumber
	}
	if data.UpdateSendInvoicesFrom {
		cypher += `,tbp.sendInvoicesFrom=$sendInvoicesFrom`
		params["sendInvoicesFrom"] = data.SendInvoicesFrom
	}
	if data.UpdateSendInvoicesBcc {
		cypher += `,tbp.sendInvoicesBcc=$sendInvoicesBcc`
		params["sendInvoicesBcc"] = data.SendInvoicesBcc
	}
	if data.UpdateCanPayWithPigeon {
		cypher += `,tbp.canPayWithPigeon=$canPayWithPigeon`
		params["canPayWithPigeon"] = data.CanPayWithPigeon
	}
	if data.UpdateCanPayWithBankTransfer {
		cypher += `,tbp.canPayWithBankTransfer=$canPayWithBankTransfer`
		params["canPayWithBankTransfer"] = data.CanPayWithBankTransfer
	}
	if data.UpdateCheck {
		cypher += `,tbp.check=$check`
		params["check"] = data.Check
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *tenantWriteRepository) UpdateTenantSettings(ctx context.Context, tenant string, data TenantSettingsFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantWriteRepository.UpdateTenantBillingProfile")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (t:Tenant {name:$tenant})
				MERGE (t)-[:HAS_SETTINGS]->(ts:TenantSettings {tenant:$tenant})
				ON CREATE SET
					ts.id=randomUUID(),
					ts.createdAt=$now
				SET
					ts.updatedAt=datetime()`
	params := map[string]any{
		"tenant": tenant,
		"now":    utils.Now(),
	}
	if data.UpdateInvoicingEnabled {
		cypher += ", ts.invoicingEnabled=$invoicingEnabled"
		params["invoicingEnabled"] = data.InvoicingEnabled
	}
	if data.UpdateInvoicingPostpaid {
		cypher += ", ts.invoicingPostpaid=$invoicingPostpaid"
		params["invoicingPostpaid"] = data.InvoicingPostpaid
	}
	if data.UpdateBaseCurrency {
		cypher += ", ts.baseCurrency=$baseCurrency"
		params["baseCurrency"] = data.BaseCurrency.String()
	}
	if data.UpdateLogoRepositoryFileId {
		cypher += ", ts.logoRepositoryFileId=$logoRepositoryFileId"
		params["logoRepositoryFileId"] = data.LogoRepositoryFileId
	}
	if data.UpdateWorkspaceLogo {
		cypher += ", ts.workspaceLogo=$workspaceLogo"
		params["workspaceLogo"] = data.WorkspaceLogo
	}
	if data.UpdateWorkspaceName {
		cypher += ", ts.workspaceName=$workspaceName"
		params["workspaceName"] = data.WorkspaceName
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
		model2.NodeLabelTenantBillingProfile,
		model2.NodeLabelBankAccount,
		model2.NodeLabelTimelineEvent,
		model2.NodeLabelContact,
		model2.NodeLabelCustomField,
		model2.NodeLabelJobRole,
		model2.NodeLabelEmail,
		model2.NodeLabelLocation,
		model2.NodeLabelInteractionEvent,
		model2.NodeLabelInteractionSession,
		model2.NodeLabelNote,
		model2.NodeLabelLogEntry,
		model2.NodeLabelOrganization,
		model2.NodeLabelBillingProfile,
		model2.NodeLabelAction,
		model2.NodeLabelPageView,
		model2.NodeLabelPhoneNumber,
		model2.NodeLabelTag,
		model2.NodeLabelIssue,
		model2.NodeLabelUser,
		model2.NodeLabelAttachment,
		model2.NodeLabelMeeting,
		model2.NodeLabelSocial,
		model2.NodeLabelActionItem,
		model2.NodeLabelComment,
		model2.NodeLabelContract,
		model2.NodeLabelDeletedContract,
		model2.NodeLabelServiceLineItem,
		model2.NodeLabelOpportunity,
		model2.NodeLabelInvoicingCycle,
		model2.NodeLabelExternalSystem,
		model2.NodeLabelInvoice,
		model2.NodeLabelInvoiceLine,
		model2.NodeLabelReminder,
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
