package session_manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

const (
	REQUEST_TIMEOUT         = 60 * time.Second
	IP_DATA_LOOKBACK_PERIOD = 90 // days
)

type LeadSource struct {
	Channel        enum.Channel
	SearchEngine   enum.SearchEngine
	SocialPlatform enum.SocialPlatform
	AdPlatform     enum.AdPlatform
	ReferrerHost   string
	ReferrerPath   string
	Campaign       *CampaignMetadata
}

type CampaignMetadata struct {
	IsPaid      bool
	AdPlatform  enum.AdPlatform
	UTMSource   string
	UTMMedium   string
	UTMCampaign string
	UTMTerm     string
	UTMContent  string
}

func (s *sessionManager) NewSession(ctx context.Context, msg *nats.Msg) {
	span, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.NewSession")
	defer span.Finish()

	if msg == nil {
		span.TraceError(errors.New("nil nats message"))
		return
	}
	span.TagString("nats.subject", msg.Subject)
	span.TagString("nats.reply", msg.Reply)

	message := &pb.WebtrackerSessionNew{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		err := fmt.Errorf("failed to parse message: %w", err)
		span.TraceError(err)
		s.handleProcessingError(ctx, msg, err)
		return
	}

	s.processNewSession(ctx, message)
}

func (s *sessionManager) processNewSession(ctx context.Context, message *pb.WebtrackerSessionNew) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.processNewSession")
	defer spans.Finish()

	// check if known IP
	ipRecord, err := s.checkForExistingIPRecord(ctx, message.Ip)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	// if known IP, return
	if ipRecord != nil {
		return nil
	}

	// process new IP
	ipRecord, err = s.processNewIP(ctx, message)
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	if ipRecord == nil || ipRecord.Domain == "" {
		return nil
	}

	// publish visitor.identified event
	event := &pb.WebtrackerVisitorIdentified{
		SessionId: message.SessionId,
		TrackerId: message.TrackerId,
		VisitorId: message.VisitorId,
		Ip:        ipRecord.IPAddress,
		Domain:    ipRecord.Domain,
	}
	return s.createVisitorIdentifiedEvent(ctx, event)
}

func (s *sessionManager) createVisitorIdentifiedEvent(ctx context.Context, event *pb.WebtrackerVisitorIdentified) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.publishVisitorIdentifiedEvent")
	defer span.Finish()

	payload, err := proto.Marshal(event)
	if err != nil {
		span.TraceError(err)
		return err
	}

	outbox := &models.OutboxEvent{
		ID:        utils.GenerateEventID(),
		EventType: enum.EventWebtrackerVisitorIdentified,
		EntityID:  event.TrackerId,
		Publisher: enum.SessionManager,
		Tenant:    utils.GetTenantFromContext(ctx),
		SessionID: event.SessionId,
		Payload:   payload,
		Status:    enum.OutboxPending,
		CreatedAt: utils.Now(),
	}

	return s.repositories.Outbox.Create(ctx, outbox)
}

func (s *sessionManager) checkForExistingIPRecord(ctx context.Context, ipAddress string) (*models.IPIntelligence, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.checkForExistingIPRecord")
	defer span.Finish()

	ipRecord, err := s.repositories.IPIntelligence.FindByIP(ctx, ipAddress)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	cutoffDate := time.Now().AddDate(0, 0, IP_DATA_LOOKBACK_PERIOD)
	if ipRecord != nil && ipRecord.UpdatedAt.Before(cutoffDate) {
		return nil, nil
	}

	return ipRecord, nil
}

func (s *sessionManager) processNewIP(ctx context.Context, message *pb.WebtrackerSessionNew) (*models.IPIntelligence, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.processNewIP")
	defer span.Finish()

	// check if bot, return early if not trusted IP
	userAgent := utils.ParseUserAgent(message.UserAgent)
	isSuspicious := utils.IsSuspiciousURL(message.Referrer)

	ipProfile, err := s.profileIP(ctx, message.Ip)
	if err != nil {
		span.TraceError(err)
		return nil, nil
	}

	if ipProfile != nil {
		if ipProfile.IsThreat || userAgent.IsBot || isSuspicious {
			return nil, nil
		}
	}

	// attempt to identify
	domain, err := s.identifyIP(ctx, message.Ip)
	if err != nil {
		span.TraceError(err)
		return nil, nil
	}

	if domain == "" {
		return nil, nil
	}

	// get latest intel record from DB
	return s.repositories.IPIntelligence.FindByIP(ctx, message.Ip)
}

func (s *sessionManager) identifyIP(ctx context.Context, ipAddress string) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.identifyIP")
	defer span.Finish()

	request := &pb.IPAddressIdentifyRequest{
		IpAddress: ipAddress,
	}

	reqData, err := proto.Marshal(request)
	if err != nil {
		span.TraceError(err)
		return "", err
	}

	// Send request to service
	msg := nats.NewMsg(enum.EventAskSnitcher.String())
	msg.Header = nats.Header{
		enum.TENANT_HEADER:  []string{utils.GetTenantFromContext(ctx)},
		enum.USER_ID_HEADER: []string{utils.GetUserIdFromContext(ctx)},
	}
	msg.Data = reqData

	resp, err := s.natsConn.Conn.RequestMsg(msg, REQUEST_TIMEOUT)
	if err != nil {
		span.TraceError(err)
		return "", err
	}

	// Unmarshal response
	response := &pb.IPAddressIdentifyResponse{}
	if err := proto.Unmarshal(resp.Data, response); err != nil {
		span.TraceError(err)
		return "", err
	}

	if response.ErrorMessage != "" {
		span.TraceError(errors.New(response.ErrorMessage))
		return "", err
	}

	return response.Domain, nil
}

func (s *sessionManager) profileIP(ctx context.Context, ipAddress string) (*pb.IPAddressVerifyResponse, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.isTrustedIP")
	defer span.Finish()

	request := &pb.IPAddressVerifyRequest{
		IpAddress: ipAddress,
	}

	reqData, err := proto.Marshal(request)
	if err != nil {
		span.TraceError(err)
		return nil, err
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
		return nil, err
	}

	// Unmarshal response
	response := &pb.IPAddressVerifyResponse{}
	if err := proto.Unmarshal(resp.Data, response); err != nil {
		span.TraceError(err)
		return nil, err
	}

	return response, nil
}
