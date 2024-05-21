package service

import (
	"context"
	"fmt"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ReminderService interface {
	CreateReminder(ctx context.Context, userId, orgId, content string, dueDate time.Time) (string, error)
	UpdateReminder(ctx context.Context, id string, content *string, dueDate *time.Time, dismissed *bool) error
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

func (s *reminderService) CreateReminder(ctx context.Context, userId, orgId, content string, dueDate time.Time) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.CreateReminder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("userId", userId), log.String("orgId", orgId), log.String("content", content), log.String("dueDate", dueDate.String()))

	createdAt := time.Now().UTC()
	due := dueDate.UTC()
	grpcRequest := &reminderpb.CreateReminderGrpcRequest{
		LoggedInUserId: userId,
		OrganizationId: orgId,
		Content:        content,
		DueDate:        utils.ConvertTimeToTimestampPtr(&due),
		Dismissed:      false,
		Tenant:         common.GetTenantFromContext(ctx),
		CreatedAt:      utils.ConvertTimeToTimestampPtr(&createdAt),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*reminderpb.ReminderGrpcResponse](func() (*reminderpb.ReminderGrpcResponse, error) {
		return s.grpcClients.ReminderClient.CreateReminder(ctx, grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelReminder, span)

	return response.Id, nil
}

func (s *reminderService) UpdateReminder(ctx context.Context, id string, content *string, dueDate *time.Time, dismissed *bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderService.UpdateReminder")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, id)

	if content == nil && dueDate == nil && dismissed == nil {
		return nil
	}

	updatedAt := time.Now().UTC()
	grpcRequest := &reminderpb.UpdateReminderGrpcRequest{
		AppSource:  constants.AppSourceCustomerOsApi,
		Content:    utils.IfNotNilString(content),
		DueDate:    utils.ConvertTimeToTimestampPtr(dueDate),
		Dismissed:  utils.IfNotNilBool(dismissed),
		UpdatedAt:  utils.ConvertTimeToTimestampPtr(&updatedAt),
		ReminderId: id,
		Tenant:     common.GetTenantFromContext(ctx),
	}
	fieldsMask := make([]reminderpb.ReminderFieldMask, 0, 3)
	if content != nil {
		fieldsMask = append(fieldsMask, reminderpb.ReminderFieldMask_REMINDER_PROPERTY_CONTENT)
	}
	if dueDate != nil {
		fieldsMask = append(fieldsMask, reminderpb.ReminderFieldMask_REMINDER_PROPERTY_DUE_DATE)
	}
	if dismissed != nil {
		fieldsMask = append(fieldsMask, reminderpb.ReminderFieldMask_REMINDER_PROPERTY_DISMISSED)
	}
	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	grpcRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*reminderpb.ReminderGrpcResponse](func() (*reminderpb.ReminderGrpcResponse, error) {
		return s.grpcClients.ReminderClient.UpdateReminder(ctx, grpcRequest)
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
