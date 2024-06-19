package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
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
	commongrpc "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	emailpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/email"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
)

type EmailService interface {
	CreateEmailAddressViaEvents(ctx context.Context, email, appSource string) (string, error)
	GetAllFor(ctx context.Context, entityType entity.EntityType, entityId string) (*entity.EmailEntities, error)
	GetAllForEntityTypeByIds(ctx context.Context, entityType entity.EntityType, entityIds []string) (*entity.EmailEntities, error)
	UpdateEmailFor(ctx context.Context, entityType entity.EntityType, entityId string, entity *entity.EmailEntity) (*entity.EmailEntity, error)
	DetachFromEntity(ctx context.Context, entityType entity.EntityType, entityId, email string) (bool, error)
	DetachFromEntityById(ctx context.Context, entityType entity.EntityType, entityId, emailId string) (bool, error)
	DeleteById(ctx context.Context, emailId string) (bool, error)
	GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error)
	GetByEmailAddress(ctx context.Context, email string) (*entity.EmailEntity, error)

	mapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity
}

type emailService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
	grpcClients  *grpc_client.Clients
}

func NewEmailService(log logger.Logger, repositories *repository.Repositories, services *Services, grpcClients *grpc_client.Clients) EmailService {
	return &emailService{
		log:          log,
		repositories: repositories,
		services:     services,
		grpcClients:  grpcClients,
	}
}

func (s *emailService) getDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *emailService) GetAllFor(ctx context.Context, entityType entity.EntityType, entityId string) (*entity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetAllFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.String("entityId", entityId))

	records, err := s.repositories.EmailRepository.GetAllFor(ctx, common.GetContext(ctx).Tenant, entityType, entityId)
	if err != nil {
		return nil, err
	}

	emailEntities := make(entity.EmailEntities, 0, len(records))
	for _, dbRecord := range records {
		emailEntity := s.mapDbNodeToEmailEntity(dbRecord.Values[0].(dbtype.Node))
		s.addDbRelationshipToEmailEntity(dbRecord.Values[1].(dbtype.Relationship), emailEntity)
		emailEntities = append(emailEntities, *emailEntity)
	}

	return &emailEntities, nil
}

func (s *emailService) GetAllForEntityTypeByIds(ctx context.Context, entityType entity.EntityType, entityIds []string) (*entity.EmailEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetAllForEntityTypeByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", entityType.String()), log.Object("entityIds", entityIds))

	emails, err := s.repositories.EmailRepository.GetAllForIds(ctx, common.GetContext(ctx).Tenant, entityType, entityIds)
	if err != nil {
		return nil, err
	}

	emailEntities := make(entity.EmailEntities, 0, len(emails))
	for _, v := range emails {
		emailEntity := s.mapDbNodeToEmailEntity(*v.Node)
		s.addDbRelationshipToEmailEntity(*v.Relationship, emailEntity)
		emailEntity.DataloaderKey = v.LinkedNodeId
		emailEntities = append(emailEntities, *emailEntity)
	}
	return &emailEntities, nil
}

func (s *emailService) UpdateEmailFor(ctx context.Context, entityType entity.EntityType, entityId string, inputEntity *entity.EmailEntity) (*entity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.UpdateEmailFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("entityType", string(entityType)), log.String("entityId", entityId))

	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	var err error
	var emailNode *dbtype.Node
	var emailRelationship *dbtype.Relationship
	var detachCurrentEmail = false

	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		currentEmailNode, err := s.repositories.EmailRepository.GetByIdAndRelatedEntity(ctx, entityType, common.GetTenantFromContext(ctx), inputEntity.Id, entityId)
		if err != nil {
			return nil, err
		}

		currentEmail := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*currentEmailNode), "email")
		currentRawEmail := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*currentEmailNode), "rawEmail")

		var emailExists = false
		if currentRawEmail == "" {
			emailExists, err = s.repositories.EmailRepository.Exists(ctx, common.GetContext(ctx).Tenant, inputEntity.RawEmail)
			if err != nil {
				return nil, err
			}
		}

		if len(inputEntity.RawEmail) == 0 || inputEntity.RawEmail == currentEmail || inputEntity.RawEmail == currentRawEmail ||
			(currentRawEmail == "" && !emailExists) {
			// proceed with update
			emailNode, emailRelationship, err = s.repositories.EmailRepository.UpdateEmailForInTx(ctx, tx, common.GetContext(ctx).Tenant, entityType, entityId, *inputEntity)
			if err != nil {
				return nil, err
			}
			emailId := utils.GetPropsFromNode(*emailNode)["id"].(string)
			if inputEntity.Primary == true {
				err := s.repositories.EmailRepository.SetOtherEmailsNonPrimaryInTx(ctx, tx, common.GetContext(ctx).Tenant, entityType, entityId, emailId)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// proceed with email address replace
			// merge new email address
			emailNode, emailRelationship, err = s.repositories.EmailRepository.MergeEmailToInTx(ctx, tx, common.GetContext(ctx).Tenant, entityType, entityId, *inputEntity)
			if err != nil {
				return nil, err
			}
			emailId := utils.GetPropsFromNode(*emailNode)["id"].(string)
			if inputEntity.Primary == true {
				err := s.repositories.EmailRepository.SetOtherEmailsNonPrimaryInTx(ctx, tx, common.GetContext(ctx).Tenant, entityType, entityId, emailId)
				if err != nil {
					return nil, err
				}
			}
			detachCurrentEmail = true
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	if detachCurrentEmail {
		_, err = s.DetachFromEntityById(ctx, entityType, entityId, inputEntity.Id)
	}

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	}

	var emailEntity = s.mapDbNodeToEmailEntity(*emailNode)
	s.addDbRelationshipToEmailEntity(*emailRelationship, emailEntity)
	return emailEntity, nil
}

