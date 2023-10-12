package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/aggregate"
	cmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command"
	cmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/graph_db/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

type ActionForecastMetadata struct {
	Likelihood string `json:"likelihood"`
	Reason     string `json:"reason"`
}

type OrganizationEventHandler struct {
	repositories         *repository.Repositories
	organizationCommands *cmdhnd.OrganizationCommands
	log                  logger.Logger
}

type eventMetadata struct {
	UserId string `json:"user-id"`
}

func (h *OrganizationEventHandler) OnOrganizationCreate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationCreate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationCreateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var err error
		err = h.repositories.OrganizationRepository.CreateOrganizationInTx(ctx, tx, organizationId, eventData)
		if err != nil {
			h.log.Errorf("Error while saving organization %s: %s", organizationId, err.Error())
			return nil, err
		}
		if eventData.ExternalSystem.Available() {
			err = h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, constants.NodeLabel_Organization, eventData.ExternalSystem)
			if err != nil {
				h.log.Errorf("Error while link organization %s with external system %s: %s", organizationId, eventData.ExternalSystem.ExternalSystemId, err.Error())
				return nil, err
			}
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// set customer os id
	customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
	if customerOsErr != nil {
		tracing.TraceErr(span, customerOsErr)
		h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
	}

	// Set organization owner
	evtMetadata := eventMetadata{}
	if err = json.Unmarshal(evt.Metadata, &evtMetadata); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "json.Unmarshal")
	} else {
		if evtMetadata.UserId != "" {
			err = h.repositories.OrganizationRepository.ReplaceOwner(ctx, eventData.Tenant, organizationId, evtMetadata.UserId)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed to replace owner of organization %s with user %s", organizationId, evtMetadata.UserId)
			}
		}
	}

	// Set create action
	_, err = h.repositories.ActionRepository.MergeByActionType(ctx, eventData.Tenant, organizationId, entity.ORGANIZATION, entity.ActionCreated, "", "", eventData.CreatedAt)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed creating likelihood update action for organization %s: %s", organizationId, err.Error())
	}

	// Request last touch point update
	err = h.organizationCommands.RefreshLastTouchpointCommand.Handle(ctx, cmd.NewRefreshLastTouchpointCommand(eventData.Tenant, organizationId, evtMetadata.UserId))
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("RefreshLastTouchpointCommand failed: %v", err.Error())
	}

	return err
}

func (h *OrganizationEventHandler) setCustomerOsId(ctx context.Context, tenant, organizationId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.setCustomerOsId")
	defer span.Finish()
	span.LogFields(log.String("Tenant", tenant), log.String("OrganizationId", organizationId))

	orgDbNode, err := h.repositories.OrganizationRepository.GetOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	organizationEntity := graph_db.MapDbNodeToOrganizationEntity(*orgDbNode)

	if organizationEntity.CustomerOsId != "" {
		return nil
	}
	var customerOsId string
	maxAttempts := 20
	for attempt := 1; attempt < maxAttempts+1; attempt++ {
		customerOsId = generateNewRandomCustomerOsId()
		customerOsIdsEntity := postgresentity.CustomerOsIds{
			Tenant:       tenant,
			CustomerOSID: customerOsId,
			Entity:       postgresentity.Organization,
			EntityId:     organizationId,
			Attempts:     attempt,
		}
		innerErr := h.repositories.CustomerOsIdsRepository.Reserve(customerOsIdsEntity)
		if innerErr == nil {
			break
		}
	}
	return h.repositories.OrganizationRepository.SetCustomerOsIdIfMissing(ctx, tenant, organizationId, customerOsId)
}

func (h *OrganizationEventHandler) OnOrganizationUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	var err error
	if eventData.IgnoreEmptyFields {
		err = h.repositories.OrganizationRepository.UpdateOrganizationIgnoreEmptyInputParams(ctx, organizationId, eventData)
	} else {
		err = h.repositories.OrganizationRepository.UpdateOrganization(ctx, organizationId, eventData)
		// set customer os id
		customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
		if customerOsErr != nil {
			tracing.TraceErr(span, customerOsErr)
			h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
		}
	}
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Error while saving organization %s: %s", organizationId, err.Error())
		return err
	}

	if eventData.ExternalSystem.Available() {
		session := utils.NewNeo4jWriteSession(ctx, *h.repositories.Drivers.Neo4jDriver)
		defer session.Close(ctx)

		_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			//var err error
			if eventData.ExternalSystem.Available() {
				innerErr := h.repositories.ExternalSystemRepository.LinkWithEntityInTx(ctx, tx, eventData.Tenant, organizationId, constants.NodeLabel_Organization, eventData.ExternalSystem)
				if innerErr != nil {
					h.log.Errorf("Error while link organization %s with external system %s: %s", organizationId, eventData.ExternalSystem.ExternalSystemId, err.Error())
					return nil, innerErr
				}
			}
			return nil, nil
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (h *OrganizationEventHandler) OnPhoneNumberLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnPhoneNumberLinkedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkPhoneNumberEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.PhoneNumberRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.PhoneNumberId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnEmailLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnEmailLinkedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkEmailEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.EmailRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.EmailId, eventData.Label, eventData.Primary, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnLocationLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLocationLinkedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkLocationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.LocationRepository.LinkWithOrganization(ctx, eventData.Tenant, organizationId, eventData.LocationId, eventData.UpdatedAt)

	return err
}

