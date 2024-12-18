package service

import (
	"context"
	"fmt"
	"time"

	mailsherpa "github.com/customeros/mailsherpa/mailvalidate"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type OrganizationService interface {
	GetById(ctx context.Context, tenant, organizationId string) (*neo4jentity.OrganizationEntity, error)

	CreateFromGlobalOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, globalOrgId uint64) (string, error)
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, dataFields data_fields.OrganizationFields) (string, error)
	LinkWithDomain(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId, domain string) error

	Hide(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string) error
	Show(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string) error

	GetHiddenOrganizationIds(ctx context.Context, hiddenAfter time.Time) ([]string, error)
	GetMergedOrganizationIds(ctx context.Context, mergedAfter time.Time) ([]string, error)
	RequestRefreshLastTouchpoint(ctx context.Context, organizationId string) error
	RefreshLastTouchpoint(ctx context.Context, organizationId string) error
	CheckOrganizationExistsWithEmail(ctx context.Context, email string) (bool, string, error)
	CheckOrganizationExistsWithLinkedIn(ctx context.Context, url, alias, externalId string) (bool, string, error)
	GetPrimaryOrganizationsWithJobRoleForContacts(ctx context.Context, contactIds []string) (*neo4jentity.OrganizationWithJobRoleEntities, error)
}

type organizationService struct {
	services *Services
}

func NewOrganizationService(services *Services) OrganizationService {
	return &organizationService{
		services: services,
	}
}

