package service

import (
	"encoding/json"
	"fmt"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type LocationService interface {
	GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error)
	GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.LocationEntities, error)
	GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.LocationEntities, error)
	ExtractAndEnrichLocation(ctx context.Context, tenant, address string) (*data_fields.LocationFields, error)
	Create(ctx context.Context, commit *utils.TxWithPostCommit, locationFields data_fields.LocationFields, linkWith *LinkWith) (string, error)
}

type locationService struct {
	log      logger.Logger
	services *Services
}

func NewLocationService(log logger.Logger, services *Services) LocationService {
	return &locationService{
		log:      log,
		services: services,
	}
}

func (s *locationService) GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error) {
	dbNodes, err := s.services.Neo4jRepositories.LocationReadRepository.GetAllForContact(ctx, common.GetTenantFromContext(ctx), contactId)
	if err != nil {
		return nil, err
	}

	locationEntities := neo4jentity.LocationEntities{}
	for _, dbNode := range dbNodes {
		locationEntities = append(locationEntities, *neo4jmapper.MapDbNodeToLocationEntity(dbNode))
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForContacts(ctx context.Context, contactIds []string) (*neo4jentity.LocationEntities, error) {
	locations, err := s.services.Neo4jRepositories.LocationReadRepository.GetAllForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
	if err != nil {
		return nil, err
	}
	locationEntities := neo4jentity.LocationEntities{}
	for _, v := range locations {
		locationEntity := neo4jmapper.MapDbNodeToLocationEntity(v.Node)
		locationEntity.DataloaderKey = v.LinkedNodeId
		locationEntities = append(locationEntities, *locationEntity)
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForOrganization(ctx context.Context, organizationId string) (*neo4jentity.LocationEntities, error) {
	dbNodes, err := s.services.Neo4jRepositories.LocationReadRepository.GetAllForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		return nil, err
	}

	locationEntities := neo4jentity.LocationEntities{}
	for _, dbNode := range dbNodes {
		locationEntities = append(locationEntities, *neo4jmapper.MapDbNodeToLocationEntity(dbNode))
	}
	return &locationEntities, nil
}

func (s *locationService) GetAllForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.LocationEntities, error) {
	locations, err := s.services.Neo4jRepositories.LocationReadRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	locationEntities := neo4jentity.LocationEntities{}
	for _, v := range locations {
		locationEntity := neo4jmapper.MapDbNodeToLocationEntity(v.Node)
		locationEntity.DataloaderKey = v.LinkedNodeId
		locationEntities = append(locationEntities, *locationEntity)
	}
	return &locationEntities, nil
}

func (s *locationService) ExtractAndEnrichLocation(ctx context.Context, tenant, address string) (*data_fields.LocationFields, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationEventHandler.ExtractAndEnrichLocation")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("address", address))

	if strings.TrimSpace(address) == "" {
		return nil, errors.New("address is empty")
	}

	// Step 1: Check if mapping exists
	locationMapping, err := s.services.PostgresRepositories.AiLocationMappingRepository.GetLatestLocationMappingByInput(ctx, address)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get location mapping"))
	}
	if locationMapping != nil {
		var location data_fields.LocationFields
		err = json.Unmarshal([]byte(locationMapping.ResponseJson), &location)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal location"))
			return nil, err
		}
		return &location, nil
	}

	// Step 2: Use AI to enrich the location
	prompt := fmt.Sprintf(s.services.GlobalConfig.InternalServices.AiApiConfig.AnthropicPrompts.LocationEnrichmentPrompt, address)
	promptLog := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      common.GetAppSourceFromContext(ctx),
		Provider:       constants.Anthropic,
		Model:          enum.AIModelAnthropicHaiku.String(),
		PromptType:     constants.PromptTypeExtractLocationValue,
		Tenant:         &tenant,
		PromptTemplate: &s.services.GlobalConfig.InternalServices.AiApiConfig.AnthropicPrompts.LocationEnrichmentPrompt,
		Prompt:         prompt,
	}
	promptStoreLogId, err := s.services.PostgresRepositories.AiPromptLogRepository.Store(promptLog)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error storing prompt log: %v", err)
	}

	aiResult, err := s.services.AIService.AskAI(ctx, enum.AIModelAnthropicHaiku, &prompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get AI response"))
		s.log.Errorf("Error invoking AI: %s", err.Error())
		storeErr := s.services.PostgresRepositories.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with error"))
			s.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return nil, err
	}
	if aiResult == nil {
		return nil, nil
	}
	storeErr := s.services.PostgresRepositories.AiPromptLogRepository.UpdateResponse(promptStoreLogId, *aiResult)
	if storeErr != nil {
		tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with ai response"))
		s.log.Errorf("Error updating prompt log with ai response: %v", storeErr)
	}

	var location data_fields.LocationFields
	err = json.Unmarshal([]byte(*aiResult), &location)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal location"))
		return nil, err
	}

	// Step 3: Store the mapping
	locationMapping = &postgresEntity.AiLocationMapping{
		Input:         address,
		ResponseJson:  *aiResult,
		AiPromptLogId: promptStoreLogId,
	}
	err = s.services.PostgresRepositories.AiLocationMappingRepository.AddLocationMapping(ctx, *locationMapping)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store location mapping"))
		s.log.Errorf("Error storing location mapping: %v", err)
	}

	return &location, nil
}

func (s *locationService) Create(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, locationFields data_fields.LocationFields, linkWith *LinkWith) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LocationService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "locationFields", locationFields)
	tracing.LogObjectAsJson(span, "linkWith", linkWith)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	// validate link with
	if linkWith != nil {
		if !linkWith.IsContact() && !linkWith.IsOrganization() {
			err = errors.New("unsupported linkWith type")
			tracing.TraceErr(span, err)
			return "", err
		}
		// validate contact exists
		if linkWith.IsContact() {
			_, err = s.services.ContactService.GetContactById(ctx, linkWith.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
		}
		// validate organization exists
		if linkWith.IsOrganization() {
			_, err = s.services.OrganizationService.GetById(ctx, tenant, linkWith.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
		}
	}

	// set default values
	if locationFields.CreatedAt == nil {
		locationFields.CreatedAt = utils.NowPtr()
	}
	if locationFields.Source == nil {
		locationFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
	}
	if locationFields.AppSource == nil {
		locationFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
	}

	// generate location id
	locationId, err := s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelLocation)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err = s.services.Neo4jRepositories.LocationWriteRepository.CreateLocation(ctx, txWithPostCommit.Tx, tenant, locationId, locationFields)
		if err != nil {
			return "", err
		}

		if linkWith != nil {
			if linkWith.IsContact() {
				err = s.services.Neo4jRepositories.LocationWriteRepository.LinkWithContact(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, locationId)
				if err != nil {
					return "", err
				}
			} else if linkWith.IsOrganization() {
				err = s.services.Neo4jRepositories.LocationWriteRepository.LinkWithOrganization(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, locationId)
				if err != nil {
					return "", err
				}
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			innerErr := s.services.RabbitMQService.PublishEvent(ctx, locationId, model.CONTACT, dto.CreateLocation{locationFields})
			if innerErr != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateLocation"))
			}

			if linkWith != nil {
				s.services.RabbitMQService.PublishEventCompleted(ctx, tenant, linkWith.Id, linkWith.Type, utils.NewEventCompletedDetails().WithUpdate())
			}

			return nil
		})
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return locationId, nil
}
