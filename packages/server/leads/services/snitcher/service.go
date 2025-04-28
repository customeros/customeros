package snitcher

import (
	"context"
	"errors"
	"fmt"
	"time"

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

type SnitcherService struct {
	config        *config.SnitcherConfig
	natsConn      *nats_internal.NATSConnections
	repositories  *repository.Repositories
	subscriptions []*nats.Subscription
}

func NewSnitcherService(config *config.SnitcherConfig, repos *repository.Repositories, natsConn *nats_internal.NATSConnections) *SnitcherService {
	return &SnitcherService{
		config:       config,
		natsConn:     natsConn,
		repositories: repos,
	}
}

var SUBSCRIBED_SUBJECT = enum.EventAskSnitcher.String()

const (
	HTTP_TIMEOUT      = 60 * time.Second
	MAX_RESPONSE_SIZE = 1 * 1024 * 1024
)

// Start begins listening for  events
func (s *SnitcherService) Start(ctx context.Context) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SnitcherService.Start")
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
func (s *SnitcherService) Close() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
}

func (s *SnitcherService) handleNatsMessage(ctx context.Context, msg *nats.Msg) {
	ctx = utils.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartServiceSpan(ctx, "SnitcherService.handleNatsMessage")
	defer spans.Finish()

	if msg == nil {
		spans.TraceError(errors.New("nil nats message"))
		return
	}
	spans.TagString("nats.subject", msg.Subject)
	spans.TagString("nats.reply", msg.Reply)

	resp := &pb.IPAddressIdentifyResponse{}

	request := &pb.IPAddressIdentifyRequest{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		errMsg := "Failed to parse request"
		resp.ErrorMessage = errMsg
		s.sendResponse(ctx, msg, resp)
		spans.TraceError(err)
		return
	}

	resp = s.AskSnitcher(ctx, request.IpAddress)
	if resp == nil {
		spans.TraceError(errors.New("empty response"))
		return
	}

	s.sendResponse(ctx, msg, resp)
}

func (s *SnitcherService) sendResponse(ctx context.Context, req *nats.Msg, resp *pb.IPAddressIdentifyResponse) {
	spans, _ := telemetry.StartServiceSpan(ctx, "SnitcherService.sendResponse")
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
