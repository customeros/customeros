package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type TenantSettingsService interface {
	GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error)
	GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error)
	UpdateTenantSettings(ctx context.Context, dataFields data_fields.TenantSettingsFields) error

	GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error)
	GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error)
	GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error)

	CreateBankAccount(ctx context.Context, dataFields data_fields.BankAccountFields) (string, error)
	UpdateBankAccount(ctx context.Context, bankAccountId string, dataFields data_fields.BankAccountFields) error
	DeleteBankAccount(ctx context.Context, bankAccountId string) error
}

type tenantSettingsService struct {
	log      logger.Logger
	services *Services
}

func NewTenantSettingsService(log logger.Logger, services *Services) TenantSettingsService {
	return &tenantSettingsService{
		log:      log,
		services: services,
	}
}

func (s *tenantSettingsService) GetTenantSettings(ctx context.Context) (*neo4jentity.TenantSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.GetTenantSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantSettingsService) GetTenantSettingsForTenant(ctx context.Context, tenant string) (*neo4jentity.TenantSettingsEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.GetTenantSettingsForTenant")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, tenant)

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode), nil
}

func (s *tenantSettingsService) UpdateTenantSettings(ctx context.Context, dataFields data_fields.TenantSettingsFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.UpdateTenantSettings")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if dataFields.IsEmpty() {
		// nothing to update, return
		return nil
	}

	// update tenant settings in neo4j
	err = s.services.Neo4jRepositories.TenantWriteRepository.UpdateTenantSettings(ctx, tenant, dataFields)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error("Unable to update tenant settings", err)
		return err
	}

	// send event to RabbitMQ
	err = s.services.RabbitMQService.PublishEvent(ctx, tenant, model.TENANT_SETTINGS, dto.UpdateTenantSettings{dataFields})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateTenantSettings"))
	}

	return nil
}

func (s *tenantSettingsService) GetTenantBillingProfiles(ctx context.Context) (*neo4jentity.TenantBillingProfileEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.GetTenantBillingProfiles")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNodes, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfiles(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBillingProfiles: %w", err)
	}

	tenantBillingProfiles := neo4jentity.TenantBillingProfileEntities{}
	for _, dbNode := range dbNodes {
		tenantBillingProfiles = append(tenantBillingProfiles, *neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode))
	}

	return &tenantBillingProfiles, nil
}

func (s *tenantSettingsService) GetTenantBillingProfile(ctx context.Context, id string) (*neo4jentity.TenantBillingProfileEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.GetTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("id", id))

	dbNode, err := s.services.Neo4jRepositories.TenantReadRepository.GetTenantBillingProfileById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBillingProfile: %w", err)
	}

	return neo4jmapper.MapDbNodeToTenantBillingProfileEntity(dbNode), nil
}

func (s *tenantSettingsService) GetDefaultTenantBillingProfile(ctx context.Context) (*neo4jentity.TenantBillingProfileEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.GetDefaultTenantBillingProfile")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenantBillingProfiles, err := s.GetTenantBillingProfiles(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetDefaultTenantBillingProfile: %w", err)
	}
	if tenantBillingProfiles == nil || len(*tenantBillingProfiles) == 0 {
		return nil, nil
	} else {
		return &(*tenantBillingProfiles)[0], nil
	}
}

func (s *tenantSettingsService) CreateBankAccount(ctx context.Context, dataFields data_fields.BankAccountFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.CreateBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// generate new id
	bankAccountId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelBankAccount)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	dataFields.ID = bankAccountId

	err = s.services.Neo4jRepositories.BankAccountWriteRepository.CreateBankAccount(ctx, tenant, dataFields)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error("Unable to create bank account", err)
		return "", err
	}

	// send event to RabbitMQ
	err = s.services.RabbitMQService.PublishEvent(ctx, tenant, model.TENANT_SETTINGS, dto.CreateBankAccount{dataFields})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateBankAccount"))
	}

	return bankAccountId, nil
}

func (s *tenantSettingsService) UpdateBankAccount(ctx context.Context, bankAccountId string, dataFields data_fields.BankAccountFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.UpdateBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// verify if bank account exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, bankAccountId, model.NodeLabelBankAccount)
	if err != nil || !exists {
		err = errors.New("bank account not found")
		tracing.TraceErr(span, err)
		return err
	}
	dataFields.ID = bankAccountId
	tracing.TagEntity(span, bankAccountId)

	err = s.services.Neo4jRepositories.BankAccountWriteRepository.UpdateBankAccount(ctx, tenant, dataFields)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error("Unable to create bank account", err)
		return err
	}

	// send event to RabbitMQ
	err = s.services.RabbitMQService.PublishEvent(ctx, tenant, model.TENANT_SETTINGS, dto.UpdateBankAccount{dataFields})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateBankAccount"))
	}

	return nil
}

func (s *tenantSettingsService) DeleteBankAccount(ctx context.Context, bankAccountId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "TenantSettingsService.DeleteBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, bankAccountId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	err = s.services.Neo4jRepositories.BankAccountWriteRepository.DeleteBankAccount(ctx, tenant, bankAccountId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error("Unable to delete bank account", err)
		return err
	}

	// send event to RabbitMQ
	err = s.services.RabbitMQService.PublishEvent(ctx, tenant, model.TENANT_SETTINGS, dto.DeleteBankAccount{ID: bankAccountId})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message DeleteBankAccount"))
	}

	return nil
}
