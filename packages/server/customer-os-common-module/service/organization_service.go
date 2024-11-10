package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type OrganizationService interface {
	GetById(ctx context.Context, tenant, organizationId string) (*neo4jentity.OrganizationEntity, error)

	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId *string, input *repository.OrganizationSaveFields) (*string, error)
	LinkWithDomain(ctx context.Context, tx *neo4j.ManagedTransaction, organizationId, domain string) error

	Hide(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error
	Show(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error

	GetLatestOrganizationsWithJobRolesForContacts(ctx context.Context, contactIds []string) (*neo4jentity.OrganizationWithJobRoleEntities, error)

	GetHiddenOrganizationIds(ctx context.Context, hiddenAfter time.Time) ([]string, error)
	GetMergedOrganizationIds(ctx context.Context, mergedAfter time.Time) ([]string, error)
	RequestRefreshLastTouchpoint(ctx context.Context, organizationId string) error
	RefreshLastTouchpoint(ctx context.Context, organizationId string) error
}

type organizationService struct {
	services *Services
}

func NewOrganizationService(services *Services) OrganizationService {
	return &organizationService{
		services: services,
	}
}

func (s *organizationService) GetById(ctx context.Context, tenant, organizationId string) (*neo4jentity.OrganizationEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))

	dbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return neo4jmapper.MapDbNodeToOrganizationEntity(dbNode), nil
}

