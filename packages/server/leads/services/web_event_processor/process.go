package web_event_processor

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/dto"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

func (s *webEventProcessor) Process(ctx context.Context, event *dto.WebTrackerEvent, webtrackerID string) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.Process")
	defer span.Finish()

	// check if bot, return early if not trusted IP
	userAgent := utils.ParseUserAgent(event.UserAgent)
	trusted, err := s.isTrustedIP(ctx, event.IP)
	isSuspicious := utils.IsSuspiciousURL(event.Referrer)

	if !trusted || userAgent.IsBot || isSuspicious {
		return
	}

	// attach to session
	sessionID, err := s.attachToSession(ctx, event, webtrackerID)
	if err != nil {
		err = fmt.Errorf("unable to build tracking record")
		span.TraceError(err)
		return
	}

	fmt.Printf(sessionID)

	// attempt to identify

	// write web event to outbox
}

func (s *webEventProcessor) isTrustedIP(ctx context.Context, ipAddress string) (bool, error) {
	span, ctx := telemetry.StartRestSpan(ctx, "webEventProcessor.isTrustedIP")
	defer span.Finish()

	request := &pb.IPAddressVerifyRequest{
		IpAddress: ipAddress,
	}

	reqData, err := proto.Marshal(request)
	if err != nil {
		span.TraceError(err)
		return false, err
	}

	// Send request to service
	msg := nats.NewMsg(enum.EventAskIPData.String())
	msg.Header = nats.Header{
		enum.TENANT_HEADER:  []string{utils.GetTenantFromContext(ctx)},
		enum.USER_ID_HEADER: []string{utils.GetUserIdFromContext(ctx)},
	}
	msg.Data = reqData

	resp, err := s.natsConn.Conn.RequestMsg(msg, REQUEST_TIMEOUT)
	if err != nil {
		span.TraceError(err)
		return false, err
	}

	// Unmarshal response
	response := &pb.IPAddressVerifyResponse{}
	if err := proto.Unmarshal(resp.Data, response); err != nil {
		span.TraceError(err)
		return false, err
	}

	return !response.IsThreat, nil
}

func (s *webEventProcessor) attachToSession(ctx context.Context, event *dto.WebTrackerEvent, webtrackerID string) (string, error) {
	span, ctx := telemetry.StartRestSpan(ctx, "webEventProcessor.attachToSession")
	defer span.Finish()

	// check to see if sessionID exists in Nats KV
	sessionID, err := s.natsConn.SessionCache.Get(ctx, event.VisitorID)
	if err != nil {
		span.TraceError(err)
		return "", err
	}
	if sessionID != "" {
		return sessionID, nil
	}

	// generate new ID
	sessionID = utils.GenerateNanoIDWithPrefix("sess", 21)
	err = s.natsConn.SessionCache.Set(ctx, event.VisitorID, sessionID)
	if err != nil {
		span.TraceError(err)
		return "", err
	}

	// publish webtracker.session.new event
	err = s.publishNewSessionEvent(ctx, webtrackerID, sessionID, event.IP)
	if err != nil {
		span.TraceError(err)
		return "", err
	}

	return sessionID, nil
}

func (s *webEventProcessor) publishNewSessionEvent(ctx context.Context, trackerID, sessionID, IPAddress string) error {
	span, ctx := telemetry.StartRestSpan(ctx, "webEventProcessor.publishNewSessionEvent")
	defer span.Finish()

	message := &pb.WebtrackerSessionNew{
		Id:        utils.GenerateEventID(),
		TrackerId: trackerID,
		SessionId: sessionID,
		Ip:        IPAddress,
	}

	data, err := proto.Marshal(message)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create message with headers
	msg := nats.NewMsg(enum.EventWebtrackerSessionNew.String())
	msg.Data = data
	msg.Header.Set("X-Tenant", utils.GetTenantFromContext(ctx))
	msg.Header.Set("X-UserId", utils.GetUserIdFromContext(ctx))

	// Publish to the stored subject
	_, err = s.natsConn.JS.PublishMsg(msg)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to publish page view event: %w", err)
	}

	return nil
}