func (s *organizationService) CreateFromGlobalOrganization(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, globalOrgId uint64) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.CreateFromGlobalOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Uint64("globalOrgId", globalOrgId))

	globalOrganization, err := s.services.PostgresRepositories.GlobalOrganizationRepository.GetById(ctx, globalOrgId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if globalOrganization == nil {
		err = errors.New(fmt.Sprintf("Global organization with id %d not found", globalOrgId))
		tracing.TraceErr(span, err)
		return "", err
	}

	dataFields := data_fields.OrganizationFields{
		GlobalOrgId:      utils.ToPtr(globalOrgId),
		Name:             utils.StringPtr(globalOrganization.Name),
		PrimaryDomain:    utils.StringPtr(globalOrganization.PrimaryDomain),
		Description:      utils.StringPtr(globalOrganization.Description),
		Website:          utils.StringPtr(globalOrganization.Website),
		LogoUrl:          utils.StringPtr(globalOrganization.LogoUrl),
		IconUrl:          utils.StringPtr(globalOrganization.IconUrl),
		LinkedInUrl:      utils.StringPtr(globalOrganization.LinkedInUrl),
		LinkedInAlias:    utils.StringPtr(globalOrganization.LinkedInAlias),
		ValueProposition: utils.StringPtr(globalOrganization.ValueProposition),
		TargetAudience:   utils.StringPtr(globalOrganization.TargetAudience),
		Domains:          utils.StringToSlice(globalOrganization.OtherDomains),
	}

	if globalOrganization.YearFounded > 0 {
		dataFields.YearFounded = utils.Int64Ptr(int64(globalOrganization.YearFounded))
	}
	if globalOrganization.EmployeeCount > 0 {
		dataFields.Employees = utils.Int64Ptr(int64(globalOrganization.EmployeeCount))
	}

	return s.Save(ctx, txWithPostCommit, nil, dataFields)
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

func (s *organizationService) Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, input data_fields.OrganizationFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	var existingOrganizationEntity *neo4jentity.OrganizationEntity
	createFlow := false
	organizationId := ""

	// prepare primary domain from website
	primaryDomain := utils.IfNotNilString(input.PrimaryDomain)
	adjustedWebsite := utils.IfNotNilString(input.Website)
	if input.GlobalOrgId == nil {
		if utils.IfNotNilString(input.Website) != "" {
			primaryDomain, adjustedWebsite = s.services.DomainService.GetPrimaryDomainForOrganizationWebsite(ctx, *input.Website)
			span.LogFields(log.String("process.primaryDomainFromWebsite", primaryDomain))
			span.LogFields(log.String("process.adjustedWebsite", adjustedWebsite))
		}
	}

	if primaryDomain != "" && adjustedWebsite == "" {
		adjustedWebsite = primaryDomain
	}

	// prepare domains in advance
	err = s.services.DomainService.MergeDomain(ctx, primaryDomain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to merge domain"))
	}
	for _, domain := range input.Domains {
		err = s.services.DomainService.MergeDomain(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to merge domain"))
		}
	}

	if id == nil {
		createFlow = true
		span.LogFields(log.String("process.flow", "create"))
	} else {
		span.LogFields(log.String("process.flow", "update"))
	}

	if createFlow {
		domains := input.Domains
		if primaryDomain != "" {
			domains = append(domains, primaryDomain)
		}
		domains = utils.RemoveEmpties(domains)
		domains = utils.RemoveDuplicates(domains)

		// Dedup organizations by domain
		if len(domains) > 0 {
			// for each domain check that no org exists with that domain
			// if exist reject creation and return existing org id
			for _, domain := range domains {
				orgByDomainDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, domain)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "Error fetching organization by domain"))
					return "", err
				}
				// organization by domain found
				if orgByDomainDbNode != nil {
					organizationByDomainEntity := neo4jmapper.MapDbNodeToOrganizationEntity(orgByDomainDbNode)
					span.LogFields(log.String("result.duplicate.domain", domain))
					span.LogFields(log.String("result.duplicate.orgId", organizationByDomainEntity.ID))

					_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
						if organizationByDomainEntity.IsHidden() {
							err = s.Show(ctx, txWithPostCommit, organizationByDomainEntity.ID)
							if err != nil {
								tracing.TraceErr(span, err)
								return nil, nil
							}
						} else {
							// just update organization updatedAt
							err = s.services.Neo4jRepositories.CommonWriteRepository.TouchEntity(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelOrganization, organizationByDomainEntity.ID)
							if err != nil {
								tracing.TraceErr(span, err)
								s.services.Logger.Errorf("Failed to update organization updated at property: %v", err.Error())
							}
							txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
								s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
								return nil
							})
						}
						return nil, nil
					})

					return organizationByDomainEntity.ID, nil
				}
			}
		}

		// Dedup organization by linked in url
		if utils.IfNotNilString(input.LinkedInUrl) != "" {
			linkedInUrl := utils.IfNotNilString(input.LinkedInUrl)
			if (neo4jentity.SocialEntity{Url: linkedInUrl}).IsLinkedin() {
				linkedInAlreadyUsed, existingOrganizationId, err := s.CheckOrganizationExistsWithLinkedIn(ctx, linkedInUrl, utils.IfNotNilString(input.LinkedInAlias), "")
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to check organization exists with linkedIn"))
					return "", err
				}
				if linkedInAlreadyUsed {
					organizationByLinkedInEntity, err := s.GetById(ctx, tenant, existingOrganizationId)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "unable to get organization by id"))
						return "", err
					}

					_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
						if organizationByLinkedInEntity.IsHidden() {
							err = s.Show(ctx, txWithPostCommit, organizationByLinkedInEntity.ID)
							if err != nil {
								tracing.TraceErr(span, err)
								return nil, nil
							}
						} else {
							// just update organization updatedAt
							err = s.services.Neo4jRepositories.CommonWriteRepository.TouchEntity(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelOrganization, organizationByLinkedInEntity.ID)
							if err != nil {
								tracing.TraceErr(span, err)
								s.services.Logger.Errorf("Failed to update organization updated at property: %v", err.Error())
							}
							txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
								s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationByLinkedInEntity.ID, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
								return nil
							})
						}
						return nil, nil
					})

					return organizationByLinkedInEntity.ID, nil
				}
			}
		}

		generatedId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelOrganization)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		organizationId = generatedId

	} else {
		organizationId = *id
		existingOrganizationEntity, err = s.GetById(ctx, tenant, organizationId)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, organizationId)

	// validate stage and relationship combination all the time (from input or existing computed )
	stageStr := input.GetStageStr()
	relationshipStr := input.GetRelationshipStr()
	if stageStr != "" || relationshipStr != "" {
		if stageStr == "" && existingOrganizationEntity != nil && existingOrganizationEntity.Stage != "" {
			stageStr = existingOrganizationEntity.Stage.String()
		}
		if relationshipStr == "" && existingOrganizationEntity != nil && existingOrganizationEntity.Relationship != "" {
			relationshipStr = existingOrganizationEntity.Relationship.String()
		}
	}
	if !neo4jentity.OrganizationStageAndRelationshipCompatible(ctx, stageStr, relationshipStr) {
		err := errors.New("Stage and Relationship are not compatible")
		tracing.TraceErr(span, err)
		return "", err
	}

	// adapt fields for creating new organization
	if createFlow {
		if utils.IfNotNilString(input.AppSource) == "" {
			input.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if utils.IfNotNilString(input.Source) == "" {
			if input.ExternalSystemAvailable() {
				input.Source = utils.StringPtr(input.ExternalSystem.ExternalSystemId)
			} else {
				input.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
			}
		}
		// if no name is provided, we try to extract if from domain
		if utils.IfNotNilString(input.Name) == "" {
			domain := primaryDomain
			if domain == "" && len(input.Domains) > 0 {
				domain = input.Domains[0]
			}
			if domain != "" {
				input.Name = utils.StringPtr(utils.CapitalizeAllParts(utils.GetDomainWithoutTLD(domain), []string{"-", "_", "."}))
			}
			// if still empty, extract from website
			if utils.IfNotNilString(input.Name) == "" && utils.IfNotNilString(input.Website) != "" {
				websiteDomain := utils.ExtractDomain(utils.IfNotNilString(input.Website))
				input.Name = utils.StringPtr(utils.CapitalizeAllParts(utils.GetDomainWithoutTLD(websiteDomain), []string{"-", "_", "."}))
			}
		}
		input.Hide = utils.BoolPtr(false)
		if input.Stage == nil {
			input.Stage = utils.ToPtr(neo4jenum.Lead)
		}
		if input.Relationship == nil {
			input.Relationship = utils.ToPtr(neo4jenum.OrganizationRelationshipProspect)
		}
	}

	// generate customerOsId if not provided or if it is empty in the db
	if createFlow || (existingOrganizationEntity != nil && existingOrganizationEntity.CustomerOsId == "") {
		customerOsId, err := s.generateCustomerOSId(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		input.CustomerOsId = utils.StringPtr(customerOsId)
	}

	// Clean and update organization name if not updated manually
	if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
		if input.Name != nil {
			input.Name = utils.StringPtr(utils.CleanName(*input.Name))
		}
	}

	newDomains := make([]string, 0)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err = s.services.Neo4jRepositories.OrganizationWriteRepository.Save(ctx, txWithPostCommit.Tx, tenant, organizationId, input)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to save organization"))
			return nil, err
		}

		if createFlow {
			_, err = s.services.Neo4jRepositories.ActionWriteRepository.MergeByActionType(ctx, txWithPostCommit.Tx, tenant, organizationId, model.ORGANIZATION, neo4jenum.ActionCreated, "", "", utils.Now(), common.GetAppSourceFromContext(ctx))
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to merge action"))
				return nil, err
			}
			err = s.services.Neo4jRepositories.OrganizationWriteRepository.RefreshContactCountByOrgId(ctx, txWithPostCommit.Tx, tenant, organizationId)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to refresh contact count by organization id"))
			}
		}

		if utils.IfNotNilString(input.Website) != "" && adjustedWebsite != "" {
			input.Website = utils.StringPtr(adjustedWebsite)
			if primaryDomain != "" {
				newDomains = append(newDomains, primaryDomain)
				err = s.LinkWithDomain(ctx, txWithPostCommit, organizationId, primaryDomain)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to link with domain"))
					return nil, err
				}
			}
		}

		if input.Domains != nil && len(input.Domains) > 0 {
			for _, domain := range input.Domains {
				err = s.LinkWithDomain(ctx, txWithPostCommit, organizationId, domain)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
				newDomains = append(newDomains, domain)
			}
		}

		if input.ExternalSystemAvailable() {
			externalSystemData := neo4jmodel.ExternalSystem{
				ExternalSystemId: input.ExternalSystem.ExternalSystemId,
				ExternalUrl:      input.ExternalSystem.ExternalUrl,
				ExternalId:       input.ExternalSystem.ExternalId,
				ExternalIdSecond: input.ExternalSystem.ExternalIdSecond,
				ExternalSource:   input.ExternalSystem.ExternalSource,
				SyncDate:         input.ExternalSystem.SyncDate,
			}
			err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntityInTx(ctx, txWithPostCommit.Tx, tenant, organizationId, model.NodeLabelOrganization, externalSystemData)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		if utils.IfNotNilString(input.OwnerId) != "" {
			err = s.services.Neo4jRepositories.OrganizationWriteRepository.ReplaceOwner(ctx, txWithPostCommit.Tx, tenant, organizationId, *input.OwnerId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

		// if linked in provided, merge social with organization in create only mode
		if createFlow && utils.IfNotNilString(input.LinkedInUrl) != "" {
			linkedInUrl := utils.IfNotNilString(input.LinkedInUrl)
			if (neo4jentity.SocialEntity{Url: linkedInUrl}).IsLinkedin() {
				_, err := s.services.SocialService.AddSocialToEntity(ctx, txWithPostCommit,
					LinkWith{
						Id:   organizationId,
						Type: model.ORGANIZATION,
					},
					neo4jentity.SocialEntity{
						Url:       linkedInUrl,
						Alias:     utils.IfNotNilString(input.LinkedInAlias),
						Source:    neo4jentity.DecodeDataSource(utils.IfNotNilString(input.Source)),
						AppSource: utils.IfNotNilString(input.AppSource),
					})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to merge social with organization"))
				}
			}
		}

		// in update only mode if slack channel changed, update reported relationship
		if !createFlow && input.SlackChannelId != nil {
			if existingOrganizationEntity.SlackChannelId != utils.IfNotNilString(input.SlackChannelId) {
				if utils.IfNotNilString(input.SlackChannelId) == "" {
					err := s.services.Neo4jRepositories.IssueWriteRepository.RemoveReportedByOrganizationWithGroupId(ctx, txWithPostCommit.Tx, tenant, organizationId, existingOrganizationEntity.SlackChannelId)
					if err != nil {
						tracing.TraceErr(span, err)
					}
				} else {
					err := s.services.Neo4jRepositories.IssueWriteRepository.ReportedByOrganizationWithGroupId(ctx, txWithPostCommit.Tx, tenant, organizationId, utils.IfNotNilString(input.SlackChannelId))
					if err != nil {
						tracing.TraceErr(span, err)
					}
				}
			}
		}

		// add post commit actions to send events
		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if createFlow {
				err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.CreateOrganization{input})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateOrganization"))
				}
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithCreate())
			} else {
				err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.UpdateOrganization{input})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateOrganization"))
				}
				if common.GetAppSourceFromContext(ctx) != constants.AppSourceCustomerOsApi {
					s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())
				}
			}
			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// request enrich organization by primary domain if in creaate mode or organization not enriched yet
			if createFlow || existingOrganizationEntity.EnrichDetails.EnrichedAt == nil {
				// select primary domain from new domains
				localPrimaryDomain := primaryDomain
				if localPrimaryDomain == "" && len(newDomains) > 0 {
					newDomains = utils.RemoveEmpties(newDomains)
					newDomains = utils.RemoveDuplicates(newDomains)
					for _, domain := range newDomains {
						domainEntity, err := s.services.DomainService.GetDomain(ctx, domain)
						if err != nil {
							tracing.TraceErr(span, err)
						} else if domainEntity != nil && domainEntity.IsPrimary != nil && *domainEntity.IsPrimary {
							localPrimaryDomain = domain
							break
						}
					}
				}
				// invoke enrich organization by domain
				if localPrimaryDomain != "" {
					err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.RequestEnrichOrganization{Url: localPrimaryDomain})
					if err != nil {
						tracing.TraceErr(span, err)
					}
				}
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// request last touchpoint refresh for new organizations
			if createFlow {
				err = s.RequestRefreshLastTouchpoint(ctx, organizationId)
				if err != nil {
					tracing.TraceErr(span, err)
				}
			}

			return nil
		})
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return organizationId, nil
}