func (s *organizationService) Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId *string, input *repository.OrganizationSaveFields) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("input", input))

	var err error
	var existing *neo4jentity.OrganizationEntity
	createFlow := false

	// prepare primary domain from website
	primaryDomainFromWebsite := ""
	adjustedWebsite := input.Website
	if input.UpdateWebsite && input.Website != "" {
		primaryDomainFromWebsite, adjustedWebsite = s.services.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, input.Website)
		span.LogFields(log.String("primaryDomainFromWebsite", primaryDomainFromWebsite))
		span.LogFields(log.String("adjustedWebsite", adjustedWebsite))
	}

	// prepare domains in advance
	err = s.services.DomainService.MergeDomain(ctx, primaryDomainFromWebsite)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to merge domain"))
	}
	for _, domain := range input.Domains {
		err = s.services.DomainService.MergeDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to merge domain"))
		}
	}

	// if the org is new, we are looking for existing orgs with the same domain based on the website, we show it and we return it
	if organizationId == nil {
		createFlow = true
		domains := input.Domains
		if input.UpdateWebsite && input.Website != "" && primaryDomainFromWebsite != "" {
			domains = append(domains, primaryDomainFromWebsite)
		}
		domains = utils.RemoveEmpties(domains)
		domains = utils.RemoveDuplicates(domains)

		if len(domains) > 0 {
			// for each domain check that no org exists with that domain
			// if exist reject creation and return error
			for _, domain := range domains {
				orgDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, domain)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
				if orgDbNode != nil {
					organizationEntity := neo4jmapper.MapDbNodeToOrganizationEntity(orgDbNode)
					if organizationEntity.Hide {
						err = s.Show(ctx, tx, tenant, organizationEntity.ID)
						if err != nil {
							tracing.TraceErr(span, err)
							return nil, nil
						}
					}
					return &organizationEntity.ID, nil
				}
			}
		}

		// set source fields if not provided
		if input.SourceFields.AppSource == "" {
			input.SourceFields.AppSource = common.GetAppSourceFromContext(ctx)
		}
		if input.SourceFields.Source == "" {
			if input.ExternalSystem.Available() {
				input.SourceFields.Source = input.ExternalSystem.ExternalSystemId
			} else {
				input.SourceFields.Source = neo4jentity.DataSourceOpenline.String()
			}
		}
	} else {
		existing, err = s.GetById(ctx, tenant, *organizationId)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	//validate stage and relationship combination all the time ( from input or existing computed )
	stage := input.Stage.String()
	relationship := input.Relationship.String()
	if stage == "" && existing != nil && existing.Stage != "" {
		stage = existing.Stage.String()
	}
	if relationship == "" && existing != nil && existing.Relationship != "" {
		relationship = existing.Relationship.String()
	}
	if !neo4jentity.OrganizationStageAndRelationshipCompatible(stage, relationship) {
		err := errors.New("Stage and Relationship are not compatible")
		tracing.TraceErr(span, err)
		return nil, err
	}

	//generate customerOsId if not provided or if it is empty in the db
	if createFlow || (existing != nil && existing.CustomerOsId == "") {
		customerOsId, err := s.generateCustomerOSId(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		input.CustomerOsId = customerOsId
		input.UpdateCustomerOsId = true
	}

	if createFlow {
		//if no name is provided, we try to extract if from domain
		if utils.IfNotNilString(input.Name) == "" {
			domain := primaryDomainFromWebsite
			if domain == "" && len(input.Domains) > 0 {
				domain = input.Domains[0]
			}
			if domain != "" {
				input.Name = utils.CapitalizeAllParts(utils.GetDomainWithoutTLD(domain), []string{"-", "_", "."})
				input.UpdateName = true
			}
			// if still empty, extract from website
			if input.Name == "" && input.Website != "" {
				websiteDomain := utils.ExtractDomain(input.Website)
				input.Name = utils.CapitalizeAllParts(utils.GetDomainWithoutTLD(websiteDomain), []string{"-", "_", "."})
				input.UpdateName = true
			}
		}

		input.SourceFields.Source = constants.SourceOpenline

		if !input.UpdateHide {
			input.Hide = false
			input.UpdateHide = true
		}

		generatedId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelOrganization)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		organizationId = &generatedId
	}

	newDomains := make([]string, 0)

	_, err = utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		err = s.services.Neo4jRepositories.OrganizationWriteRepository.Save(ctx, &tx, tenant, *organizationId, *input)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if existing == nil {
			_, err = s.services.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, &tx, tenant, *organizationId, model.ORGANIZATION, neo4jenum.ActionCreated, "", "", utils.Now(), input.SourceFields.AppSource)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		if input.UpdateWebsite && adjustedWebsite != "" {
			input.Website = adjustedWebsite
			if primaryDomainFromWebsite != "" {
				newDomains = append(newDomains, primaryDomainFromWebsite)
				err = s.LinkWithDomain(ctx, &tx, *organizationId, primaryDomainFromWebsite)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
			}
		}

		if input.Domains != nil && len(input.Domains) > 0 {
			for _, domain := range input.Domains {
				err = s.LinkWithDomain(ctx, &tx, *organizationId, domain)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
				newDomains = append(newDomains, domain)
			}
		}

		if input.ExternalSystem.Available() {
			externalSystemData := neo4jmodel.ExternalSystem{
				ExternalSystemId: input.ExternalSystem.ExternalSystemId,
				ExternalUrl:      input.ExternalSystem.ExternalUrl,
				ExternalId:       input.ExternalSystem.ExternalId,
				ExternalIdSecond: input.ExternalSystem.ExternalIdSecond,
				ExternalSource:   input.ExternalSystem.ExternalSource,
				SyncDate:         input.ExternalSystem.SyncDate,
			}
			err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, tenant, *organizationId, model.NodeLabelOrganization, externalSystemData)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		if input.UpdateOwnerId {
			err = s.services.Neo4jRepositories.OrganizationWriteRepository.ReplaceOwner(ctx, &tx, tenant, *organizationId, input.OwnerId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// create events section
	if createFlow {
		err = s.services.RabbitMQService.PublishEvent(ctx, *organizationId, model.ORGANIZATION, dto.New_CreateOrganization_From_OrganizationFields(*input))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateOrganization"))
		}
		utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), *organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())
	} else {
		err = s.services.RabbitMQService.PublishEvent(ctx, *organizationId, model.ORGANIZATION, dto.New_UpdateOrganization_From_OrganizationFields(*input))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateOrganization"))
		}
		if input.SourceFields.AppSource != constants.AppSourceCustomerOsApi {
			utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), *organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		}
	}

	// select primary domain from new domains
	primaryDomain := primaryDomainFromWebsite
	if primaryDomain == "" && len(newDomains) > 0 {
		newDomains = utils.RemoveEmpties(newDomains)
		newDomains = utils.RemoveDuplicates(newDomains)
		for _, domain := range newDomains {
			domainEntity, err := s.services.DomainService.GetDomain(ctx, domain)
			if err != nil {
				tracing.TraceErr(span, err)
			} else if domainEntity != nil && domainEntity.IsPrimary != nil && *domainEntity.IsPrimary {
				primaryDomain = domain
				break
			}
		}
	}

	if primaryDomain != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.EnrichOrganization(ctx, &organizationpb.EnrichOrganizationGrpcRequest{
				Tenant:         tenant,
				OrganizationId: *organizationId,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      input.SourceFields.AppSource,
				Url:            primaryDomain,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// request last touchpoint refresh for new organizations
	if createFlow {
		err = s.RequestRefreshLastTouchpoint(ctx, *organizationId)
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return organizationId, nil
}

func (s *organizationService) Hide(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Hide")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	organization, err := s.GetById(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if organization == nil {
		err = fmt.Errorf("opportunity not found")
		tracing.TraceErr(span, err)
		return err
	}

	fields := repository.OrganizationSaveFields{Hide: true, UpdateHide: true}
	err = s.services.Neo4jRepositories.OrganizationWriteRepository.Save(ctx, tx, tenant, organizationId, fields)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithDelete())

	return nil
}

