package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContactService interface {
	CreateContact(ctx context.Context, tenant string, contactFields neo4jrepository.ContactFields, socialUrl string, externalSystem neo4jmodel.ExternalSystem) (string, error)
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

func (s *contactService) CreateContact(ctx context.Context, tenant string, contactFields neo4jrepository.ContactFields, socialUrl string, externalSystem neo4jmodel.ExternalSystem) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactService.CreateContact")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "contactFields", contactFields)
	tracing.LogObjectAsJson(span, "externalSystem", externalSystem)
	span.LogKV("socialUrl", socialUrl)

	if tenant == "" {
		tenant = common.GetTenantFromContext(ctx)
	}

	// prepare contact id
	contactId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelContact)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tracing.TagEntity(span, contactId)

	// if createdAt missing, set it to now
	contactFields.CreatedAt = utils.TimeOrNow(contactFields.CreatedAt)

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		innerErr := s.services.Neo4jRepositories.ContactWriteRepository.CreateContactInTx(ctx, tx, tenant, contactId, contactFields)
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

	// TODO move link with social url to a (sync process + event)

	// TODO create new proto event for contact creation (after deprecating existing event)
	// send contact to events
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
				Source:    contactFields.SourceFields.Source,
				AppSource: contactFields.SourceFields.AppSource,
			},
			LoggedInUserId: common.GetUserIdFromContext(ctx),
			SocialUrl:      socialUrl,
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
		return contactId, err
	}

	span.LogFields(log.Bool("response.contactCreated", true))
	return contactId, nil
}
