package service

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	phonenumberpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/phone_number"
	userpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/user"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

type PhoneNumberService interface {
	CreatePhoneNumberViaEvents(ctx context.Context, phoneNumber, appSource string) (string, error)
	UpdatePhoneNumberFor(ctx context.Context, entityType entity.EntityType, entityId string, input model.PhoneNumberRelationUpdateInput) error
	DetachFromEntityByPhoneNumber(ctx context.Context, entityType entity.EntityType, entityId, phoneNumber string) (bool, error)
	DetachFromEntityById(ctx context.Context, entityType entity.EntityType, entityId, phoneNumberId string) (bool, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType entity.EntityType, ids []string) (*entity.PhoneNumberEntities, error)
	GetById(ctx context.Context, phoneNumberId string) (*neo4jentity.PhoneNumberEntity, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.PhoneNumberEntity, error)
	Update(ctx context.Context, input model.PhoneNumberUpdateInput) error

	mapDbNodeToPhoneNumberEntity(node dbtype.Node) *entity.PhoneNumberEntity
}

type phoneNumberService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewPhoneNumberService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) PhoneNumberService {
	return &phoneNumberService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *phoneNumberService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *phoneNumberService) CreatePhoneNumberViaEvents(ctx context.Context, phoneNumber, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.CreatePhoneNumberViaEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber), log.String("appSource", appSource))

	phoneNumber = strings.TrimSpace(phoneNumber)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*phonenumberpb.PhoneNumberIdGrpcResponse](func() (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
		return s.grpcClients.PhoneNumberClient.UpsertPhoneNumber(ctx, &phonenumberpb.UpsertPhoneNumberGrpcRequest{
			Tenant:      common.GetTenantFromContext(ctx),
			PhoneNumber: phoneNumber,
			SourceFields: &commonpb.SourceFields{
				Source:    string(neo4jentity.DataSourceOpenline),
				AppSource: utils.StringFirstNonEmpty(appSource, constants.AppSourceCustomerOsApi),
			},
			LoggedInUserId: common.GetUserIdFromContext(ctx),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelPhoneNumber, span)
	return response.Id, nil
}

func (s *phoneNumberService) GetAllForEntityTypeByIds(ctx context.Context, entityType entity.EntityType, ids []string) (*entity.PhoneNumberEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetAllForEntityTypeByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("ids", ids))

	phoneNumbers, err := s.repositories.PhoneNumberRepository.GetAllForIds(ctx, common.GetTenantFromContext(ctx), entityType, ids)
	if err != nil {
		return nil, err
	}

	phoneNumberEntities := make(entity.PhoneNumberEntities, 0, len(phoneNumbers))
	for _, v := range phoneNumbers {
		phoneNumberEntity := s.mapDbNodeToPhoneNumberEntity(*v.Node)
		s.addDbRelationshipToPhoneNumberEntity(*v.Relationship, phoneNumberEntity)
		phoneNumberEntity.DataloaderKey = v.LinkedNodeId
		phoneNumberEntities = append(phoneNumberEntities, *phoneNumberEntity)
	}
	return &phoneNumberEntities, nil
}

func (s *phoneNumberService) UpdatePhoneNumberFor(ctx context.Context, entityType entity.EntityType, entityId string, input model.PhoneNumberRelationUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.UpdatePhoneNumberFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", string(entityType)), log.String("entityId", entityId))

	phoneNumberEntity, err := s.GetById(ctx, input.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if entityType == entity.CONTACT {
		contactExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), entityId, neo4jutil.NodeLabelContact)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if !contactExists {
			err = errors.New("Contact not found")
			tracing.TraceErr(span, err)
			return err
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
			return s.grpcClients.ContactClient.LinkPhoneNumberToContact(ctx, &contactpb.LinkPhoneNumberToContactGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				ContactId:      entityId,
				PhoneNumberId:  phoneNumberEntity.Id,
				Primary:        utils.IfNotNilBool(input.Primary),
				Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			graphql.AddErrorf(ctx, "Failed to add phone number %s to contact %s", input.ID, entityId)
			return err
		}
	} else if entityType == entity.ORGANIZATION {
		organizationExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), entityId, neo4jutil.NodeLabelOrganization)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if !organizationExists {
			err = errors.New("Organization not found")
			tracing.TraceErr(span, err)
			return err
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.grpcClients.OrganizationClient.LinkPhoneNumberToOrganization(ctx, &organizationpb.LinkPhoneNumberToOrganizationGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				OrganizationId: entityId,
				PhoneNumberId:  phoneNumberEntity.Id,
				Primary:        utils.IfNotNilBool(input.Primary),
				Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			graphql.AddErrorf(ctx, "Failed to add phone number %s to organization %s", input.ID, entityId)
			return err
		}
	} else if entityType == entity.USER {
		userExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), entityId, neo4jutil.NodeLabelUser)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if !userExists {
			err = errors.New("User not found")
			tracing.TraceErr(span, err)
			return err
		}

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*userpb.UserIdGrpcResponse](func() (*userpb.UserIdGrpcResponse, error) {
			return s.grpcClients.UserClient.LinkPhoneNumberToUser(ctx, &userpb.LinkPhoneNumberToUserGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				UserId:         entityId,
				PhoneNumberId:  phoneNumberEntity.Id,
				Primary:        utils.IfNotNilBool(input.Primary),
				Label:          utils.IfNotNilString(input.Label, func() string { return input.Label.String() }),
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			graphql.AddErrorf(ctx, "Failed to add phone number %s to user %s", input.ID, entityId)
			return err
		}
	}

	return nil
}