func (s *organizationService) Show(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Show")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	organization, err := s.GetById(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if organization == nil {
		err = fmt.Errorf("opportunity not found")
		tracing.TraceErr(span, err)
		return err
	}

	fields := repository.OrganizationSaveFields{Hide: false, UpdateHide: true}
	err = s.services.Neo4jRepositories.OrganizationWriteRepository.Save(ctx, tx, tenant, organizationId, fields)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())

	return nil
}

func (s *organizationService) generateCustomerOSId(ctx context.Context, tenant string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.generateCustomerOSId")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	var customerOsId string
	maxAttempts := 20
	for attempt := 1; attempt < maxAttempts+1; attempt++ {
		customerOsId = generateNewRandomCustomerOsId()

		exists, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByCustomerOsId(ctx, tenant, customerOsId)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}

		if exists == nil {
			break
		}
	}

	return customerOsId, nil
}

func generateNewRandomCustomerOsId() string {
	charset := "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	customerOsID := "C-" + utils.GenerateRandomStringFromCharset(3, charset) + "-" + utils.GenerateRandomStringFromCharset(3, charset)
	return customerOsID
}

func (s *organizationService) GetLatestOrganizationsWithJobRolesForContacts(ctx context.Context, contactIds []string) (*neo4jentity.OrganizationWithJobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetLatestOrganizationsWithJobRolesForContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactIds", contactIds))

	dbResults, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetLatestOrganizationWithJobRoleForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	orgWithJobRoleEntities := make(neo4jentity.OrganizationWithJobRoleEntities, 0)
	for _, v := range dbResults {
		orgWithJobRoleEntity := neo4jentity.OrganizationWithJobRole{}
		orgWithJobRoleEntity.Organization = *neo4jmapper.MapDbNodeToOrganizationEntity(v.Pair.First)
		orgWithJobRoleEntity.JobRole = *neo4jmapper.MapDbNodeToJobRoleEntity(v.Pair.Second)
		orgWithJobRoleEntity.DataloaderKey = v.LinkedNodeId
		orgWithJobRoleEntities = append(orgWithJobRoleEntities, orgWithJobRoleEntity)
	}
	return &orgWithJobRoleEntities, nil
}

