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
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContactService interface {
	SaveContact(ctx context.Context, id *string, contactFields neo4jrepository.ContactFields, socialUrl string, externalSystem neo4jmodel.ExternalSystem) (string, error)
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

	// TODO remove below block of sending events to event-store-db once fully deprecated
	if createFlow {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.services.GrpcClients.ContactClient.UpsertContact(ctx, &contactpb.UpsertContactGrpcRequest{
				Id:              contactId,
				Tenant:          tenant,
				FirstName:       contactFields.FirstName,
				LastName:        contactFields.LastName,
				Name:            contactFields.Name,
				Prefix:          contactFields.Prefix,
				CreatedAt:       utils.ConvertTimeToTimestampPtr(&contactFields.CreatedAt),
				Description:     contactFields.Description,
				Timezone:        contactFields.Timezone,
				ProfilePhotoUrl: contactFields.ProfilePhotoUrl,
				SourceFields: &commonpb.SourceFields{
					Source:    contactFields.SourceFields.GetSource(),
					AppSource: contactFields.SourceFields.GetAppSource(),
				},
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				Username:       contactFields.Username,
				ExternalSystemFields: &commonpb.ExternalSystemFields{
					ExternalSystemId: externalSystem.ExternalSystemId,
					ExternalUrl:      externalSystem.ExternalUrl,
					ExternalId:       externalSystem.ExternalId,
					ExternalSource:   externalSystem.ExternalSource,
					ExternalIdSecond: externalSystem.ExternalIdSecond,
					SyncDate:         utils.ConvertTimeToTimestampPtr(externalSystem.SyncDate),
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to create contact in events platform"))
		}
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
		_, err := s.services.SocialService.MergeSocialWithEntity(ctx, tenant, contactId, model.CONTACT,
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
