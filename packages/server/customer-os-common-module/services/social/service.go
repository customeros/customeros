package social

import (
	"context"
	"fmt"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type socialService struct {
	log     logger.Logger
	neo4j   *neoRepo.Repositories
	events  *events.EventsService
	contact interfaces.ContactService
}

func NewSocialService(log logger.Logger, neo4j *neoRepo.Repositories, events *events.EventsService, contact interfaces.ContactService) interfaces.SocialService {
	return &socialService{
		log:     log,
		neo4j:   neo4j,
		events:  events,
		contact: contact,
	}
}

func (s *socialService) SetContactService(contact interfaces.ContactService) {
	s.contact = contact
}

func (s *socialService) IsInitialized() bool {
	if s.neo4j == nil || s.events == nil || s.contact == nil {
		return false
	}
	return true
}

func (s *socialService) GetAllForEntities(ctx context.Context, tenant string, linkedEntityType model.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetAllForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("linkedEntityType", string(linkedEntityType)), log.Object("linkedEntityIds", linkedEntityIds))

	socials, err := s.neo4j.SocialReadRepository.GetAllForEntities(ctx, tenant, linkedEntityType, linkedEntityIds)
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

func (s *socialService) GetAllLinkedinForEntities(ctx context.Context, tenant string, linkedEntityType model.EntityType, linkedEntityIds []string) (*neo4jentity.SocialEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.GetAllLinkedinForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(log.String("linkedEntityType", string(linkedEntityType)), log.Object("linkedEntityIds", linkedEntityIds))

	socials, err := s.neo4j.SocialReadRepository.GetAllLinkedinForEntities(ctx, tenant, linkedEntityType, linkedEntityIds)
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
	socialDbNode, err := s.neo4j.SocialReadRepository.GetById(ctx, tenant, socialEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// reset alias and external id if linkedin url changed
	var alias, externalId *string
	currentSocialEntity := neo4jmapper.MapDbNodeToSocialEntity(socialDbNode)
	if currentSocialEntity.IsLinkedin() {
		if currentSocialEntity.Url != socialEntity.Url {
			alias = utils.StringPtr("")
			externalId = utils.StringPtr("")
		}
	}

	// update social in DB
	updatedSocialNode, err := s.neo4j.SocialWriteRepository.Update(ctx, tenant, socialEntity.Id, socialEntity.Url, alias, externalId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to update social"))
		return nil, err
	}

	updateSocialDto := dto.UpdateSocial{
		Url: socialEntity.Url,
	}
	if alias != nil {
		updateSocialDto.Alias = *alias
	}
	if externalId != nil {
		updateSocialDto.ExternalId = *externalId
	}
	err = s.events.Publisher.PublishEvent(ctx, socialEntity.Id, model.SOCIAL, updateSocialDto)

	// get linked entities
	linkedEntities, err := s.neo4j.CommonReadRepository.GetDbNodesLinkedTo(ctx, tenant, socialEntity.Id, model.SOCIAL.Neo4jLabel(), "HAS")
	if err != nil {
		tracing.TraceErr(span, err)
	}
	// notify linked entities updated.
	for _, linkedEntity := range linkedEntities {
		labels := linkedEntity.Labels
		props := utils.GetPropsFromNode(*linkedEntity)
		id := utils.GetStringPropOrEmpty(props, "id")

		if utils.Contains(labels, model.CONTACT.Neo4jLabel()) {
			err = s.events.Publisher.PublishEvent(ctx, id, model.CONTACT, dto.UpdateSocialForContact{
				SocialId:  socialEntity.Id,
				SocialUrl: socialEntity.Url,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateSocialForContact"))
			}
			if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
				s.events.Publisher.PublishEventCompleted(ctx, tenant, id, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
			}
		} else if utils.Contains(labels, model.ORGANIZATION.Neo4jLabel()) {
			err = s.events.Publisher.PublishEvent(ctx, id, model.ORGANIZATION, dto.UpdateSocialForOrganization{
				SocialId:  socialEntity.Id,
				SocialUrl: socialEntity.Url,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateSocialForOrganization"))
			}

			if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
				s.events.Publisher.PublishEventCompleted(ctx, tenant, id, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
			}
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

	err := s.neo4j.SocialWriteRepository.PermanentlyDelete(ctx, tenant, socialId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to permanently delete social"))
		return err
	}

	err = s.events.Publisher.PublishEvent(ctx, socialId, model.SOCIAL, dto.Delete{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message DeleteSocial"))
	}

	return err
}

func (s *socialService) AddSocialToEntity(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkWith common_srv.LinkWith, socialEntity neo4jentity.SocialEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.AddSocialToEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(
		log.String("linkWith.id", linkWith.Id),
		log.String("linkWith.type", string(linkWith.Type)),
		log.String("socialEntity.url", socialEntity.Url))

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

	socialId := ""

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// validate linked entity exists
		exists, err := s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, *txWithPostCommit.Tx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
			return nil, err
		}
		if !exists {
			err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
			tracing.TraceErr(span, err)
			return nil, err
		}

		// prepare social url
		socialUrl := normalizeSocialUrl(socialEntity.Url)
		span.LogFields(log.String("socialUrl.normalized", socialUrl))

		// Check social not used by another entity
		if socialEntity.IsLinkedin() {
			if linkWith.Type == model.ORGANIZATION {
				alias := socialEntity.Alias
				if alias == "" {
					// use identifier as alias
					alias = socialEntity.ExtractLinkedinCompanyIdentifierFromUrl()
				}
				orgsDbNodes, err := s.neo4j.OrganizationReadRepository.GetOrganizationsByLinkedIn(ctx, tenant, socialUrl, alias, socialEntity.ExternalId)
				if err != nil {
					tracing.TraceErr(span, err)
					return "", err
				}
				if len(orgsDbNodes) > 0 {
					if orgsDbNodes[0].Props["id"] == linkWith.Id {
						// social already linked to organization
						span.LogFields(log.Bool("result.alreadyLinked", true))
						return "", nil
					} else {
						span.LogFields(log.String("result.error", fmt.Sprintf("linkedin url %s already used by organization %s", socialUrl, orgsDbNodes[0].Props["id"])))
						err = coserrors.ErrLinkedInUsed
						return "", err
					}
				}
			} else if linkWith.Type == model.CONTACT {
				linkedInUsed, existingContactId, err := s.contact.CheckContactExistsWithLinkedIn(ctx, socialUrl, socialEntity.Alias, socialEntity.ExternalId)
				if err != nil {
					tracing.TraceErr(span, err)
					return "", err
				}
				if linkedInUsed {
					if existingContactId == linkWith.Id {
						// social already linked to contact
						span.LogFields(log.Bool("result.alreadyLinked", true))
						return "", nil
					} else {
						span.LogFields(log.String("result.error", fmt.Sprintf("linkedin url %s already used by contact %s", socialUrl, existingContactId)))
						err = coserrors.ErrLinkedInUsed
						return "", err
					}
				}
			}
		}

		// get or generate social entity id
		createSocialFlow := false
		socialId = socialEntity.Id
		if socialId == "" {
			createSocialFlow = true
			socialId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelSocial)
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
		err = s.neo4j.SocialWriteRepository.MergeSocialForEntity(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel(), data)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}

		// reset contact enrich attempts
		if socialEntity.IsLinkedin() {
			switch linkWith.Type {
			case model.CONTACT:
				_ = s.neo4j.ContactWriteRepository.ResetEnrichAttempts(ctx, txWithPostCommit.Tx, tenant, linkWith.Id)
			case model.ORGANIZATION:
				_ = s.neo4j.OrganizationWriteRepository.ResetEnrichAttempts(ctx, txWithPostCommit.Tx, tenant, linkWith.Id)
			}
		}

		// send events for social entity
		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createSocialFlow {
				err = s.events.Publisher.PublishEvent(ctx, socialId, model.SOCIAL, dto.CreateSocial{
					Url:           socialUrl,
					Alias:         socialEntity.Alias,
					ExtId:         socialEntity.ExternalId,
					FollowerCount: socialEntity.FollowersCount,
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateSocial"))
				}
			}
			return nil
		})

		// send events for linked entity
		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			switch linkWith.Type {
			case model.CONTACT:
				err = s.events.Publisher.PublishEvent(ctx, linkWith.Id, model.CONTACT, dto.AddSocialToContact{
					SocialId: socialId,
					Social:   socialUrl,
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddSocialToContact"))
				}
				s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
			case model.ORGANIZATION:
				err = s.events.Publisher.PublishEvent(ctx, linkWith.Id, model.ORGANIZATION, dto.AddSocialToOrganization{
					SocialId: socialId,
					Social:   socialUrl,
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message AddSocialToOrganization"))
				}
				s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
			}
			return nil
		})

		return nil, nil
	})

	return socialId, err
}

func (s *socialService) RemoveSocialFromEntity(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, linkWith common_srv.LinkWith, socialId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SocialService.RemoveSocialFromEntity")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(
		log.String("linkWith.id", linkWith.Id),
		log.String("linkWith.type", string(linkWith.Type)),
		log.String("socialEntity.id", socialId))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	// get social entity
	socialEntity, err := s.GetById(ctx, socialId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get social entity"))
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// validate linked entity exists
		exists, err := s.neo4j.CommonReadRepository.ExistsByIdInTx(ctx, *txWithPostCommit.Tx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel())
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to check linked entity exists"))
			return nil, err
		}
		if !exists {
			err = errors.Errorf("linked entity %s with id %s not found", linkWith.Type.String(), linkWith.Id)
			tracing.TraceErr(span, err)
			return nil, err
		}

		// neo query to remove social from entity
		err = s.neo4j.SocialWriteRepository.RemoveSocialForEntityById(ctx, tenant, linkWith.Id, linkWith.Type.Neo4jLabel(), socialId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to remove social from entity"))
			return nil, err
		}

		// send events for linked entity
		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			switch linkWith.Type {
			case model.CONTACT:
				err = s.events.Publisher.PublishEvent(ctx, linkWith.Id, model.CONTACT, dto.RemoveSocialFromContact{
					SocialId: socialId,
					Social:   socialEntity.Url,
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveSocialFromContact"))
				}
				s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, model.CONTACT, utils.NewEventCompletedDetails().WithUpdate())
			case model.ORGANIZATION:
				err = s.events.Publisher.PublishEvent(ctx, linkWith.Id, model.ORGANIZATION, dto.RemoveSocialFromOrganization{
					SocialId: socialId,
					Social:   socialEntity.Url,
				})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message RemoveSocialFromOrganization"))
				}
				s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
			}
			return nil
		})

		return nil, nil
	})

	return err
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

	socialNode, err := s.neo4j.SocialReadRepository.GetById(ctx, tenant, socialId)
	if err != nil {
		return nil, err
	}
	return neo4jmapper.MapDbNodeToSocialEntity(socialNode), nil
}
