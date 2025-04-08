package tenant_settings

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type tenantSettingsService struct {
	log    logger.Logger
	neo4j  *neo4j_repository.Repositories
	events *events.EventsService
}

func NewTenantSettingsService(log logger.Logger, neo4j *neo4j_repository.Repositories, events *events.EventsService) interfaces.TenantSettingsService {
	return &tenantSettingsService{
		log:    log,
		neo4j:  neo4j,
		events: events,
	}
}

func (s *tenantSettingsService) GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.GetTenantSettings")
	defer spans.Finish()

	dbNode, err := s.neo4j.TenantReadRepository.GetTenantSettings(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantSettingsService) GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.GetTenantSettingsForTenant")
	defer spans.Finish()

	dbNode, err := s.neo4j.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantSettingsService) UpdateTenantSettings(ctx context.Context, dataFields data_fields.TenantSettingsFields) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.UpdateTenantSettings")
	defer spans.Finish()

	spans.LogObjectAsJson("dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if dataFields.IsEmpty() {
		// nothing to update, return
		return nil
	}

	// update tenant settings in neo4j
	err = s.neo4j.TenantWriteRepository.UpdateTenantSettings(ctx, tenant, dataFields)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("Unable to update tenant settings", err)
		return err
	}

	// send event to RabbitMQ
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT_SETTINGS, dto.UpdateTenantSettings{dataFields})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message UpdateTenantSettings"))
	}

	return nil
}

func (s *tenantSettingsService) GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.GetTenantBillingProfiles")
	defer spans.Finish()

	dbNodes, err := s.neo4j.TenantReadRepository.GetTenantBillingProfiles(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("GetTenantBillingProfiles: %w", err)
	}

	tenantBillingProfiles := neo4jentity.TenantBillingProfileEntities{}
	for _, dbNode := range dbNodes {
		tenantBillingProfiles = append(tenantBillingProfiles, *neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode))
	}

	return &tenantBillingProfiles, nil
}

func (s *tenantSettingsService) GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.GetTenantBillingProfile")
	defer spans.Finish()

	spans.LogKV("id", id)

	dbNode, err := s.neo4j.TenantReadRepository.GetTenantBillingProfileById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("GetTenantBillingProfile: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode), nil
}

func (s *tenantSettingsService) GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.GetDefaultTenantBillingProfile")
	defer spans.Finish()

	tenantBillingProfiles, err := s.GetTenantBillingProfiles(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("GetDefaultTenantBillingProfile: %w", err)
	}
	if tenantBillingProfiles == nil || len(*tenantBillingProfiles) == 0 {
		return nil, nil
	} else {
		return &(*tenantBillingProfiles)[0], nil
	}
}

func (s *tenantSettingsService) CreateBankAccount(ctx context.Context, dataFields data_fields.BankAccountFields) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.CreateBankAccount")
	defer spans.Finish()

	spans.LogObjectAsJson("dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// generate new id
	bankAccountId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelBankAccount)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	dataFields.ID = bankAccountId

	// set default fields
	if dataFields.AppSource == nil {
		dataFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
	}
	if dataFields.Source == nil {
		dataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
	}

	err = s.neo4j.BankAccountWriteRepository.CreateBankAccount(ctx, tenant, dataFields)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("Unable to create bank account", err)
		return "", err
	}

	// send event to RabbitMQ
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT_SETTINGS, dto.CreateBankAccount{dataFields})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message CreateBankAccount"))
	}

	return bankAccountId, nil
}

func (s *tenantSettingsService) UpdateBankAccount(ctx context.Context, bankAccountId string, dataFields data_fields.BankAccountFields) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.UpdateBankAccount")
	defer spans.Finish()

	spans.LogObjectAsJson("dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if dataFields.IsEmpty() {
		// nothing to update, return
		return nil
	}

	// verify if bank account exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, bankAccountId, model.NodeLabelBankAccount)
	if err != nil || !exists {
		err = errors.New("bank account not found")
		spans.TraceError(err)
		return err
	}
	dataFields.ID = bankAccountId
	spans.TagEntity(bankAccountId)

	err = s.neo4j.BankAccountWriteRepository.UpdateBankAccount(ctx, tenant, dataFields)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("Unable to create bank account", err)
		return err
	}

	// send event to RabbitMQ
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT_SETTINGS, dto.UpdateBankAccount{dataFields})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message UpdateBankAccount"))
	}

	return nil
}

func (s *tenantSettingsService) DeleteBankAccount(ctx context.Context, bankAccountId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.DeleteBankAccount")
	defer spans.Finish()

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// verify if bank account exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, bankAccountId, model.NodeLabelBankAccount)
	if err != nil || !exists {
		err = errors.New("bank account not found")
		spans.TraceError(err)
		return err
	}
	spans.TagEntity(bankAccountId)

	err = s.neo4j.BankAccountWriteRepository.DeleteBankAccount(ctx, tenant, bankAccountId)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("Unable to delete bank account", err)
		return err
	}

	// send event to RabbitMQ
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT_SETTINGS, dto.DeleteBankAccount{ID: bankAccountId})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message DeleteBankAccount"))
	}

	return nil
}

func (s *tenantSettingsService) CreateTenantBillingProfile(ctx context.Context, dataFields data_fields.TenantBillingProfileFields) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.CreateTenantBillingProfile")
	defer spans.Finish()

	spans.LogObjectAsJson("dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// generate new id
	profileId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelTenantBillingProfile)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	dataFields.ID = profileId
	spans.TagEntity(profileId)

	// set default fields
	if dataFields.AppSource == nil {
		dataFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
	}
	if dataFields.Source == nil {
		dataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
	}

	// save to db
	err = s.neo4j.TenantWriteRepository.CreateTenantBillingProfile(ctx, tenant, dataFields)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("Unable to create bank account", err)
		return "", err
	}

	// send event to RabbitMQ
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT_SETTINGS, dto.CreateTenantBillingProfile{dataFields})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message CreateTenantBillingProfile"))
	}

	return profileId, nil
}

func (s *tenantSettingsService) UpdateTenantBillingProfile(ctx context.Context, profileId string, dataFields data_fields.TenantBillingProfileFields) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TenantSettingsService.UpdateTenantBillingProfile")
	defer spans.Finish()

	spans.LogObjectAsJson("dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if dataFields.IsEmpty() {
		// nothing to update, return
		return nil
	}

	// verify if bank account exists
	exists, err := s.neo4j.CommonReadRepository.ExistsById(ctx, tenant, profileId, model.NodeLabelTenantBillingProfile)
	if err != nil || !exists {
		err = errors.New("tenant billing profile not found")
		spans.TraceError(err)
		return err
	}
	dataFields.ID = profileId
	spans.TagEntity(profileId)

	err = s.neo4j.TenantWriteRepository.UpdateTenantBillingProfile(ctx, tenant, dataFields)
	if err != nil {
		spans.TraceError(err)
		s.log.Error("Unable to create bank account", err)
		return err
	}

	// send event to RabbitMQ
	err = s.events.Publisher.PublishFanoutEvent(ctx, tenant, model.TENANT_SETTINGS, dto.UpdateTenantBillingProfile{dataFields})
	if err != nil {
		spans.TraceError(errors.Wrap(err, "unable to publish message UpdateTenantBillingProfile"))
	}

	return nil
}
