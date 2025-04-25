package visitor_identity

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/interfaces"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

type VisitorIdentityService struct {
	natsConn      *nats_internal.NATSConnections
	repositories  *repository.Repositories
	subscriptions []*nats.Subscription
}

func NewVisitorIdentityService(
	natsConn *nats_internal.NATSConnections,
	repositories *repository.Repositories,
) interfaces.NatsService {
	return &VisitorIdentityService{
		natsConn:      natsConn,
		repositories:  repositories,
		subscriptions: make([]*nats.Subscription, 0),
	}
}

var SUBSCRIBED_SUBJECT = enum.EventAskSnitcher.String()

// Start begins listening for  events
func (s *VisitorIdentityService) Start(ctx context.Context) error {
	// Create a subscription for handling requests
	sub, err := s.natsConn.Conn.Subscribe(SUBSCRIBED_SUBJECT, func(msg *nats.Msg) {
		ctx = utils.WithCustomContextFromNats(ctx, msg)
		spans, ctx := telemetry.StartServiceSpan(ctx, "VisitorIdentityService.Start")
		defer spans.Finish()

		resp := &pb.IdentifyVisitorResponse{}

		request := &pb.IdentifyVisitorRequest{}
		err := proto.Unmarshal(msg.Data, request)
		if err != nil || request == nil {
			errMsg := "Failed to parse request"
			resp.ErrorMessage = errMsg
			s.sendResponse(ctx, msg, resp)
			spans.TraceError(err)
			return
		}

		// Process the request
		resp = s.identifyWebVisitor(ctx, request)
		if resp == nil {
			spans.TraceError(errors.New("empty response"))
			return
		}

		s.sendResponse(ctx, msg, resp)
	})
	if err != nil {
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
func (s *VisitorIdentityService) Stop() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
}

func (s *VisitorIdentityService) sendResponse(ctx context.Context, req *nats.Msg, resp *pb.IdentifyVisitorResponse) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "VisitorIdentityService.sendResponse")
	defer spans.Finish()

	respMessage, err := proto.Marshal(resp)
	if err != nil {
		spans.TraceError(err)
		return
	}
	err = req.Respond(respMessage)
	if err != nil {
		spans.TraceError(err)
		return

	}
	return
}
