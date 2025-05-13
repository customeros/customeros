package webpage_intent_profiler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/enum"
	nats_internal "github.com/customeros/customeros/packages/server/ai/internal/nats"
	"github.com/customeros/customeros/packages/server/ai/internal/proto/pb"
	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
	"github.com/customeros/customeros/packages/server/ai/internal/utils"
	ai "github.com/customeros/customeros/packages/server/ai/services/ask_ai"
)

const (
	MODEL             = enum.AIModelAnthropicSonnet
	MODEL_TEMPERATURE = 0.2
	MAX_TOKENS        = 1024
	RETRIES           = 2
)

var ErrInvalidScore = errors.New("Invalid score")

type WebpageIntentProfile struct {
	ProblemRecognitionScore int `json:"problem_recognition"`
	SolutionResearchScore   int `json:"solution_research"`
	EvaluationScore         int `json:"evaluation"`
	PurchaseReadinessScrore int `json:"purchase_readiness"`
}

func (s *webpageIntentProfiler) handleWebpageIntentRequest(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageIntentProfiler.handleWebpageIntentRequest")
	defer span.Finish()

	message, err := s.parseNatsMessage(ctx, msg)
	if err != nil {
		span.TraceError(err)
		return err
	}

	err = s.validateNatsMessage(ctx, message)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// askAI
	aiRequest := buildAIRequest(message)
	resp, err := s.ai.AskAI(ctx, aiRequest)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// parse response
	scores, err := s.parseIntentScores(ctx, resp)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// validate response
	err = s.validateScores(ctx, scores)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// publish response message
	return s.publishWebpageProfiledEvent(ctx, message.ContentId, scores)
}

func (s *webpageIntentProfiler) publishWebpageProfiledEvent(ctx context.Context, contentID string, record *WebpageIntentProfile) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageIntentProfiler.publishWebpageProfiledEvent")
	defer span.Finish()

	event := &pb.WebpageProfiled{
		ContentId:               contentID,
		ProblemRecognitionScore: uint32(record.ProblemRecognitionScore),
		SolutionResearchScore:   uint32(record.SolutionResearchScore),
		EvaluationScore:         uint32(record.EvaluationScore),
		PurchaseReadinessScore:  uint32(record.PurchaseReadinessScrore),
	}

	data, err := proto.Marshal(event)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create message with headers
	msg := nats.NewMsg(enum.EventWebpageProfiled.String())
	msg.Data = data
	msg.Header.Set(nats_internal.HEADER_TENANT, utils.GetTenantFromContext(ctx))
	msg.Header.Set(nats_internal.HEADER_USERID, utils.GetUserIdFromContext(ctx))

	// Publish to the stored subject
	_, err = s.natsConn.JS.PublishMsg(msg)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to publish stored email: %w", err)
	}
	return nil
}

func (s *webpageIntentProfiler) validateScores(ctx context.Context, record *WebpageIntentProfile) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageIntentProfiler.validateScores")
	defer span.Finish()

	if record == nil {
		err := ai.ErrEmptyResponse
		span.TraceError(err)
		return err
	}

	if record.ProblemRecognitionScore < 0 || record.ProblemRecognitionScore > 5 {
		err := errors.Wrap(ErrInvalidScore, "Problem Recognition Score")
		span.TraceError(err)
		return err
	}

	if record.SolutionResearchScore < 0 || record.SolutionResearchScore > 5 {
		err := errors.Wrap(ErrInvalidScore, "Solution Research Score")
		span.TraceError(err)
		return err
	}

	if record.EvaluationScore < 0 || record.EvaluationScore > 5 {
		err := errors.Wrap(ErrInvalidScore, "Evaluation Score")
		span.TraceError(err)
		return err
	}

	if record.PurchaseReadinessScrore < 0 || record.PurchaseReadinessScrore > 5 {
		err := errors.Wrap(ErrInvalidScore, "Purchase Readiness Score")
		span.TraceError(err)
		return err
	}

	return nil
}

func (s *webpageIntentProfiler) parseIntentScores(ctx context.Context, resp *string) (*WebpageIntentProfile, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageIntentProfiler.parseIntentScores")
	defer span.Finish()

	if resp == nil {
		span.TraceError(ai.ErrEmptyResponse)
		return nil, ai.ErrEmptyResponse
	}

	response := &WebpageIntentProfile{}
	err := json.Unmarshal([]byte(*resp), response)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	return response, nil
}

func (s *webpageIntentProfiler) validateNatsMessage(ctx context.Context, message *pb.RequestWebpageIntent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageIntentProfiler.validateNatsMessage")
	defer span.Finish()

	switch {
	case message.ContentId == "":
		return errors.New("ContentID cannot be empty")
	case message.Content == "":
		return errors.New("Content is empty")
	default:
		return nil
	}
}

func (s *webpageIntentProfiler) parseNatsMessage(ctx context.Context, msg *nats.Msg) (*pb.RequestWebpageIntent, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageIntentProfiler.handleWebpageIntentRequest")
	defer span.Finish()

	message := &pb.RequestWebpageIntent{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	return message, nil
}

func buildAIRequest(message *pb.RequestWebpageIntent) *interfaces.AskAIRequest {
	systemPrompt := `I will provide you with the scraped content of a webpage along with company metadata.

Your job is to score the content across 4 stages of the buyer's journey, evaluating how well it addresses buyer needs at each stage:

1. Problem Recognition (identifying challenges/pain points)
2. Solution Research (educating about approaches/methodologies)
3. Evaluation (comparing options, features, case studies)
4. Purchase Readiness (pricing, demos, trials, CTAs)

Score each stage 1-5:
	1: Not relevant for this stage
    2: Slightly relevant
    3: Moderately relevant
    4: Very relevant
    5: Highly relevant

Return ONLY a JSON object in this exact format:
{
  "problem_recognition": 3,
  "solution_research": 5,
  "evaluation": 2,
  "purchase_readiness": 1
}`

	var p strings.Builder

	fmt.Fprintf(&p, "Domain: %s\n", message.Domain)
	fmt.Fprintf(&p, "Website URL: %s\n", message.Url)
	fmt.Fprintf(&p, "URL Content: %s\n", message.Content)
	prompt := p.String()

	askAI := &interfaces.AskAIRequest{
		Model:            MODEL,
		SystemPrompt:     &systemPrompt,
		Prompt:           &prompt,
		ModelTemperature: utils.Float32Ptr(MODEL_TEMPERATURE),
		MaxOutputTokens:  utils.Int32Ptr(MAX_TOKENS),
		OutputFormat:     enum.AIOutputJson,
		Retries:          utils.IntPtr(RETRIES),
	}

	return askAI
}
