package location

import (
	"encoding/json"
	"fmt"
	"strings"

	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	common_srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type locationService struct {
	log      logger.Logger
	neo4j    *neoRepo.Repositories
	postgres *repository.Repositories
	events   *events.EventsService
	cfg      *config.AnthropicPrompts
	ai       interfaces.AIService
	contact  interfaces.ContactService
	org      interfaces.OrganizationService
}

func NewLocationService(log logger.Logger, neo4j *neoRepo.Repositories, postgres *repository.Repositories, events *events.EventsService, cfg *config.AnthropicPrompts, ai interfaces.AIService, contact interfaces.ContactService, org interfaces.OrganizationService) interfaces.LocationService {
	return &locationService{
		log:      log,
		neo4j:    neo4j,
		postgres: postgres,
		events:   events,
		cfg:      cfg,
		ai:       ai,
		contact:  contact,
		org:      org,
	}
}

func (s *locationService) SetContactService(contact interfaces.ContactService) {
	s.contact = contact
}

func (s *locationService) SetOrganizationService(org interfaces.OrganizationService) {
	s.org = org
}

func (s *locationService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *locationService) GetAllForContact(ctx context.Context, contactId string) (*neo4jentity.LocationEntities, error) {
	dbNodes, err := s.neo4j.LocationReadRepository.GetAllForContact(ctx, common.GetTenantFromContext(ctx), contactId)
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
	locations, err := s.neo4j.LocationReadRepository.GetAllForContacts(ctx, common.GetTenantFromContext(ctx), contactIds)
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
	dbNodes, err := s.neo4j.LocationReadRepository.GetAllForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId)
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
	locations, err := s.neo4j.LocationReadRepository.GetAllForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
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
	locationMapping, err := s.postgres.AiLocationMappingRepository.GetLatestLocationMappingByInput(ctx, address)
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
	prompt := fmt.Sprintf(s.cfg.LocationEnrichmentPrompt, address)
	promptLog := postgresEntity.AiPromptLog{
		CreatedAt:      utils.Now(),
		AppSource:      common.GetAppSourceFromContext(ctx),
		Provider:       constants.Anthropic,
		Model:          enum.AIModelAnthropicHaiku.String(),
		PromptType:     constants.PromptTypeExtractLocationValue,
		Tenant:         &tenant,
		PromptTemplate: &s.cfg.LocationEnrichmentPrompt,
		Prompt:         prompt,
	}
	promptStoreLogId, err := s.postgres.AiPromptLogRepository.Store(promptLog)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error storing prompt log: %v", err)
	}

	aiResult, err := s.ai.AskAI(ctx, enum.AIModelAnthropicHaiku, &prompt)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get AI response"))
		s.log.Errorf("Error invoking AI: %s", err.Error())
		storeErr := s.postgres.AiPromptLogRepository.UpdateError(promptStoreLogId, err.Error())
		if storeErr != nil {
			tracing.TraceErr(span, errors.Wrap(storeErr, "failed to update prompt log with error"))
			s.log.Errorf("Error updating prompt log with error: %v", storeErr)
		}
		return nil, err
	}
	if aiResult == nil {
		return nil, nil
	}
	storeErr := s.postgres.AiPromptLogRepository.UpdateResponse(promptStoreLogId, *aiResult)
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
	err = s.postgres.AiLocationMappingRepository.AddLocationMapping(ctx, *locationMapping)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to store location mapping"))
		s.log.Errorf("Error storing location mapping: %v", err)
	}

	return &location, nil
}

func (s *locationService) Create(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, locationFields data_fields.LocationFields, linkWith *common_srv.LinkWith) (string, error) {
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
			_, err = s.contact.GetContactById(ctx, linkWith.Id)
			if err != nil {
				tracing.TraceErr(span, err)
				return "", err
			}
		}
		// validate organization exists
		if linkWith.IsOrganization() {
			_, err = s.org.GetById(ctx, tenant, linkWith.Id)
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
	locationId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelLocation)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err = s.neo4j.LocationWriteRepository.CreateLocation(ctx, txWithPostCommit.Tx, tenant, locationId, locationFields)
		if err != nil {
			return "", err
		}

		if linkWith != nil {
			if linkWith.IsContact() {
				err = s.neo4j.LocationWriteRepository.LinkWithContact(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, locationId)
				if err != nil {
					return "", err
				}
			} else if linkWith.IsOrganization() {
				err = s.neo4j.LocationWriteRepository.LinkWithOrganization(ctx, txWithPostCommit.Tx, tenant, linkWith.Id, locationId)
				if err != nil {
					return "", err
				}
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			innerErr := s.events.Publisher.PublishEvent(ctx, locationId, model.CONTACT, dto.CreateLocation{locationFields})
			if innerErr != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateLocation"))
			}

			if linkWith != nil {
				s.events.Publisher.PublishEventCompleted(ctx, tenant, linkWith.Id, linkWith.Type, utils.NewEventCompletedDetails().WithUpdate())
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