func (s *organizationService) Hide(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Hide")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		fields := data_fields.OrganizationFields{Hide: utils.BoolPtr(true)}
		err = s.services.Neo4jRepositories.OrganizationWriteRepository.Save(ctx, txWithPostCommit.Tx, tenant, organizationId, fields)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithDelete())
			return nil
		})
		return nil, nil
	})

	return err
}

func (s *organizationService) Show(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.Show")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		fields := data_fields.OrganizationFields{Hide: utils.BoolPtr(false)}
		err = s.services.Neo4jRepositories.OrganizationWriteRepository.Save(ctx, txWithPostCommit.Tx, tenant, organizationId, fields)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithCreate())
			err = s.RequestRefreshLastTouchpoint(ctx, organizationId)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			return nil
		})
		return nil, nil
	})

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

func (s *organizationService) GetPrimaryOrganizationsWithJobRoleForContacts(ctx context.Context, contactIds []string) (*neo4jentity.OrganizationWithJobRoleEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.GetPrimaryOrganizationsWithJobRoleForContacts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contactIds", contactIds))

	dbResults, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetPrimaryOrganizationsWithJobRoleForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
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

func (s *organizationService) LinkWithDomain(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, organizationId, domain string) error {
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

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		domainLinkedSuccessfully, err := s.services.Neo4jRepositories.OrganizationWriteRepository.LinkWithDomain(ctx, txWithPostCommit.Tx, tenant, organizationId, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to link domain in neo4j"))
			return nil, err
		}

		if domainLinkedSuccessfully {
			txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
				// send organization enrich request
				err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.RequestEnrichOrganization{Url: domain})
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to publish event RequestEnrichOrganization"))
				}
				return nil
			})

			txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
				// send event to rabbitmq
				err = s.services.RabbitMQService.PublishEvent(ctx, organizationId, model.ORGANIZATION, dto.NewAddDomainEvent(domain))
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to publish event AddDomain"))
				}

				// send event to events platform
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

				return nil
			})
		}

		return nil, nil
	})

	return err
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

	// fetch the real touchpoint
	// if it doesn't exist, check for the Created Action
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

	s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, organizationId, model.ORGANIZATION, utils.NewEventCompletedDetails().WithUpdate())

	return nil
}

