package organization

import (
	"context"
	"fmt"

	nats_core "github.com/customeros/customeros/packages/server/core-crm/nats"
	"github.com/customeros/customeros/packages/server/core-crm/proto/pb"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type organizationService struct {
	log             logger.Logger
	natsConn        *nats_core.NATSConnections
	postgres        *postgres_repository.Repositories
	neo4j           *neo4j_repository.Repositories
	events          *events.EventsService
	domain          interfaces.DomainService
	industry        interfaces.IndustryService
	user            interfaces.UserService
	social          interfaces.SocialService
	currencyService interfaces.CurrencyService
	contractService interfaces.ContractService
	subscriptions   []*nats.Subscription
}

func NewOrganizationService(log logger.Logger,
	natsConn *nats_core.NATSConnections,
	postgres *postgres_repository.Repositories,
	neo4j *neo4j_repository.Repositories,
	events *events.EventsService,
	domain interfaces.DomainService,
	industry interfaces.IndustryService,
	social interfaces.SocialService,
	user interfaces.UserService,
	currencyService interfaces.CurrencyService,
) interfaces.OrganizationService {
	return &organizationService{
		log:             log,
		postgres:        postgres,
		neo4j:           neo4j,
		events:          events,
		domain:          domain,
		industry:        industry,
		user:            user,
		social:          social,
		currencyService: currencyService,
		subscriptions:   make([]*nats.Subscription, 0),
	}
}

var SUBSCRIBED_SUBJECT = "core.organization.>"

const QUEUE_GROUP = "organization-service"

func (s *organizationService) Start(ctx context.Context) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "organizationService.Start")
	defer spans.Finish()

	// Create a subscription for handling requests
	sub, err := s.natsConn.Conn.QueueSubscribe(SUBSCRIBED_SUBJECT, QUEUE_GROUP, func(msg *nats.Msg) {
		s.handleNatsMessage(ctx, msg)
	})
	if err != nil {
		spans.TraceError(err)
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Keep track of subscription for cleanup
	s.subscriptions = append(s.subscriptions, sub)

	// Listen for context cancellation to clean up
	go func() {
		<-ctx.Done()
		for _, sub := range s.subscriptions {
			sub.Unsubscribe()
		}
	}()

	return nil
}

// Close gracefully shuts down the service
func (s *organizationService) Stop() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
}

func (s *organizationService) handleNatsMessage(ctx context.Context, msg *nats.Msg) {
	ctx = common.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartListenerSpan(ctx, "organizationService.handleNatsMessage")
	defer spans.Finish()

	if msg == nil {
		spans.TraceError(errors.New("nil nats message"))
		return
	}
	spans.TagString("nats.subject", msg.Subject)
	spans.TagString("nats.reply", msg.Reply)

	resp := &pb.OrganizationSaveResponse{}

	request := &pb.OrganizationSaveRequest{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		s.sendResponse(ctx, msg, resp)
		spans.TraceError(err)
		return
	}

	if resp == nil {
		spans.TraceError(errors.New("empty response"))
		return
	}

	s.sendResponse(ctx, msg, resp)
}

func (s *organizationService) sendResponse(ctx context.Context, req *nats.Msg, resp *pb.OrganizationSaveResponse) {
	spans, _ := telemetry.StartServiceSpan(ctx, "organizationService.sendResponse")
	defer spans.Finish()

	respMessage, err := proto.Marshal(resp)
	if err != nil {
		spans.TraceError(err)
		return
	}
	req.Respond(respMessage)
}

func (s *organizationService) SetSocialService(social interfaces.SocialService) {
	s.social = social
}

func (s *organizationService) SetContractService(contractService interfaces.ContractService) {
	s.contractService = contractService
}

func (s *organizationService) IsInitialized() bool {
	return utils.IsInitialized(s)
}
