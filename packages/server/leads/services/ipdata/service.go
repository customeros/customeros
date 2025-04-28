package ipdata

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/internal/config"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

type IPDataService struct {
	config        *config.IPDataConfig
	natsConn      *nats_internal.NATSConnections
	repositories  *repository.Repositories
	subscriptions []*nats.Subscription
}

func NewIPDataService(
	config *config.IPDataConfig,
	natsConn *nats_internal.NATSConnections,
	repos *repository.Repositories,
) *IPDataService {
	return &IPDataService{
		config:       config,
		natsConn:     natsConn,
		repositories: repos,
	}
}

var SUBSCRIBED_SUBJECT = enum.EventAskIPData.String()

// Start begins listening for  events
func (s *IPDataService) Start(ctx context.Context) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "IPDataService.Start")
	defer spans.Finish()

	// Create a subscription for handling requests
	sub, err := s.natsConn.Conn.Subscribe(SUBSCRIBED_SUBJECT, func(msg *nats.Msg) {
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
func (s *IPDataService) Close() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
}

func (s *IPDataService) handleNatsMessage(ctx context.Context, msg *nats.Msg) {
	ctx = utils.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartServiceSpan(ctx, "ipDataService.handleNatsMessage")
	defer spans.Finish()

	if msg == nil {
		spans.TraceError(errors.New("nil nats message"))
		return
	}
	spans.TagString("nats.subject", msg.Subject)
	spans.TagString("nats.reply", msg.Reply)

	resp := &pb.IPAddressVerifyResponse{}

	request := &pb.IPAddressIdentifyRequest{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		errMsg := "Failed to parse request"
		resp.ErrorMessage = errMsg
		s.sendResponse(ctx, msg, resp)
		spans.TraceError(err)
		return
	}

	resp = s.AskIPData(ctx, request.IpAddress)
	if resp == nil {
		spans.TraceError(errors.New("empty response"))
		return
	}

	s.sendResponse(ctx, msg, resp)
}

func (s *IPDataService) sendResponse(ctx context.Context, req *nats.Msg, resp *pb.IPAddressVerifyResponse) {
	spans, _ := telemetry.StartServiceSpan(ctx, "IPDataService.sendResponse")
	defer spans.Finish()

	respMessage, err := proto.Marshal(resp)
	if err != nil {
		spans.TraceError(err)
		return
	}
	err = req.Respond(respMessage)
	if err != nil {
		spans.TraceError(err)
	}
}
