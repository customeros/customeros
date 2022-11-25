package mocks

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/helper"
	"github.com/stretchr/testify/mock"
)

type MockedAppInfoRepo struct {
	mock.Mock
}

func (m *MockedAppInfoRepo) FindOneById(ctx context.Context, id string) helper.QueryResult {
	args := m.Called(ctx, id)
	return args.Get(0).(helper.QueryResult)
}

func (m *MockedAppInfoRepo) FindAll(ctx context.Context) helper.QueryResult {
	args := m.Called(ctx)
	return args.Get(0).(helper.QueryResult)
}
