package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
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
)

type OrganizationService interface {
	GetById(ctx context.Context, tenant, organizationId string) (*neo4jentity.OrganizationEntity, error)

	Save(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId *string, input *repository.OrganizationSaveFields) (*string, error)
	AddDomainFromWebsite(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId string, website string) (string, error)

	Hide(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error
	Show(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error

	Archive(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error

	GetLatestOrganizationsWithJobRolesForContacts(ctx context.Context, contactIds []string) (*neo4jentity.OrganizationWithJobRoleEntities, error)
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

	// if the org is new, we are looking for existing orgs with the same domain based on the website, we show it and we return it
	if organizationId == nil {
		domains := input.Domains
		if input.UpdateWebsite && input.Website != "" {
			websiteDomain := s.services.DomainService.ExtractDomainFromOrganizationWebsite(ctx, input.Website)
			if websiteDomain != "" {
				domains = append(domains, websiteDomain)
			}
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
	}

	if organizationId != nil {
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
	if organizationId == nil || (existing != nil && existing.CustomerOsId == "") {
		customerOsId, err := s.generateCustomerOSId(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		input.CustomerOsId = customerOsId
		input.UpdateCustomerOsId = true
	}

	if organizationId == nil {
		//if no name is provided, we try to extract if from domain
		if utils.IfNotNilString(input.Name) == "" && len(input.Domains) > 0 {
			input.Name = utils.CapitalizeAllParts(utils.GetDomainPrefix(input.Domains[0]), []string{"-", "_"})
			input.UpdateName = true
		}

		input.SourceFields.Source = constants.SourceOpenline

		if !input.UpdateHide {
			input.Hide = false
			input.UpdateHide = true
		}

		generatedId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, commonModel.NodeLabelOrganization)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		organizationId = &generatedId
	}

	domainForEnrich := ""

	_, err = utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (any, error) {

		err = s.services.Neo4jRepositories.OrganizationWriteRepository.Save(ctx, &tx, tenant, *organizationId, *input)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if existing == nil {
			_, err = s.services.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, &tx, tenant, *organizationId, commonModel.ORGANIZATION, neo4jenum.ActionCreated, "", "", utils.Now(), input.SourceFields.AppSource)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		if input.UpdateWebsite && input.Website != "" {
			domain, err := s.AddDomainFromWebsite(ctx, &tx, tenant, *organizationId, input.Website)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			if domain != "" && domainForEnrich == "" {
				domainForEnrich = domain
			}
		}

		if input.Domains != nil && len(input.Domains) > 0 {
			for _, domain := range input.Domains {
				err = s.services.Neo4jRepositories.OrganizationWriteRepository.LinkWithDomain(ctx, &tx, tenant, *organizationId, domain)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
				if domainForEnrich == "" {
					domainForEnrich = domain
				}
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
			err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, tx, tenant, *organizationId, commonModel.NodeLabelOrganization, externalSystemData)
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
	// completion event
	if input.SourceFields.AppSource != constants.AppSourceCustomerOsApi {
		details := utils.NewEventCompletedDetails()

		if existing == nil {
			details.WithCreate()
		} else {
			details.WithUpdate()
		}

		utils.EventCompleted(ctx, tenant, commonModel.ORGANIZATION.String(), *organizationId, s.services.GrpcClients, details)
	}

	// request organization enriching by domain event
	if domainForEnrich != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.EnrichOrganization(ctx, &organizationpb.EnrichOrganizationGrpcRequest{
				Tenant:         tenant,
				OrganizationId: *organizationId,
				Url:            domainForEnrich,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      input.SourceFields.AppSource,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	// request last touchpoint refresh for new organizations
	if existing == nil {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
				Tenant:         tenant,
				OrganizationId: *organizationId,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				AppSource:      input.SourceFields.AppSource,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
		}
	}

	return organizationId, nil
}

func (s *organizationService) AddDomainFromWebsite(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, organizationId string, website string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.AddDomainFromWebsite")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, organizationId)
	span.LogFields(log.String("website", website))

	domain := s.services.DomainService.ExtractDomainFromOrganizationWebsite(ctx, website)
	if domain == "" {
		return "", nil
	}

	providers := s.services.Cache.GetPersonalEmailProviders()
	if providers == nil || len(providers) == 0 {
		err := fmt.Errorf("personal email providers not loaded")
		tracing.TraceErr(span, err)
		return "", err
	}

	if s.services.Cache.IsPersonalEmailProvider(domain) {
		span.LogFields(log.String("result", "personal email provider"))
		return "", nil
	}

	err := s.services.Neo4jRepositories.OrganizationWriteRepository.LinkWithDomain(ctx, tx, tenant, organizationId, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return domain, nil
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

	utils.EventCompleted(ctx, tenant, commonModel.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithDelete())

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

	utils.EventCompleted(ctx, tenant, commonModel.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())

	return nil
}

func (s *organizationService) Archive(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Archive")
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

	err = s.services.Neo4jRepositories.OrganizationWriteRepository.Archive(ctx, tx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	utils.EventCompleted(ctx, tenant, commonModel.ORGANIZATION.String(), organizationId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithDelete())

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
