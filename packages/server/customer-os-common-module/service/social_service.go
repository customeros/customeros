package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
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
	"strings"
)

type SocialService interface {
	GetById(ctx context.Context, socialId string) (*neo4jentity.SocialEntity, error)
	AddSocialToEntity(ctx context.Context, linkWith LinkWith, socialEntity neo4jentity.SocialEntity) (string, error)
	Update(ctx context.Context, entity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error)
	PermanentlyDelete(ctx context.Context, tenant, socialId string) error
	GetAllForEntities(ctx context.Context, tenant string, linkedEntityType model.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error)
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
	socialEntities := make(neo4jentity.SocialEntities, 0)
	for _, v := range socials {
		socialEntity := neo4jmapper.MapDbNodeToSocialEntity(v.Node)
		socialEntity.DataloaderKey = v.LinkedNodeId
		socialEntities = append(socialEntities, *socialEntity)
	}
	return &socialEntities, nil
}

func (s *socialService) Update(ctx context.Context, socialEntity neo4jentity.SocialEntity) (*neo4jentity.SocialEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	// get current social entity
	socialDbNode, err := s.services.Neo4jRepositories.SocialReadRepository.GetById(ctx, tenant, socialEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	currentSocialEntity := neo4jmapper.MapDbNodeToSocialEntity(socialDbNode)
	if currentSocialEntity.IsLinkedin() {
		if currentSocialEntity.Alias != "" || currentSocialEntity.ExternalId != "" {
			return currentSocialEntity, coserrors.ErrOperationNotAllowed
		}
	}

	// update social in DB
	updatedSocialNode, err := s.services.Neo4jRepositories.SocialWriteRepository.Update(ctx, tenant, socialEntity)
	if err != nil {
		return nil, err
	}

	err = s.services.RabbitMQService.PublishEvent(ctx, socialEntity.Id, model.SOCIAL, dto.UpdateSocial{
		Url: socialEntity.Url,
	})

	// get linked entities
	linkedEntities, err := s.services.Neo4jRepositories.CommonReadRepository.GetDbNodesLinkedTo(ctx, tenant, socialEntity.Id, model.SOCIAL.Neo4jLabel(), "HAS")
	if err != nil {
		tracing.TraceErr(span, err)
	}
	// notify linked entities updated.
	for _, linkedEntity := range linkedEntities {
		labels := linkedEntity.Labels
		props := utils.GetPropsFromNode(*linkedEntity)
		id := utils.GetStringPropOrEmpty(props, "id")

		if utils.Contains(labels, model.CONTACT.Neo4jLabel()) {
			err = s.services.RabbitMQService.PublishEvent(ctx, id, model.CONTACT, dto.UpdateSocialForContact{
				SocialId:  socialEntity.Id,
				SocialUrl: socialEntity.Url,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateSocialForContact"))
			}
			utils.EventCompleted(ctx, tenant, model.CONTACT.String(), id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		} else if utils.Contains(labels, model.ORGANIZATION.Neo4jLabel()) {
			err = s.services.RabbitMQService.PublishEvent(ctx, id, model.ORGANIZATION, dto.UpdateSocialForOrganization{
				SocialId:  socialEntity.Id,
				SocialUrl: socialEntity.Url,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateSocialForOrganization"))
			}
			utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		}
	}

	return neo4jmapper.MapDbNodeToSocialEntity(updatedSocialNode), nil
}

func (s *socialService) PermanentlyDelete(ctx context.Context, tenant string, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.PermanentlyDelete")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, socialId)

	// get linked entities
	// TODO get linked entities to send update events to rabbit and eventstore

	err := s.services.Neo4jRepositories.SocialWriteRepository.PermanentlyDelete(ctx, tenant, socialId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to permanently delete social"))
		return err
	}

	err = s.services.RabbitMQService.PublishEvent(ctx, socialId, model.SOCIAL, dto.Delete{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message DeleteSocial"))
	}

	return err
}

func (s *socialService) AddSocialToEntity(ctx context.Context, linkWith LinkWith, socialEntity neo4jentity.SocialEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.AddSocialToEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("linkWith.id", linkWith.Id), log.String("linkWith.type", string(linkWith.Type)))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate social url
	if socialEntity.Url == "" {
		err = errors.New("social url is required")
		tracing.TraceErr(span, err)
		return "", err
	}

	// validate linked entity exists
	exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
		return "", err
	}
	if !exists {
		err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
		tracing.TraceErr(span, err)
		return "", err
	}

	// prepare social url
	socialUrl := normalizeSocialUrl(socialEntity.Url)
	span.LogFields(log.String("socialUrl.normalized", socialUrl))

	// get or generate social entity id
	createSocialFlow := false
	socialId := socialEntity.Id
	if socialId == "" {
		createSocialFlow = true
		socialId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelSocial)
		if err != nil {
			return "", err
		}
	}
	tracing.TagEntity(span, socialId)

	// save social to neo4j
	data := neo4jrepository.SocialFields{
		SocialId:       socialId,
		Url:            socialUrl,
		Alias:          socialEntity.Alias,
		ExternalId:     socialEntity.ExternalId,
		FollowersCount: socialEntity.FollowersCount,
		CreatedAt:      utils.NowIfZero(socialEntity.CreatedAt),
		SourceFields: neo4jmodel.SourceFields{
			Source:    neo4jmodel.GetSource(socialEntity.Source.String()),
			AppSource: neo4jmodel.GetAppSource(socialEntity.AppSource),
		},
	}
	err = s.services.Neo4jRepositories.SocialWriteRepository.MergeSocialForEntity(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel(), data)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createSocialFlow {
		err = s.services.RabbitMQService.PublishEvent(ctx, socialId, model.SOCIAL, dto.CreateSocial{
			Url:           socialUrl,
			Alias:         socialEntity.Alias,
			ExtId:         socialEntity.ExternalId,
			FollowerCount: socialEntity.FollowersCount,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateSocial"))
		}
	}

	switch linkWith.Type {
	case model.CONTACT:
		err = s.services.RabbitMQService.PublishEvent(ctx, linkWith.Id, model.CONTACT, dto.AddSocialToContact{
			SocialId: socialId,
			Social:   socialUrl,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddSocialToContact"))
		}
		utils.EventCompleted(ctx, tenant, model.CONTACT.String(), linkWith.Id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	case model.ORGANIZATION:
		err = s.services.RabbitMQService.PublishEvent(ctx, linkWith.Id, model.ORGANIZATION, dto.AddSocialToOrganization{
			SocialId: socialId,
			Social:   socialUrl,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddSocialToOrganization"))
		}
		utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), linkWith.Id, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
	}

	return socialId, nil
}

func normalizeSocialUrl(url string) string {
	socialUrl := strings.TrimSpace(url)
	// adjust social url value
	if strings.HasPrefix(socialUrl, "linkedin.com") {
		socialUrl = "https://www." + socialUrl
	}
	if strings.Contains(socialUrl, "linkedin.com") && !strings.HasSuffix(socialUrl, "") {
		socialUrl = socialUrl + "/"
	}
	return socialUrl
}

func (s *socialService) GetById(ctx context.Context, socialId string) (*neo4jentity.SocialEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, socialId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)

	socialNode, err := s.services.Neo4jRepositories.SocialReadRepository.GetById(ctx, tenant, socialId)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToSocialEntity(socialNode), nil
}