func (s *organizationService) CheckOrganizationExistsWithEmail(ctx context.Context, email string) (bool, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.CheckOrganizationExistsWithEmail")
	defer span.Finish()

	if email == "" {
		return false, "", nil
	}
	// if not a valid email, return false
	syntaxValidation := mailsherpa.ValidateEmailSyntax(email)
	if !syntaxValidation.IsValid {
		return false, "", nil
	}

	orgs, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsWithEmail(ctx, common.GetTenantFromContext(ctx), email)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, "", err
	}
	orgId := ""
	if len(orgs) > 0 {
		orgId = orgs[0].Props["id"].(string)
	}
	return len(orgs) > 0, orgId, nil
}

func (s *organizationService) CheckOrganizationExistsWithLinkedIn(ctx context.Context, url, alias, externalId string) (bool, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationService.CheckOrganizationExistsWithLinkedIn")
	defer span.Finish()

	if alias == "" {
		// use identifier as alias
		alias = neo4jentity.SocialEntity{Url: url}.ExtractLinkedinPersonIdentifierFromUrl()
	}
	orgs, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationsByLinkedIn(ctx, common.GetTenantFromContext(ctx), url, alias, externalId)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, "", err
	}
	orgId := ""
	if len(orgs) > 0 {
		orgId = orgs[0].Props["id"].(string)
	}
	return len(orgs) > 0, orgId, nil
}
