package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContactService interface {
	SaveContact(ctx context.Context, id *string, contactFields neo4jrepository.ContactFields, socialUrl string, externalSystem neo4jmodel.ExternalSystem) (string, error)
	HideContact(ctx context.Context, contactId string) error
	ShowContact(ctx context.Context, contactId string) error
	GetContactById(ctx context.Context, contactId string) (*neo4jentity.ContactEntity, error)
}

type contactService struct {
	log      logger.Logger
	services *Services
}

func NewContactService(log logger.Logger, services *Services) ContactService {
	return &contactService{
		log:      log,
		services: services,
	}
}

func (s *contactService) SaveContact(ctx context.Context, id *string, contactFields neo4jrepository.ContactFields, socialUrl string, externalSystem neo4jmodel.ExternalSystem) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.SaveContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "contactFields", contactFields)
	tracing.LogObjectAsJson(span, "externalSystem", externalSystem)
	span.LogKV("socialUrl", socialUrl)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// set default values
	if contactFields.SourceFields.AppSource != "" {
		common.SetAppSourceInContext(ctx, contactFields.SourceFields.AppSource)
	}

	createFlow := false
	contactId := ""

	// TODO add here any dedup logic

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")
		contactId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelContact)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		// if createdAt missing, set it to now
		contactFields.CreatedAt = utils.TimeOrNow(contactFields.CreatedAt)
	} else {
		span.LogKV("flow", "update")
		contactId = *id

		// validate contact exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, contactId, model.NodeLabelContact)
		if err != nil || !exists {
			err = errors.New("contact not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, contactId)

	_, err = utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		innerErr := s.services.Neo4jRepositories.ContactWriteRepository.SaveContactInTx(ctx, &tx, tenant, contactId, contactFields)
		if innerErr != nil {
			s.log.Errorf("Error while saving contact %s: %s", contactId, err.Error())
			return nil, innerErr
		}
		if externalSystem.Available() {
			innerErr = s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, tenant, contactId, model.NodeLabelContact, externalSystem)
			if err != nil {
				s.log.Errorf("Error while link contact %s with external system %s: %s", contactId, externalSystem.ExternalSystemId, err.Error())
				return nil, innerErr
			}
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createFlow {
		err = s.services.RabbitMQService.Publish(ctx, contactId, model.CONTACT, dto.New_CreateContact_From_ContactFields(contactFields, externalSystem))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateContact"))
		}
		utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())
	} else {
		err = s.services.RabbitMQService.Publish(ctx, contactId, model.CONTACT, dto.New_UpdateContact_From_ContactFields(contactFields, externalSystem))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateContact"))
		}
		if contactFields.SourceFields.AppSource != constants.AppSourceCustomerOsApi {
			utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		}
	}

	if createFlow && socialUrl != "" {
		_, err := s.services.SocialService.MergeSocialWithEntity(ctx,
			LinkWith{
				Id:   contactId,
				Type: model.CONTACT,
			},
			neo4jentity.SocialEntity{
				Url:       socialUrl,
				Source:    neo4jentity.DecodeDataSource(contactFields.SourceFields.Source),
				AppSource: contactFields.SourceFields.AppSource,
			})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to merge social with contact"))
		}
	}

	span.LogFields(log.Bool("response.contactCreated", true))
	return contactId, nil
}

func (s *contactService) HideContact(ctx context.Context, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.HideContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	err = s.services.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contactId, neo4jentity.ContactPropertyHide, true)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while hiding contact %s: %s", contactId, err.Error())
	}
	err = s.services.Neo4jRepositories.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelContact, contactId, string(neo4jentity.ContactPropertyHiddenAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while updating hidden at property for contact %s: %s", contactId, err.Error())
	}

	err = s.services.RabbitMQService.Publish(ctx, contactId, model.CONTACT, dto.HideContact{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message HideContact"))
	}

	utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithDelete())

	return nil
}

func (s *contactService) ShowContact(ctx context.Context, contactId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.ShowContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	err = s.services.Neo4jRepositories.ContactWriteRepository.UpdateAnyProperty(ctx, tenant, contactId, neo4jentity.ContactPropertyHide, false)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while showing contact %s: %s", contactId, err.Error())
		return err
	}

	err = s.services.RabbitMQService.Publish(ctx, contactId, model.CONTACT, dto.ShowContact{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ShowContact"))
	}

	utils.EventCompleted(ctx, tenant, model.CONTACT.String(), contactId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())

	return nil
}

func (s *contactService) GetContactById(ctx context.Context, contactId string) (*neo4jentity.ContactEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.GetContactById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, contactId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	contactDbNode, err := s.services.Neo4jRepositories.ContactReadRepository.GetContact(ctx, tenant, contactId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while getting contact %s: %s", contactId, err.Error())
		return nil, err
	}

	return neo4jmapper.MapDbNodeToContactEntity(contactDbNode), nil
}
