package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	eventstorepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/event_store"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/reminder/event"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, tenant, userId, orgId, content string, dueDate time.Time) (string, error)
	UpdateReminder(ctx context.Context, tenant, id string, content *string, dueDate *time.Time, dismissed *bool) error
	GetReminderById(ctx context.Context, id string) (*neo4jentity.ReminderEntity, error)
	RemindersForOrganization(ctx context.Context, organizationID string, dismissed *bool) ([]*neo4jentity.ReminderEntity, error)
}

type reminderService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewReminderService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) ReminderService {
	return &reminderService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *reminderService) CreateReminder(ctx context.Context, tenant, userId, organizationId, content string, dueDate time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.CreateReminder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId), log.String("organizationId", organizationId), log.String("content", content), log.String("dueDate", dueDate.String()))

	evt, err := json.Marshal(event.ReminderCreateEvent{
		BaseEvent: events.BaseEvent{
			Tenant:     tenant,
			EventName:  event.ReminderCreateV1,
			CreatedAt:  time.Now().UTC(),
			AppSource:  constants.AppSourceCustomerOsApi,
			Source:     neo4jentity.DataSourceOpenline.String(),
			EntityType: model.REMINDER,
		},
		Content:        content,
		DueDate:        dueDate.UTC(),
		UserId:         userId,
		OrganizationId: organizationId,
		Dismissed:      false,
	})

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*eventstorepb.StoreEventGrpcResponse](func() (*eventstorepb.StoreEventGrpcResponse, error) {
		return s.grpcClients.EventStoreClient.StoreEvent(ctx, &eventstorepb.StoreEventGrpcRequest{
			EventDataBytes: evt,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	return response.Id, nil
}

func (s *reminderService) UpdateReminder(ctx context.Context, tenant, reminderId string, content *string, dueDate *time.Time, dismissed *bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.UpdateReminder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, reminderId)

	if content == nil && dueDate == nil && dismissed == nil {
		return nil
	}

	fieldsMask := make([]string, 0, 3)
	if content != nil {
		fieldsMask = append(fieldsMask, event.FieldMaskReminderContent)
	}
	if dueDate != nil {
		fieldsMask = append(fieldsMask, event.FieldMaskReminderDueDate)
	}
	if dismissed != nil {
		fieldsMask = append(fieldsMask, event.FieldMaskReminderDismissed)
	}

	evt, err := json.Marshal(event.ReminderUpdateEvent{
		BaseEvent: events.BaseEvent{
			Tenant:     tenant,
			EventName:  event.ReminderUpdateV1,
			CreatedAt:  time.Now().UTC(),
			AppSource:  constants.AppSourceCustomerOsApi,
			Source:     neo4jentity.DataSourceOpenline.String(),
			EntityId:   reminderId,
			EntityType: model.REMINDER,
		},
		Content:    utils.IfNotNilString(content),
		DueDate:    dueDate.UTC(),
		Dismissed:  utils.IfNotNilBool(dismissed),
		FieldsMask: fieldsMask,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*eventstorepb.StoreEventGrpcResponse](func() (*eventstorepb.StoreEventGrpcResponse, error) {
		return s.grpcClients.EventStoreClient.StoreEvent(ctx, &eventstorepb.StoreEventGrpcRequest{
			EventDataBytes: evt,
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *reminderService) GetReminderById(ctx context.Context, id string) (*neo4jentity.ReminderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.GetReminderById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, id)

	if reminderDbNode, err := s.repositories.Neo4jRepositories.ReminderReadRepository.GetReminderById(ctx, common.GetContext(ctx).Tenant, id); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Reminder with id {%s} not found", id))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToReminderEntity(reminderDbNode), nil
	}
}

func (s *reminderService) RemindersForOrganization(ctx context.Context, organizationID string, dismissed *bool) ([]*neo4jentity.ReminderEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.RemindersForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationID)

	reminderDbNodes, err := s.repositories.Neo4jRepositories.ReminderReadRepository.GetRemindersOrderByDueDateAsc(ctx, common.GetContext(ctx).Tenant, organizationID, dismissed)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	reminderEntities := make([]*neo4jentity.ReminderEntity, 0, len(reminderDbNodes))
	if len(reminderDbNodes) == 0 {
		span.LogFields(log.String("Warning", fmt.Sprintf("Reminders for organization with id {%s} not found", organizationID)))
		return reminderEntities, nil
	}
	for _, v := range reminderDbNodes {
		reminderEntities = append(reminderEntities, neo4jmapper.MapDbNodeToReminderEntity(v))
	}
	return reminderEntities, nil
}
