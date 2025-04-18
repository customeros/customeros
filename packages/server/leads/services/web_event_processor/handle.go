package web_event_processor

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

const REQUEST_TIMEOUT = 60 * time.Second

type LeadSource struct {
	Channel          enum.Channel
	Source           enum.LeadSource
	SourcePlatform   enum.SocialPlatform
	ViewedOnPlatform enum.SocialPlatform
	GCLID            string
	DeviceType       enum.DeviceType
	LandingPage      string
	ReferrerDomain   string
	ReferrerUrl      string
	UTMSource        string
	UTMMedium        string
	UTMCampaign      string
	UTMContent       string
	UTMTerm          string
	Language         string
}

func (s *WebEventProcessor) processNewSession(ctx context.Context, message *pb.WebTrackerEvent) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.processNewSession")
	defer spans.Finish()

	// attempt to identify company
	sessionIdentity, err := s.identifyWebSession(ctx, message.Ip)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if sessionIdentity == nil {
		return nil
	}

	// determime lead source
	leadSource, err := s.determineLeadSource(ctx, message.Referrer, message.UserAgent)

	// lookup company to determine if new or existing lead

	// publish new or exisitng lead event
	// lead.new.webtracker
	// lead.existing.webtracker

	return nil
}

func (s *WebEventProcessor) processIdentifyEvent(ctx context.Context, message *pb.WebTrackerEvent) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.processIdentifyEvent")
	defer spans.Finish()

	return nil
}

func (s *WebEventProcessor) determineLeadSource(ctx context.Context, referrer, userAgent string) (string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.determineLeadSource")
	defer spans.Finish()

	if referrer == "" {
		return "", nil
	}

	return nil
}

func (s *WebEventProcessor) identifyWebSession(ctx context.Context, ipAddress string) (*pb.IdentifyVisitorResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.identifyWebSession")
	defer spans.Finish()

	if ipAddress == "" {
		return nil, nil
	}

	request := &pb.IdentifyVisitorRequest{}

	resp, err := s.sendIdentifyVisitorRequest(ctx, request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if resp == nil {
		err = errors.New("Unable to identify web session")
		spans.TraceError(err)
		return nil, err
	}

	return resp, nil
}

func (s *WebEventProcessor) sendIdentifyVisitorRequest(ctx context.Context, request *pb.IdentifyVisitorRequest) (*pb.IdentifyVisitorResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.sendIdentifyVisitorRequest")
	defer spans.Finish()

	// Marshal request to protobuf
	reqData, err := proto.Marshal(request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Send request to service
	msg := nats.NewMsg(enum.EventIdentifyVisitor.String())
	msg.Header = nats.Header{
		enum.TENANT_HEADER:  []string{utils.GetTenantFromContext(ctx)},
		enum.USER_ID_HEADER: []string{utils.GetUserIdFromContext(ctx)},
	}
	msg.Data = reqData

	resp, err := s.natsConn.Conn.RequestMsg(msg, REQUEST_TIMEOUT)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Unmarshal response
	response := &pb.IdentifyVisitorResponse{}
	if err := proto.Unmarshal(resp.Data, response); err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return response, nil
}