func (s *phoneNumberService) DetachFromEntityByPhoneNumber(ctx context.Context, entityType entity.EntityType, entityId, phoneNumber string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.DetachFromEntityByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId))

	err := s.repositories.PhoneNumberRepository.RemoveRelationship(ctx, entityType, common.GetTenantFromContext(ctx), entityId, phoneNumber)

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *phoneNumberService) DetachFromEntityById(ctx context.Context, entityType entity.EntityType, entityId, phoneNumberId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.DetachFromEntityById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId), log.String("phoneNumberId", phoneNumberId))

	err := s.repositories.PhoneNumberRepository.RemoveRelationshipById(ctx, entityType, common.GetTenantFromContext(ctx), entityId, phoneNumberId)

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *phoneNumberService) GetById(ctx context.Context, phoneNumberId string) (*neo4jentity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumberId", phoneNumberId))

	phoneNumberNode, err := s.repositories.Neo4jRepositories.PhoneNumberReadRepository.GetById(ctx, common.GetTenantFromContext(ctx), phoneNumberId)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToPhoneNumberEntity(phoneNumberNode), nil
}

func (s *phoneNumberService) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*entity.PhoneNumberEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.GetByPhoneNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("phoneNumber", phoneNumber))

	phoneNumberNode, err := s.repositories.PhoneNumberRepository.GetByPhoneNumber(ctx, common.GetTenantFromContext(ctx), phoneNumber)
	if err != nil {
		return nil, err
	}
	var phoneNumberEntity = s.mapDbNodeToPhoneNumberEntity(*phoneNumberNode)
	return phoneNumberEntity, nil
}

// Deprecated
func (s *phoneNumberService) mapDbNodeToPhoneNumberEntity(node dbtype.Node) *entity.PhoneNumberEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.PhoneNumberEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		E164:           utils.GetStringPropOrEmpty(props, "e164"),
		RawPhoneNumber: utils.GetStringPropOrEmpty(props, "rawPhoneNumber"),
		Validated:      utils.GetBoolPropOrNil(props, "validated"),
		Source:         neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &result
}

func (s *phoneNumberService) addDbRelationshipToPhoneNumberEntity(relationship dbtype.Relationship, phoneNumberEntity *entity.PhoneNumberEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	phoneNumberEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
	phoneNumberEntity.Label = utils.GetStringPropOrEmpty(props, "label")
}

func (s *phoneNumberService) Update(ctx context.Context, input model.PhoneNumberUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PhoneNumberService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	phoneNumberEntity, err := s.GetById(ctx, input.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if phoneNumberEntity.RawPhoneNumber == input.PhoneNumber {
		err = errors.New("Phone number is the same as the current one")
		tracing.TraceErr(span, err)
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*phonenumberpb.PhoneNumberIdGrpcResponse](func() (*phonenumberpb.PhoneNumberIdGrpcResponse, error) {
		return s.grpcClients.PhoneNumberClient.UpsertPhoneNumber(ctx, &phonenumberpb.UpsertPhoneNumberGrpcRequest{
			Id:          input.ID,
			Tenant:      common.GetTenantFromContext(ctx),
			PhoneNumber: input.PhoneNumber,
			SourceFields: &commonpb.SourceFields{
				Source:    string(neo4jentity.DataSourceOpenline),
				AppSource: constants.AppSourceCustomerOsApi,
			},
			LoggedInUserId: common.GetUserIdFromContext(ctx),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing %s", err.Error())
		return err
	}

	return err
}