func (s *emailService) DetachFromEntity(ctx context.Context, entityType entity.EntityType, entityId, email string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DetachFromEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email), log.String("entityId", entityId), log.String("entityType", string(entityType)))

	err := s.repositories.EmailRepository.RemoveRelationship(ctx, entityType, common.GetTenantFromContext(ctx), entityId, email)

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *emailService) DetachFromEntityById(ctx context.Context, entityType entity.EntityType, entityId, emailId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DetachFromEntityById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", emailId), log.String("entityId", entityId), log.String("entityType", string(entityType)))

	err := s.repositories.EmailRepository.RemoveRelationshipById(ctx, entityType, common.GetTenantFromContext(ctx), entityId, emailId)

	if entityType == entity.ORGANIZATION {
		s.services.OrganizationService.UpdateLastTouchpoint(ctx, entityId)
	} else if entityType == entity.CONTACT {
		s.services.OrganizationService.UpdateLastTouchpointByContactId(ctx, entityId)
	}

	return err == nil, err
}

func (s *emailService) DeleteById(ctx context.Context, emailId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.DeleteById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", emailId))

	err := s.repositories.EmailRepository.DeleteById(ctx, common.GetTenantFromContext(ctx), emailId)
	return err == nil, err
}

func (s *emailService) GetById(ctx context.Context, emailId string) (*neo4jentity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("emailId", emailId))

	emailNode, err := s.repositories.EmailRepository.GetById(ctx, emailId)
	if err != nil {
		return nil, err
	}
	var emailEntity = neo4jmapper.MapDbNodeToEmailEntity(emailNode)
	return emailEntity, nil
}

func (s *emailService) GetByEmailAddress(ctx context.Context, email string) (*entity.EmailEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.GetByEmailAddress")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email))

	emailNode, err := s.repositories.EmailRepository.GetByEmail(ctx, common.GetTenantFromContext(ctx), email)
	if err != nil {
		return nil, err
	}
	var emailEntity = s.mapDbNodeToEmailEntity(*emailNode)
	return emailEntity, nil
}

func (s *emailService) CreateEmailAddressViaEvents(ctx context.Context, email, appSource string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailService.CreateEmailAddressViaEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("email", email), log.String("appSource", appSource))

	email = strings.TrimSpace(email)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*emailpb.EmailIdGrpcResponse](func() (*emailpb.EmailIdGrpcResponse, error) {
		return s.grpcClients.EmailClient.UpsertEmail(ctx, &emailpb.UpsertEmailGrpcRequest{
			Tenant:   common.GetTenantFromContext(ctx),
			RawEmail: email,
			SourceFields: &commongrpc.SourceFields{
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

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelEmail, span)
	return response.Id, nil
}

func (s *emailService) mapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.EmailEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Email:         utils.GetStringPropOrEmpty(props, "email"),
		RawEmail:      utils.GetStringPropOrEmpty(props, "rawEmail"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),

		Validated:      utils.GetBoolPropOrNil(props, "validated"),
		IsReachable:    utils.GetStringPropOrNil(props, "isReachable"),
		IsValidSyntax:  utils.GetBoolPropOrNil(props, "isValidSyntax"),
		CanConnectSMTP: utils.GetBoolPropOrNil(props, "canConnectSmtp"),
		AcceptsMail:    utils.GetBoolPropOrNil(props, "acceptsMail"),
		HasFullInbox:   utils.GetBoolPropOrNil(props, "hasFullInbox"),
		IsCatchAll:     utils.GetBoolPropOrNil(props, "isCatchAll"),
		IsDeliverable:  utils.GetBoolPropOrNil(props, "isDeliverable"),
		IsDisabled:     utils.GetBoolPropOrNil(props, "isDisabled"),
		Error:          utils.GetStringPropOrNil(props, "validationError"),
	}
	return &result
}

func (s *emailService) addDbRelationshipToEmailEntity(relationship dbtype.Relationship, emailEntity *entity.EmailEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	emailEntity.Primary = utils.GetBoolPropOrFalse(props, "primary")
	emailEntity.Label = utils.GetStringPropOrEmpty(props, "label")
}
