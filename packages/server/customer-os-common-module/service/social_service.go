package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	socialpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/social"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type SocialService interface {
	Update(ctx context.Context, tenant string, entity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error)
	GetAllForEntities(ctx context.Context, tenant string, linkedEntityType model.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error)
	Remove(ctx context.Context, tenant, socialId string) error
	MergeSocialWithEntity(ctx context.Context, tenant, linkedEntityId string, linkedEntityType model.EntityType, socialEntity neo4jentity.SocialEntity) (string, error)
}

type socialService struct {
	log      logger.Logger
	services *Services
}

func NewSocialService(log logger.Logger, services *Services) SocialService {
	return &socialService{
		log:      log,
		services: services,
	}
}

func (s *socialService) GetAllForEntities(ctx context.Context, tenant string, linkedEntityType model.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("linkedEntityType", string(linkedEntityType)), log.Object("linkedEntityIds", linkedEntityIds))

	socials, err := s.services.Neo4jRepositories.SocialReadRepository.GetAllForEntities(ctx, tenant, linkedEntityType, linkedEntityIds)
	if err != nil {
		return nil, err
	}
	socialEntities := make(neo4jentity.SocialEntities, 0, len(socials))
	for _, v := range socials {
		socialEntity := neo4jmapper.MapDbNodeToSocialEntity(v.Node)
		socialEntity.DataloaderKey = v.LinkedNodeId
		socialEntities = append(socialEntities, *socialEntity)
	}
	return &socialEntities, nil
}

func (s *socialService) Update(ctx context.Context, tenant string, socialEntity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	updatedLocationNode, err := s.services.Neo4jRepositories.SocialWriteRepository.Update(ctx, tenant, socialEntity)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToSocialEntity(updatedLocationNode), nil
}

func (s *socialService) Remove(ctx context.Context, tenant string, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Remove")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, socialId)

	return s.services.Neo4jRepositories.SocialWriteRepository.PermanentlyDelete(ctx, tenant, socialId)
}

func (s *socialService) MergeSocialWithEntity(ctx context.Context, tenant, linkedEntityId string, linkedEntityType model.EntityType, socialEntity neo4jentity.SocialEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.MergeSocialWithEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("linkedEntityId", linkedEntityId), log.String("linkedEntityType", string(linkedEntityType)))

	if tenant == "" {
		tenant = common.GetTenantFromContext(ctx)
	}

	socialId := socialEntity.Id
	if socialId == "" {
		var err error
		socialId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelSocial)
		if err != nil {
			return "", err
		}
	}
	tracing.TagEntity(span, socialId)

	// save social to neo4j
	data := neo4jrepository.SocialFields{
		SocialId:       socialId,
		Url:            socialEntity.Url,
		Alias:          socialEntity.Alias,
		ExternalId:     socialEntity.ExternalId,
		FollowersCount: socialEntity.FollowersCount,
		CreatedAt:      utils.TimeOrNow(socialEntity.CreatedAt),
		SourceFields: neo4jmodel.Source{
			Source:    neo4jmodel.GetSource(socialEntity.Source.String()),
			AppSource: neo4jmodel.GetAppSource(socialEntity.AppSource),
		},
	}
	err := s.services.Neo4jRepositories.SocialWriteRepository.MergeSocialForEntity(ctx, tenant, linkedEntityId, linkedEntityType.Neo4jLabel(), data)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// send event to event store
	if linkedEntityType == model.CONTACT {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*socialpb.SocialIdGrpcResponse](func() (*socialpb.SocialIdGrpcResponse, error) {
			return s.services.GrpcClients.ContactClient.AddSocial(ctx, &contactpb.ContactAddSocialGrpcRequest{
				Tenant:         tenant,
				ContactId:      linkedEntityId,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				SourceFields: &commonpb.SourceFields{
					Source:    neo4jmodel.GetSource(socialEntity.Source.String()),
					AppSource: neo4jmodel.GetAppSource(socialEntity.AppSource),
				},
				Url:            socialEntity.Url,
				CreatedAt:      utils.ConvertTimeToTimestampPtr(utils.TimePtr(utils.TimeOrNow(socialEntity.CreatedAt))),
				SocialId:       socialId,
				Alias:          socialEntity.Alias,
				ExternalId:     socialEntity.ExternalId,
				FollowersCount: socialEntity.FollowersCount,
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error while sending event"))
			return socialId, err
		}
	}

	return socialId, nil
}