func (s *organizationService) LinkWithDomain(ctx context.Context, tx *neo4j.ManagedTransaction, organizationId, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.LinkWithDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, organizationId)
	span.LogKV("domain", domain)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if !s.services.DomainService.AcceptedDomainForOrganization(ctx, domain) {
		return nil
	}

	domainLinkedToOrg, err := s.services.Neo4jRepositories.OrganizationWriteRepository.LinkWithDomain(ctx, tx, tenant, organizationId, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to link domain in neo4j"))
		return err
	}

	// execute only if not in transaction and domain was linked with org
	if tx == nil && domainLinkedToOrg {
		// send event to rabbitmq
		err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.NewAddDomainEvent(domain))
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to publish event AddDomain"))
		}

		// send event to events platform
		utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

		// send organization enrich request
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.EnrichOrganization(ctx, &organizationpb.EnrichOrganizationGrpcRequest{
				Tenant:         tenant,
				OrganizationId: organizationId,
				Url:            domain,
				AppSource:      common.GetAppSourceFromContext(ctx),
			})
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to request enrich organization"))
		}
	}

	return nil
}

func (s *organizationService) GetHiddenOrganizationIds(ctx context.Context, hiddenAfter time.Time) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetHiddenOrganizationIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("hiddenAfter", hiddenAfter))

	organizationIds, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetHiddenOrganizationIds(ctx, common.GetTenantFromContext(ctx), hiddenAfter)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return organizationIds, nil
}

func (s *organizationService) GetMergedOrganizationIds(ctx context.Context, mergedAfter time.Time) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetMergedOrganizationIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("mergedAfter", mergedAfter))

	organizationIds, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetMergedOrganizationIds(ctx, common.GetTenantFromContext(ctx), mergedAfter)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return organizationIds, nil
}

func (s *organizationService) RequestRefreshLastTouchpoint(ctx context.Context, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.RequestRefreshLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, organizationId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.RequestRefreshLastTouchpoint{})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to publish event RequestRefreshLastTouchpoint"))
		return err
	}

	return nil
}

