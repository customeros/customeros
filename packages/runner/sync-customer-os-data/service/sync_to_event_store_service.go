package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	email_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/email"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const syncToEventStoreBatchSize = 100

type SyncToEventStoreService interface {
	SyncEmails(ctx context.Context)
}

type syncToEventStoreService struct {
	repositories *repository.Repositories
	services     *Services
	grpcClients  *grpc_client.Clients
	batchSize    int
}

func NewSyncToEventStoreService(repositories *repository.Repositories, services *Services, grpcClients *grpc_client.Clients) SyncToEventStoreService {
	return &syncToEventStoreService{
		repositories: repositories,
		services:     services,
		grpcClients:  grpcClients,
		batchSize:    syncToEventStoreBatchSize,
	}
}

func (s *syncToEventStoreService) SyncEmails(ctx context.Context) {
	logrus.Infof("start sync emails to eventstore at %v", utils.Now())
	completedCount := 0
	failedCount := 0

	completedCount, failedCount, _ = s.upsertEmailsIntoEventStore(ctx)

	logrus.Infof("completed %v and faled %v emails upserting to eventstore at %v", completedCount, failedCount, utils.Now())
}

func (s *syncToEventStoreService) upsertEmailsIntoEventStore(ctx context.Context) (int, int, error) {
	processedRecords := 0
	failedRecords := 0
	records, err := s.repositories.EmailRepository.GetAllCrossTenantsWithRawEmail(ctx, s.batchSize)
	if err != nil {
		return 0, 0, err
	}
	for _, v := range records {
		_, err := s.grpcClients.EmailClient.UpsertEmail(context.Background(), &email_grpc_service.UpsertEmailGrpcRequest{
			Id:            utils.GetStringPropOrEmpty(v.Node.Props, "id"),
			Tenant:        v.LinkedNodeId,
			RawEmail:      utils.GetStringPropOrEmpty(v.Node.Props, "rawEmail"),
			Source:        utils.GetStringPropOrEmpty(v.Node.Props, "source"),
			SourceOfTruth: utils.GetStringPropOrEmpty(v.Node.Props, "sourceOfTruth"),
			AppSource:     utils.GetStringPropOrEmpty(v.Node.Props, "appSource"),
			CreatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "createdAt")),
			UpdatedAt:     utils.ConvertTimeToTimestampPtr(utils.GetTimePropOrNil(v.Node.Props, "updatedAt")),
		})
		if err != nil {
			failedRecords++
			logrus.Errorf("Failed to call method: %v", err)
		} else {
			processedRecords++
		}
	}

	return processedRecords, failedRecords, nil
}