func (h *OrganizationEventHandler) OnDomainLinkedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnDomainLinkedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationLinkDomainEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	if strings.TrimSpace(eventData.Domain) == "" {
		return nil
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.OrganizationRepository.LinkWithDomain(ctx, eventData.Tenant, organizationId, strings.TrimSpace(eventData.Domain))

	return err
}

func (h *OrganizationEventHandler) OnSocialAddedToOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnSocialAddedToOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationAddSocialEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.SocialRepository.CreateSocialFor(ctx, eventData.Tenant, organizationId, "Organization", eventData)

	return err
}

func (h *OrganizationEventHandler) OnRenewalLikelihoodUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRenewalLikelihoodUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateRenewalLikelihoodEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		err = errors.Wrap(err, "GetJsonData")
		tracing.TraceErr(span, err)
		return err
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	err := h.repositories.OrganizationRepository.UpdateRenewalLikelihood(ctx, organizationId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	if eventData.PreviousLikelihood != eventData.RenewalLikelihood {
		if string(eventData.RenewalLikelihood) != "" {
			userDbNode, err := h.repositories.UserRepository.GetUser(ctx, eventData.Tenant, eventData.UpdatedBy)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("GetUser failed for id: %s", eventData.UpdatedBy, err.Error())
			}
			message := "Renewal likelihood set to " + eventData.RenewalLikelihood.CamelCaseString()
			if userDbNode != nil {
				userEntity := graph_db.MapDbNodeToUserEntity(*userDbNode)
				message += " by " + userEntity.FirstName + " " + userEntity.LastName
			}
			metadata, err := utils.ToJson(ActionForecastMetadata{
				Likelihood: string(eventData.RenewalLikelihood),
				Reason:     utils.IfNotNilString(eventData.Comment),
			})
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("ToJson failed: %s", err.Error())
			}
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, organizationId, entity.ORGANIZATION, entity.ActionRenewalLikelihoodUpdated, message, metadata, eventData.UpdatedAt)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating likelihood update action for organization %s: %s", organizationId, err.Error())
			}
		}

		err = h.organizationCommands.RequestRenewalForecastCommand.Handle(ctx, cmd.NewRequestRenewalForecastCommand(eventData.Tenant, organizationId, ""))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestRenewalForecastCommand failed: %v", err.Error())
		}
	}

	return err
}

func (h *OrganizationEventHandler) OnRenewalForecastUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRenewalForecastUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateRenewalForecastEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		err = errors.Wrap(err, "GetJsonData")
		tracing.TraceErr(span, err)
		return err
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	err := h.repositories.OrganizationRepository.UpdateRenewalForecast(ctx, organizationId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// If the amount has changed, create an action
	if eventData.Amount != nil && !utils.Float64PtrEquals(eventData.Amount, eventData.PreviousAmount) {
		message := ""
		strAmount := utils.FormatCurrencyAmount(*eventData.Amount)
		if eventData.UpdatedBy == "" && string(eventData.RenewalLikelihood) != "" {
			if eventData.RenewalLikelihood == models.RenewalLikelihoodHIGH {
				message = fmt.Sprintf("Renewal forecast set by default to $%s based on the billing amount", strAmount)
			} else {
				message = fmt.Sprintf("Renewal forecast set by default to $%s, by discounting the billing amount using the renewal likelihood", strAmount)
			}
		} else if eventData.UpdatedBy != "" {
			userDbNode, err := h.repositories.UserRepository.GetUser(ctx, eventData.Tenant, eventData.UpdatedBy)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("GetUser failed for id: %s", eventData.UpdatedBy, err.Error())
			}
			message = fmt.Sprintf("Renewal forecast set to $%s", strAmount)
			if userDbNode != nil {
				userEntity := graph_db.MapDbNodeToUserEntity(*userDbNode)
				if userEntity.FirstName != "" || userEntity.LastName != "" {
					message += " by " + userEntity.FirstName + " " + userEntity.LastName
				}
			}
		}
		metadata, err := utils.ToJson(ActionForecastMetadata{
			Likelihood: string(eventData.RenewalLikelihood),
			Reason:     utils.IfNotNilString(eventData.Comment),
		})
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("ToJson failed: %s", err.Error())
		}
		if message != "" {
			_, err = h.repositories.ActionRepository.Create(ctx, eventData.Tenant, organizationId, entity.ORGANIZATION, entity.ActionRenewalForecastUpdated, message, metadata, eventData.UpdatedAt)
			if err != nil {
				tracing.TraceErr(span, err)
				h.log.Errorf("Failed creating forecast update action for organization %s: %s", organizationId, err.Error())
			}
		}
	}

	if eventData.UpdatedBy != "" && eventData.Amount == nil {
		err = h.organizationCommands.RequestRenewalForecastCommand.Handle(ctx, cmd.NewRequestRenewalForecastCommand(eventData.Tenant, organizationId, ""))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestRenewalForecastCommand failed: %v", err.Error())
		}
	}

	return err
}

