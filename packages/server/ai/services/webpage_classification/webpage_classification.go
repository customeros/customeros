package webpage_classification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
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
	MODEL_TEMPERATURE = 0.2
	MAX_TOKENS        = 1024
	RETRIES           = 2
)

type WebpageClassification struct {
	PrimaryTopic        string           `json:"primary_topic"`
	SecondaryTopics     []string         `json:"secondary_topics"`
	SolutionFocus       []string         `json:"solution_focus"`
	ContentType         enum.ContentType `json:"content_type"`
	IndustryVertical    string           `json:"industry_vertical"`
	KeyPainPoints       []string         `json:"key_pain_points"`
	ValueProposition    string           `json:"value_proposition"`
	ReferencedCustomers []string         `json:"referenced_customers"`
}

func (s *webpageClassification) handleWebpageClassification(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageClassification.handleWebpageClassification")
	defer span.Finish()

	message, err := s.parseNatsMessage(ctx, msg)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// askAI
	aiRequest := buildAIRequest(message)
	resp, err := s.askAI.AskAI(ctx, aiRequest)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// parse response
	webpageClassification, err := s.parseWebpageClassification(ctx, resp)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// validate response
	err = s.validateWebpageClassification(ctx, webpageClassification)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// publish response message
	err = s.publishWebpageClassifiedEvent(ctx, message.ContentId, webpageClassification)
	if err != nil {
		span.TraceError(err)
		return err
	}

	return nil
}

func (s *webpageClassification) publishWebpageClassifiedEvent(ctx context.Context, contentID string, record *WebpageClassification) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageClassification.publishWebpageClassifiedEvent")
	defer span.Finish()

	event := &pb.WebpageClassified{
		ContentId:           contentID,
		PrimaryTopic:        record.PrimaryTopic,
		SecondaryTopics:     record.SecondaryTopics,
		SolutionFocus:       record.SolutionFocus,
		ContentType:         string(record.ContentType),
		IndustryVertical:    record.IndustryVertical,
		KeyPainPoints:       record.KeyPainPoints,
		ValueProposition:    record.ValueProposition,
		ReferencedCustomers: record.ReferencedCustomers,
	}

	data, err := proto.Marshal(event)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create message with headers
	msg := nats.NewMsg(enum.EventWebpageClassified.String())
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

func (s *webpageClassification) validateWebpageClassification(ctx context.Context, record *WebpageClassification) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageClassification.validateWebpageClassification")
	defer span.Finish()

	if record == nil {
		err := ai.ErrEmptyResponse
		span.TraceError(err)
		return err
	}

	return nil
}

func (s *webpageClassification) parseWebpageClassification(ctx context.Context, response *string) (*WebpageClassification, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageClassification.parseWebpageClassification")
	defer span.Finish()

	if response == nil {
		span.TraceError(ai.ErrEmptyResponse)
		return nil, ai.ErrEmptyResponse
	}

	resp := &WebpageClassification{}
	err := json.Unmarshal([]byte(*response), &resp)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	return resp, nil
}

func (s *webpageClassification) parseNatsMessage(ctx context.Context, msg *nats.Msg) (*pb.RequestWebpageClassification, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageClassification.parseNatsMessage")
	defer span.Finish()

	message := &pb.RequestWebpageClassification{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	switch {
	case message.ContentId == "":
		return nil, errors.New("ContentID cannot be empty")
	case message.Content == "":
		return nil, errors.New("Content cannot be empty")
	default:
		return message, nil
	}
}

func buildAIRequest(message *pb.RequestWebpageClassification) *interfaces.AskAIRequest {
	systemPrompt := `I'm will provide you with the scraped content of a webpage along with some metadata about the company it belongs to.  Your job is to classify the content based on:
	- Primary Topic
	- Secondary Topics (0-3 values)
	- The Solution(s) the content is focused on (1-3 values)
	- The type of content it is (e.g. case study, product page, documentation, ect)
	- The Industry Vertical the content is targeting
	- The Key Pain Points the content aims to address (1-3 values)
	- The core Value Proposition of the content
	- A list of all Referenced Customers contained within the content
	- Valid content types must match one of the following:
	    - article
	    - whitepaper
	    - webinar
	    - case study
	    - product page
	    - solution page
	    - testimonial
	    - research report
	    - technical documentation

	IMPORTANT: Your response MUST be in valid JSON format exactly maching this schema:
	{
  "primary_topic": "Cloud Migration Strategy",
  "secondary_topics": [
    "Digital Transformation",
    "Infrastructure Modernization",
    "DevOps Adoption"
  ],
  "solution_focus": [
    "AWS Migration Services",
    "Container Management"
  ],
  "content_type": "Technical Whitepaper",
  "industry_vertical": "Financial Services",
  "key_pain_points": [
    "Legacy system maintenance costs",
    "Scalability limitations",
    "Security compliance requirements",
    "Slow deployment cycles"
  ],
  "value_proposition": "Reduce operational costs by 40% while improving deployment speed",
  "referenced_customers": [
    "JP Morgan Chase",
    "Bank of America",
    "Wells Fargo",
    "Capital One"
  ]
}
	Do not include any text outside the JSON object.  If you are unable to confidently determine a value, return an empty string.`

	var p strings.Builder

	fmt.Fprintf(&p, "Domain: %s\n", message.Domain)
	fmt.Fprintf(&p, "Website URL: %s\n", message.Url)
	fmt.Fprintf(&p, "URL Content: %s\n", message.Content)
	prompt := p.String()

	askAI := &interfaces.AskAIRequest{
		Model:            enum.AIModelAnthropicSonnet,
		SystemPrompt:     &systemPrompt,
		Prompt:           &prompt,
		ModelTemperature: utils.Float32Ptr(MODEL_TEMPERATURE),
		MaxOutputTokens:  utils.Int32Ptr(MAX_TOKENS),
		OutputFormat:     enum.AIOutputJson,
		Retries:          utils.IntPtr(RETRIES),
	}

	return askAI
}
