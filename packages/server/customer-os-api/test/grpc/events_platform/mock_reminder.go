package events_platform

import (
	"context"

	reminderpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/reminder"
)

type MockReminderServiceCallbacks struct {
	ReminderCreate func(context.Context, *reminderpb.CreateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error)
	ReminderUpdate func(context.Context, *reminderpb.UpdateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error)
}

var reminderCallbacks = &MockReminderServiceCallbacks{}

func SetReminderCallbacks(callbacks *MockReminderServiceCallbacks) {
	reminderCallbacks = callbacks
}

type MockReminderService struct {
	reminderpb.UnimplementedReminderGrpcServiceServer
}

func (MockReminderService) CreateReminder(context context.Context, proto *reminderpb.CreateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error) {
	if reminderCallbacks.ReminderCreate == nil {
		panic("reminderCallbacks.CreateReminder is not set")
	}
	return reminderCallbacks.ReminderCreate(context, proto)
}

func (MockReminderService) UpdateReminder(context context.Context, proto *reminderpb.UpdateReminderGrpcRequest) (*reminderpb.ReminderGrpcResponse, error) {
	if reminderCallbacks.ReminderUpdate == nil {
		panic("reminderCallbacks.UpdateReminder is not set")
	}
	return reminderCallbacks.ReminderUpdate(context, proto)
}