func (h *OrganizationEventHandler) OnBillingDetailsUpdate(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRenewalForecastUpdate")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpdateBillingDetailsEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		err = errors.Wrap(err, "GetJsonData")
		tracing.TraceErr(span, err)
		return err
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	err := h.repositories.OrganizationRepository.UpdateBillingDetails(ctx, organizationId, eventData)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	if eventData.UpdatedBy != "" {
		err = h.organizationCommands.RequestRenewalForecastCommand.Handle(ctx, cmd.NewRequestRenewalForecastCommand(eventData.Tenant, organizationId, ""))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestRenewalForecastCommand failed: %v", err.Error())
		}
		err = h.organizationCommands.RequestNextCycleDateCommand.Handle(ctx, cmd.NewRequestNextCycleDateCommand(eventData.Tenant, organizationId, ""))
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("RequestNextCycleDateCommand failed: %v", err.Error())
		}
	}

	return err
}

func (h *OrganizationEventHandler) OnOrganizationHide(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationHide")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.HideOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.OrganizationRepository.SetVisibility(ctx, eventData.Tenant, organizationId, true)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (h *OrganizationEventHandler) OnOrganizationShow(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnOrganizationShow")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.HideOrganizationEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	err := h.repositories.OrganizationRepository.SetVisibility(ctx, eventData.Tenant, organizationId, false)
	if err != nil {
		tracing.TraceErr(span, err)
	}

	// set customer os id
	customerOsErr := h.setCustomerOsId(ctx, eventData.Tenant, organizationId)
	if customerOsErr != nil {
		tracing.TraceErr(span, customerOsErr)
		h.log.Errorf("Failed to set customer os id for tenant %s organization %s", eventData.Tenant, organizationId)
	}

	return err
}

func (h *OrganizationEventHandler) OnRefreshLastTouchpoint(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnRefreshLastTouchpoint")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationRefreshLastTouchpointEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	lastTouchpointAt, lastTouchpointId, err := h.repositories.TimelineEventRepository.CalculateAndGetLastTouchpoint(ctx, eventData.Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to calculate last touchpoint: %v", err.Error())
		span.LogFields(log.Bool("last touchpoint failed", true))
		return nil
	}

	if lastTouchpointAt == nil {
		h.log.Infof("Last touchpoint not available for organization: %s", organizationId)
		span.LogFields(log.Bool("last touchpoint not found", true))
		return nil
	}

	if err = h.repositories.OrganizationRepository.UpdateLastTouchpoint(ctx, eventData.Tenant, organizationId, *lastTouchpointAt, lastTouchpointId); err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to update last touchpoint for tenant %s, organization %s: %s", eventData.Tenant, organizationId, err.Error())
	}

	return nil
}

func (h *OrganizationEventHandler) OnUpsertCustomField(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUpsertCustomField")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationUpsertCustomField
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)

	customFieldExists, err := h.repositories.CustomFieldRepository.ExistsById(ctx, eventData.Tenant, eventData.CustomFieldId)
	if err != nil {
		tracing.TraceErr(span, err)
		h.log.Errorf("Failed to check if custom field exists: %s", err.Error())
		return err
	}
	if !customFieldExists {
		err = h.repositories.CustomFieldRepository.AddCustomFieldToOrganization(ctx, eventData.Tenant, organizationId, eventData)
		if err != nil {
			tracing.TraceErr(span, err)
			h.log.Errorf("Failed to add custom field to organization: %s", err.Error())
			return err
		}
	} else {
		//TODO implement update custom field
	}

	return nil
}

func (h *OrganizationEventHandler) OnLinkWithParentOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnLinkWithParentOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationAddParentEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	return h.repositories.OrganizationRepository.LinkWithParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId, eventData.Type)
}

func (h *OrganizationEventHandler) OnUnlinkFromParentOrganization(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationEventHandler.OnUnlinkFromParentOrganization")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData events.OrganizationRemoveParentEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	organizationId := aggregate.GetOrganizationObjectID(evt.AggregateID, eventData.Tenant)
	return h.repositories.OrganizationRepository.UnlinkParentOrganization(ctx, eventData.Tenant, organizationId, eventData.ParentOrganizationId)
}
