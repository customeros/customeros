package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	contactpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contact"
	"time"
)

type ContactService interface {
	UpkeepContacts()
	FindEmails()
}

type contactService struct {
	cfg                    *config.Config
	log                    logger.Logger
	repositories           *repository.Repositories
	eventsProcessingClient *grpc_client.Clients
}

func NewContactService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, client *grpc_client.Clients) ContactService {
	return &contactService{
		cfg:                    cfg,
		log:                    log,
		repositories:           repositories,
		eventsProcessingClient: client,
	}
}

func (s *contactService) UpkeepContacts() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil.")
		return
	}

	now := utils.Now()

	s.removeDuplicatedSocials(ctx, now)
}

func (s *contactService) removeDuplicatedSocials(ctx context.Context, now time.Time) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.removeDuplicatedSocials")
	defer span.Finish()

	limit := 100

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.SocialReadRepository.GetDuplicatedSocialsForEntityType(ctx, neo4jutil.NodeLabelContact, 180, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting socials: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		//remove socials from contact
		for _, record := range records {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*contactpb.ContactIdGrpcResponse](func() (*contactpb.ContactIdGrpcResponse, error) {
				return s.eventsProcessingClient.ContactClient.RemoveSocial(ctx, &contactpb.ContactRemoveSocialGrpcRequest{
					Tenant:    record.Tenant,
					ContactId: record.LinkedEntityId,
					SocialId:  record.SocialId,
					AppSource: constants.AppSourceDataUpkeeper,
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error removing social {%s}: %s", record.SocialId, err.Error())
			}

		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *contactService) FindEmails() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel context on exit

	if s.eventsProcessingClient == nil {
		s.log.Warn("eventsProcessingClient is nil.")
		return
	}

	s.findEmailsWithBetterContact(ctx)
}

func (s *contactService) findEmailsWithBetterContact(ctx context.Context) {
	span, ctx := tracing.StartTracerSpan(ctx, "ContactService.findEmailsWithBetterContact")
	defer span.Finish()

	limit := 0 // not yet enabled

	for {
		select {
		case <-ctx.Done():
			s.log.Infof("Context cancelled, stopping")
			return
		default:
			// continue as normal
		}

		records, err := s.repositories.Neo4jRepositories.ContactReadRepository.GetContactsToEnrichEmail(ctx, 2, limit)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting socials: %v", err)
			return
		}

		// no record
		if len(records) == 0 {
			return
		}

		for _, record := range records {
			err = s.requestBetterContactToFindEmail(ctx, record)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error requesting better contact to find email: %s", err.Error())
			} else {
				// mark contact with enrich requested
				err = s.repositories.Neo4jRepositories.ContactWriteRepository.UpdateTimeProperty(ctx, record.Tenant, record.ContractId, "techFindEmailRequestedAt", utils.NowPtr())
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error updating contact' find email requested: %s", err.Error())
				}
			}
		}

		// if less than limit records are returned, we are done
		if len(records) < limit {
			return
		}

		// force exit after single iteration
		return
	}
}

func (s *contactService) requestBetterContactToFindEmail(ctx context.Context, record neo4jrepository.TenantAndContactDetails) error {
	return nil
}