func (s *organizationService) RefreshLastTouchpoint(ctx context.Context, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.RefreshLastTouchpoint")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, organizationId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	//fetch the real touchpoint
	//if it doesn't exist, check for the Created Action
	var lastTouchpointId string
	var lastTouchpointAt *time.Time
	var timelineEventNode *dbtype.Node

	// get current last touchpoint
	organizationEntity, err := s.GetById(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	currentLastTouchpointId := utils.IfNotNilString(organizationEntity.LastTouchpointId)
	span.LogFields(log.String("currentLastTouchpointId", currentLastTouchpointId))

	lastTouchpointAt, lastTouchpointId, err = s.services.Neo4jRepositories.TimelineEventReadRepository.CalculateAndGetLastTouchPoint(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.services.Logger.Errorf("Failed to calculate last touchpoint: %v", err.Error())
		span.LogFields(log.Bool("last touchpoint failed", true))
		return nil
	}

	if lastTouchpointAt == nil {
		timelineEventNode, err = s.services.Neo4jRepositories.ActionReadRepository.GetLastAction(ctx, tenant, organizationId, model.ORGANIZATION, neo4jenum.ActionCreated)
		if err != nil {
			tracing.TraceErr(span, err)
			s.services.Logger.Errorf("Failed to get created action: %v", err.Error())
			return nil
		}
		if timelineEventNode != nil {
			propsFromNode := utils.GetPropsFromNode(*timelineEventNode)
			lastTouchpointId = utils.GetStringPropOrEmpty(propsFromNode, "id")
			lastTouchpointAt = utils.GetTimePropOrNil(propsFromNode, "createdAt")
		}
	} else {
		timelineEventNode, err = s.services.Neo4jRepositories.TimelineEventReadRepository.GetTimelineEvent(ctx, tenant, lastTouchpointId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.services.Logger.Errorf("Failed to get last touchpoint: %v", err.Error())
			return nil
		}
	}

	if timelineEventNode == nil {
		s.services.Logger.Infof("Last touchpoint not available for organization: %s", organizationId)
		span.LogFields(log.Bool("last touchpoint not found", true))
		return nil
	}

	span.LogFields(log.String("lastTouchpointId", lastTouchpointId))
	// last touchpoint not changed, skip
	if lastTouchpointId == currentLastTouchpointId {
		span.LogFields(log.Bool("last touchpoint changed", false))
		s.services.Logger.Infof("Last touchpoint not changed for organization: %s", organizationId)
		return nil
	} else {
		span.LogFields(log.Bool("last touchpoint changed", true))
	}

	timelineEvent := neo4jmapper.MapDbNodeToTimelineEvent(timelineEventNode)
	if timelineEvent == nil {
		s.services.Logger.Infof("Last touchpoint not available for organization: %s", organizationId)
		span.LogFields(log.Bool("last touchpoint not found", true))
		return nil
	}

	var timelineEventType string
	switch timelineEvent.TimelineEventLabel() {
	case model.NodeLabelPageView:
		timelineEventType = neo4jenum.TouchpointTypePageView.String()
	case model.NodeLabelInteractionSession:
		timelineEventType = neo4jenum.TouchpointTypeInteractionSession.String()
	case model.NodeLabelNote:
		timelineEventType = neo4jenum.TouchpointTypeNote.String()
	case model.NodeLabelInteractionEvent:
		timelineEventInteractionEvent := timelineEvent.(*neo4jentity.InteractionEventEntity)
		if timelineEventInteractionEvent.Channel == "EMAIL" {
			interactionEventSentByUser, err := s.services.Neo4jRepositories.InteractionEventReadRepository.InteractionEventSentByUser(ctx, tenant, timelineEventInteractionEvent.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				s.services.Logger.Errorf("Failed to check if interaction event was sent by user: %v", err.Error())
			}
			if interactionEventSentByUser {
				timelineEventType = neo4jenum.TouchpointTypeInteractionEventEmailSent.String()
			} else {
				timelineEventType = neo4jenum.TouchpointTypeInteractionEventEmailReceived.String()
			}
		} else if timelineEventInteractionEvent.Channel == "VOICE" {
			timelineEventType = neo4jenum.TouchpointTypeInteractionEventPhoneCall.String()
		} else if timelineEventInteractionEvent.Channel == "CHAT" {
			timelineEventType = neo4jenum.TouchpointTypeInteractionEventChat.String()
		} else if timelineEventInteractionEvent.EventType == "meeting" {
			timelineEventType = neo4jenum.TouchpointTypeMeeting.String()
		}
	case model.NodeLabelMeeting:
		timelineEventType = neo4jenum.TouchpointTypeMeeting.String()
	case model.NodeLabelAction:
		timelineEventAction := timelineEvent.(*neo4jentity.ActionEntity)
		if timelineEventAction.Type == neo4jenum.ActionCreated {
			timelineEventType = neo4jenum.TouchpointTypeActionCreated.String()
		} else {
			timelineEventType = neo4jenum.TouchpointTypeAction.String()
		}
	case model.NodeLabelLogEntry:
		timelineEventType = neo4jenum.TouchpointTypeAction.String()
	case model.NodeLabelIssue:
		timelineEventIssue := timelineEvent.(*neo4jentity.IssueEntity)
		if timelineEventIssue.CreatedAt.Equal(timelineEventIssue.UpdatedAt) {
			timelineEventType = neo4jenum.TouchpointTypeIssueCreated.String()
		} else {
			timelineEventType = neo4jenum.TouchpointTypeIssueUpdated.String()
		}
	default:
		s.services.Logger.Infof("Last touchpoint not available for organization: %s", organizationId)
	}

	if err = s.services.Neo4jRepositories.OrganizationWriteRepository.UpdateLastTouchpoint(ctx, tenant, organizationId, lastTouchpointAt, lastTouchpointId, timelineEventType); err != nil {
		tracing.TraceErr(span, err)
		s.services.Logger.Errorf("Failed to update last touchpoint for tenant %s, organization %s: %s", tenant, organizationId, err.Error())
		return err
	}

	utils.EventCompleted(ctx, tenant, model.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}
