package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Types and Constants
type Spans struct {
	Jaeger opentracing.Span
	OTel   trace.Span
}

type contextKey string

// Component tag constants
const (
	ComponentGraphQL  = "graphql"
	ComponentPostgres = "postgres"
	ComponentNeo4j    = "neo4j"
	ComponentREST     = "rest"
	ComponentService  = "service"
	ComponentListener = "listener"
	ComponentCronJob  = "cron"
)

const (
	componentKey = "component"
)

// SpanKind constants
const (
	SpanKindInternal = "internal"
	SpanKindServer   = "server"
	SpanKindClient   = "client"
	SpanKindProducer = "producer"
	SpanKindConsumer = "consumer"
)

// Span tag constants
const (
	SpanTagTenant         = "tenant"
	SpanTagUserId         = "user.id"
	SpanTagUserEmail      = "user.email"
	SpanTagEntityId       = "entity.id"
	SpanTagEventType      = "event.type"
	SpanTagAppSource      = "app.source"
	SpanTagExternalSystem = "external.system"
	SpanTagExternalId     = "external.id"
)

// Context keys
const (
	otelSpanKey contextKey = "otel_span"
)

// SpanOptions defines options for span creation
type SpanOptions struct {
	NewRoot bool
}

// WithForceNewTrace returns a SpanOptions that forces creation of a new trace
func WithNewRoot() SpanOptions {
	return SpanOptions{
		NewRoot: true,
	}
}

// GetDefaultServiceSpanAttributes returns default attributes for service spans
func getDefaultServiceSpanAttributes(ctx context.Context) []attribute.KeyValue {
	attrs := []attribute.KeyValue{}

	if tenant := common.GetTenantFromContext(ctx); tenant != "" {
		attrs = append(attrs, attribute.String(SpanTagTenant, tenant))
	}
	if userID := common.GetUserIdFromContext(ctx); userID != "" {
		attrs = append(attrs, attribute.String(SpanTagUserId, userID))
	}
	if userEmail := common.GetUserEmailFromContext(ctx); userEmail != "" {
		attrs = append(attrs, attribute.String(SpanTagUserEmail, userEmail))
	}
	if appSource := common.GetAppSourceFromContext(ctx); appSource != "" {
		attrs = append(attrs, attribute.String(SpanTagAppSource, appSource))
	}

	return attrs
}

// SetDefaultServiceSpanAttributes sets default attributes on an OpenTelemetry span
func setDefaultServiceSpanAttributes(ctx context.Context, span trace.Span) {
	if span == nil {
		return
	}
	span.SetAttributes(getDefaultServiceSpanAttributes(ctx)...)
}

// SetDefaultJaegerSpanTags sets default tags on a Jaeger span
func setDefaultJaegerSpanTags(ctx context.Context, span opentracing.Span) {
	if span == nil {
		return
	}
	if tenant := common.GetTenantFromContext(ctx); tenant != "" {
		span.SetTag(SpanTagTenant, tenant)
	}
	if userID := common.GetUserIdFromContext(ctx); userID != "" {
		span.SetTag(SpanTagUserId, userID)
	}
	if userEmail := common.GetUserEmailFromContext(ctx); userEmail != "" {
		span.SetTag(SpanTagUserEmail, userEmail)
	}
	if appSource := common.GetAppSourceFromContext(ctx); appSource != "" {
		span.SetTag(SpanTagAppSource, appSource)
	}
}

// Core Span Operations
func StartSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	// Start Jaeger span
	var jaegerSpan opentracing.Span
	if len(opts) > 0 && opts[0].NewRoot {
		// Force new trace by creating a new root span
		jaegerSpan = opentracing.StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, jaegerSpan)
	} else {
		jaegerSpan, ctx = opentracing.StartSpanFromContext(ctx, operationName)
	}

	setDefaultJaegerSpanTags(ctx, jaegerSpan)

	// Start OpenTelemetry span
	tracer := otel.Tracer("github.com/customeros/customeros")
	var otelSpan trace.Span
	if len(opts) > 0 && opts[0].NewRoot {
		// Force new trace by using WithNewRoot() option while keeping the context
		ctx, otelSpan = tracer.Start(ctx, operationName, trace.WithNewRoot())
	} else {
		ctx, otelSpan = tracer.Start(ctx, operationName)
	}
	otelSpan.SetAttributes(attribute.String("operation.name", operationName))

	setDefaultServiceSpanAttributes(ctx, otelSpan)

	// Store spans in context
	ctx = context.WithValue(ctx, otelSpanKey, otelSpan)

	return &Spans{
		Jaeger: jaegerSpan,
		OTel:   otelSpan,
	}, ctx
}

// Finish ends both Jaeger and OpenTelemetry spans
func (s *Spans) Finish() {
	if s == nil {
		return
	}
	FinishSpans(s)
}

// FinishSpans ends both Jaeger and OpenTelemetry spans
func FinishSpans(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.Finish()
	}
	if spans.OTel != nil {
		spans.OTel.End()
	}
}

// Component-specific Span Starters
func StartCronSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentCronJob(spans)
	SetSpanKindInternal(spans)
	return spans, ctx
}

func StartGraphQLSpan(ctx context.Context, operationName string, gqlCtx *graphql.OperationContext, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentGraphQL(spans)
	SetSpanKindServer(spans)

	// Add GraphQL operation details if available
	if gqlCtx != nil {
		if spans.Jaeger != nil {
			spans.Jaeger.SetTag("graphql.operation.name", gqlCtx.OperationName)
			spans.Jaeger.SetTag("graphql.operation.type", string(gqlCtx.Operation.Operation))
		}
		if spans.OTel != nil {
			spans.OTel.SetAttributes(
				attribute.String("graphql.operation.name", gqlCtx.OperationName),
				attribute.String("graphql.operation.type", string(gqlCtx.Operation.Operation)),
			)
		}
	}

	return spans, ctx
}

func StartPostgresSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentPostgres(spans)
	SetSpanKindDatabase(spans)
	return spans, ctx
}

func StartNeo4jSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentNeo4j(spans)
	SetSpanKindDatabase(spans)
	return spans, ctx
}

func StartServiceSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentService(spans)
	SetSpanKindInternal(spans)
	return spans, ctx
}

func StartRestSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentREST(spans)
	SetSpanKindServer(spans)
	return spans, ctx
}

func StartProducerSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentService(spans)
	SetSpanKindProducer(spans)
	return spans, ctx
}

func StartListenerSpan(ctx context.Context, operationName string, opts ...SpanOptions) (*Spans, context.Context) {
	spans, ctx := StartSpan(ctx, operationName, opts...)
	TagComponentListener(spans)
	SetSpanKindConsumer(spans)
	return spans, ctx
}

// Component Tagging Helpers
func TagComponentGraphQL(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.SetTag(componentKey, ComponentGraphQL)
	}
	if spans.OTel != nil {
		spans.OTel.SetAttributes(attribute.String(componentKey, ComponentGraphQL))
	}
}

func TagComponentPostgres(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.SetTag(componentKey, ComponentPostgres)
	}
	if spans.OTel != nil {
		spans.OTel.SetAttributes(attribute.String(componentKey, ComponentPostgres))
	}
}

func TagComponentNeo4j(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.SetTag(componentKey, ComponentNeo4j)
	}
	if spans.OTel != nil {
		spans.OTel.SetAttributes(attribute.String(componentKey, ComponentNeo4j))
	}
}

func TagComponentREST(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.SetTag(componentKey, ComponentREST)
	}
	if spans.OTel != nil {
		spans.OTel.SetAttributes(attribute.String(componentKey, ComponentREST))
	}
}

func TagComponentService(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.SetTag(componentKey, ComponentService)
	}
	if spans.OTel != nil {
		spans.OTel.SetAttributes(attribute.String(componentKey, ComponentService))
	}
}

func TagComponentListener(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.SetTag(componentKey, ComponentListener)
	}
	if spans.OTel != nil {
		spans.OTel.SetAttributes(attribute.String(componentKey, ComponentListener))
	}
}

func TagComponentCronJob(spans *Spans) {
	if spans == nil {
		return
	}
	if spans.Jaeger != nil {
		spans.Jaeger.SetTag(componentKey, ComponentCronJob)
	}
	if spans.OTel != nil {
		spans.OTel.SetAttributes(attribute.String(componentKey, ComponentCronJob))
	}
}

// Span Kind Helpers
func SetSpanKindInternal(spans *Spans) {
	if spans == nil || spans.OTel == nil {
		return
	}
	spans.OTel.SetAttributes(attribute.String("span.kind", SpanKindInternal))
}

func SetSpanKindServer(spans *Spans) {
	if spans == nil || spans.OTel == nil {
		return
	}
	spans.OTel.SetAttributes(attribute.String("span.kind", SpanKindServer))
}

func SetSpanKindClient(spans *Spans) {
	if spans == nil || spans.OTel == nil {
		return
	}
	spans.OTel.SetAttributes(attribute.String("span.kind", SpanKindClient))
}

func SetSpanKindProducer(spans *Spans) {
	if spans == nil || spans.OTel == nil {
		return
	}
	spans.OTel.SetAttributes(attribute.String("span.kind", SpanKindProducer))
}

func SetSpanKindConsumer(spans *Spans) {
	if spans == nil || spans.OTel == nil {
		return
	}
	spans.OTel.SetAttributes(attribute.String("span.kind", SpanKindConsumer))
}

func SetSpanKindDatabase(spans *Spans) {
	if spans == nil || spans.OTel == nil {
		return
	}
	spans.OTel.SetAttributes(attribute.String("span.kind", SpanKindClient))
}

// Tagging Methods
func (s *Spans) TagString(key, value string) {
	if s == nil {
		return
	}
	if key == "" {
		return
	}
	if s.Jaeger != nil {
		s.Jaeger.SetTag(key, value)
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.String(key, value))
	}
}

func (s *Spans) TagInt(key string, value int) {
	if s == nil {
		return
	}
	if key == "" {
		return
	}
	if s.Jaeger != nil {
		s.Jaeger.SetTag(key, value)
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.Int(key, value))
	}
}

func (s *Spans) TagUint32(key string, value uint32) {
	if s == nil {
		return
	}
	if key == "" {
		return
	}
	if s.Jaeger != nil {
		s.Jaeger.SetTag(key, value)
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.Int64(key, int64(value)))
	}
}

func (s *Spans) TagBool(key string, value bool) {
	if s == nil {
		return
	}
	if key == "" {
		return
	}
	if s.Jaeger != nil {
		s.Jaeger.SetTag(key, value)
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.Bool(key, value))
	}
}

func (s *Spans) TagStringSlice(key string, value []string) {
	if s == nil {
		return
	}
	if key == "" {
		return
	}
	if s.Jaeger != nil {
		s.Jaeger.SetTag(key, fmt.Sprintf("%v", value))
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.StringSlice(key, value))
	}
}

func (s *Spans) TagEntity(entityId string) {
	if s == nil || entityId == "" {
		return
	}
	if s.Jaeger != nil {
		tracing.TagEntity(s.Jaeger, entityId)
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.String(SpanTagEntityId, entityId))
	}
}

func (s *Spans) TagEventType(eventType string) {
	if s == nil || eventType == "" {
		return
	}
	if s.Jaeger != nil {
		s.Jaeger.SetTag(SpanTagEventType, eventType)
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.String(SpanTagEventType, eventType))
	}
}

func (s *Spans) TagTenant(tenant string) {
	if s == nil || tenant == "" {
		return
	}
	if s.Jaeger != nil {
		tracing.TagTenant(s.Jaeger, tenant)
	}
	if s.OTel != nil {
		s.OTel.SetAttributes(attribute.String(SpanTagTenant, tenant))
	}
}

// Logging Methods
func (s *Spans) LogFields(fields ...opentracinglog.Field) {
	if s == nil {
		return
	}

	// Log to Jaeger
	if s.Jaeger != nil {
		s.Jaeger.LogFields(fields...)
	}

	// Log to OpenTelemetry as events
	if s.OTel != nil {
		// Group fields by event type
		eventFields := make(map[string][]attribute.KeyValue)
		eventType := "log" // default event type

		for _, field := range fields {
			key := field.Key()
			value := field.Value()

			// Check if this is an event type field
			if key == "event" {
				if str, ok := value.(string); ok {
					eventType = str
				}
				continue
			}

			// Convert field to attribute
			var attr attribute.KeyValue
			switch v := value.(type) {
			case string:
				attr = attribute.String(key, v)
			case int:
				attr = attribute.Int(key, v)
			case bool:
				attr = attribute.Bool(key, v)
			case float64:
				attr = attribute.Float64(key, v)
			case error:
				attr = attribute.String(key, v.Error())
			default:
				attr = attribute.String(key, fmt.Sprintf("%v", v))
			}
			eventFields[eventType] = append(eventFields[eventType], attr)
		}

		// Add each event type as a separate event
		for eventType, attrs := range eventFields {
			s.OTel.AddEvent(eventType, trace.WithAttributes(attrs...))
		}
	}
}

func (s *Spans) LogKV(alternatingKeyValues ...interface{}) {
	if s == nil {
		return
	}

	// Log to Jaeger
	if s.Jaeger != nil {
		s.Jaeger.LogKV(alternatingKeyValues...)
	}

	// Log to OpenTelemetry
	if s.OTel != nil {
		attrs := make([]attribute.KeyValue, 0, len(alternatingKeyValues)/2)
		for i := 0; i < len(alternatingKeyValues); i += 2 {
			if i+1 >= len(alternatingKeyValues) {
				break
			}
			key, ok := alternatingKeyValues[i].(string)
			if !ok {
				continue
			}
			value := alternatingKeyValues[i+1]
			switch v := value.(type) {
			case string:
				attrs = append(attrs, attribute.String(key, v))
			case int:
				attrs = append(attrs, attribute.Int(key, v))
			case bool:
				attrs = append(attrs, attribute.Bool(key, v))
			case float64:
				attrs = append(attrs, attribute.Float64(key, v))
			case error:
				attrs = append(attrs, attribute.String(key, v.Error()))
			default:
				// For unknown types, try to marshal as JSON
				if jsonBytes, err := json.Marshal(v); err == nil {
					attrs = append(attrs, attribute.String(key, string(jsonBytes)))
				} else {
					attrs = append(attrs, attribute.String(key, fmt.Sprintf("%v", v)))
				}
			}
		}
		s.OTel.AddEvent("log", trace.WithAttributes(attrs...))
	}
}

func (s *Spans) LogObjectAsJson(key string, obj interface{}) {
	if s == nil {
		return
	}

	// Log to Jaeger
	if s.Jaeger != nil {
		tracing.LogObjectAsJson(s.Jaeger, key, obj)
	}

	// Log to OpenTelemetry
	if s.OTel != nil {
		jsonBytes, err := json.Marshal(obj)
		if err != nil {
			s.OTel.AddEvent("log", trace.WithAttributes(
				attribute.String("error", fmt.Sprintf("failed to marshal object to JSON: %v", err)),
			))
			return
		}
		s.OTel.AddEvent("log", trace.WithAttributes(
			attribute.String(key, string(jsonBytes)),
		))
	}
}

func LogError(ctx context.Context, err error, fields ...opentracinglog.Field) {
	// Log to Jaeger
	jaegerSpan := opentracing.SpanFromContext(ctx)
	if jaegerSpan != nil {
		// Add standard error fields
		errorFields := []opentracinglog.Field{
			opentracinglog.Error(err),
			opentracinglog.String("event", "error"),
			opentracinglog.String("time", time.Now().Format(time.RFC3339)),
		}

		// Add OpenTelemetry specific fields
		otelFields := []opentracinglog.Field{
			opentracinglog.String("otel.status_code", "error"),
			opentracinglog.String("otel.status_message", err.Error()),
		}

		// Combine all fields
		allFields := append(errorFields, append(otelFields, fields...)...)
		jaegerSpan.LogFields(allFields...)
	}

	// Log to OpenTelemetry
	if otelSpan, ok := ctx.Value(otelSpanKey).(trace.Span); ok {
		otelSpan.RecordError(err)
		otelSpan.SetStatus(codes.Error, err.Error())
		// Add event with message and time
		otelSpan.AddEvent("error", trace.WithAttributes(
			attribute.String("message", err.Error()),
			attribute.String("time", time.Now().Format(time.RFC3339)),
		))
	}
}

func LogInfo(ctx context.Context, msg string, fields ...opentracinglog.Field) {
	// Log to Jaeger
	jaegerSpan := opentracing.SpanFromContext(ctx)
	if jaegerSpan != nil {
		// Add standard info fields
		infoFields := []opentracinglog.Field{
			opentracinglog.String("event", "info"),
			opentracinglog.String("message", msg),
			opentracinglog.String("time", time.Now().Format(time.RFC3339)),
		}

		// Add OpenTelemetry specific fields
		otelFields := []opentracinglog.Field{
			opentracinglog.String("otel.status_code", "ok"),
		}

		// Combine all fields
		allFields := append(infoFields, append(otelFields, fields...)...)
		jaegerSpan.LogFields(allFields...)
	}

	// Log to OpenTelemetry
	if otelSpan, ok := ctx.Value(otelSpanKey).(trace.Span); ok {
		otelSpan.SetStatus(codes.Ok, msg)
		// Add event with message and time
		otelSpan.AddEvent("info", trace.WithAttributes(
			attribute.String("message", msg),
			attribute.String("time", time.Now().Format(time.RFC3339)),
		))
	}
}

func LogDebug(ctx context.Context, msg string, fields ...opentracinglog.Field) {
	// Log to Jaeger
	jaegerSpan := opentracing.SpanFromContext(ctx)
	if jaegerSpan != nil {
		// Add standard debug fields
		debugFields := []opentracinglog.Field{
			opentracinglog.String("event", "debug"),
			opentracinglog.String("message", msg),
			opentracinglog.String("time", time.Now().Format(time.RFC3339)),
		}

		// Add OpenTelemetry specific fields
		otelFields := []opentracinglog.Field{
			opentracinglog.String("otel.status_code", "ok"),
		}

		// Combine all fields
		allFields := append(debugFields, append(otelFields, fields...)...)
		jaegerSpan.LogFields(allFields...)
	}

	// Log to OpenTelemetry
	if otelSpan, ok := ctx.Value(otelSpanKey).(trace.Span); ok {
		// Add event with message and time
		otelSpan.AddEvent("debug", trace.WithAttributes(
			attribute.String("message", msg),
			attribute.String("time", time.Now().Format(time.RFC3339)),
		))
	}
}

// Error Handling
func TagError(span opentracing.Span, err error) {
	if err != nil {
		// Add standard error tags
		span.SetTag("error", true)

		// Add OpenTelemetry specific tags
		span.SetTag("otel.status_code", "error")
		span.SetTag("otel.status_message", err.Error())

		span.LogFields(
			opentracinglog.Error(err),
			opentracinglog.String("event", "error"),
			opentracinglog.String("time", time.Now().Format(time.RFC3339)),
		)
	}
}

func (s *Spans) TraceError(err error) {
	if s == nil || err == nil {
		return
	}

	// Trace error in Jaeger
	if s.Jaeger != nil {
		tracing.TraceErr(s.Jaeger, err)
	}

	// Trace error in OpenTelemetry
	if s.OTel != nil {
		s.OTel.RecordError(err)
		s.OTel.SetStatus(codes.Error, err.Error())
		s.OTel.SetAttributes(
			attribute.String("event", "error"),
			attribute.String("time", time.Now().Format(time.RFC3339)),
		)
	}
}

// Recovery
func RecoverAndLog(ctx context.Context, spans *Spans, logger logger.Logger) {
	if r := recover(); r != nil {
		stack := string(debug.Stack())

		// Log to OpenTelemetry
		if spans != nil && spans.OTel != nil {
			spans.OTel.RecordError(fmt.Errorf("panic: %v", r))
			spans.OTel.SetStatus(codes.Error, fmt.Sprintf("panic: %v", r))
			spans.OTel.SetAttributes(
				attribute.String("event", "panic"),
				attribute.String("error", fmt.Sprintf("%v", r)),
				attribute.String("time", time.Now().Format(time.RFC3339)),
			)
			// Log stack trace as an event
			spans.OTel.AddEvent("panic.stack", trace.WithAttributes(
				attribute.String("stack", stack),
			))
		}

		// Log to Jaeger
		if spans != nil && spans.Jaeger != nil {
			spans.Jaeger.SetTag("error", true)
			spans.Jaeger.SetTag("event", "panic")
			spans.Jaeger.LogFields(
				opentracinglog.Error(fmt.Errorf("panic: %v", r)),
				opentracinglog.String("event", "panic"),
				opentracinglog.String("time", time.Now().Format(time.RFC3339)),
				opentracinglog.String("stack", stack),
			)
		}

		// Log to logger
		if logger != nil {
			logger.Errorf("Recovered from panic: %v\nStack trace:\n%s", r, stack)
		}

		// Do not re-panic - allow the application to continue running
	}
}

// RecoverAndLogMain handles panic recovery and logs to both Jaeger and OpenTelemetry for main function
func RecoverAndLogMain(appLogger logger.Logger) {
	if r := recover(); r != nil {
		stackTrace := string(debug.Stack())

		// Log to Jaeger
		tracing.RecoverAndLogToJaeger(appLogger)

		// Log to OpenTelemetry
		tracer := otel.Tracer("github.com/customeros/customeros")
		_, span := tracer.Start(context.Background(), "panic-recovery")
		defer span.End()

		span.SetStatus(codes.Error, fmt.Sprintf("panic: %v", r))
		span.SetAttributes(attribute.String("event", "panic"))

		// Log error and stack as events
		span.AddEvent("error", trace.WithAttributes(
			attribute.String("error", fmt.Sprintf("%v", r)),
		))
		span.AddEvent("stack", trace.WithAttributes(
			attribute.String("stack", stackTrace),
		))

		appLogger.Errorf("Recovered from panic: %v\nStack trace:\n%s", r, stackTrace)
	}
}

// InjectSpanContextIntoHTTPRequest injects both Jaeger and OpenTelemetry span contexts into an HTTP request
func InjectSpanContextIntoHTTPRequest(req *http.Request, spans *Spans) *http.Request {
	if spans == nil {
		return req
	}

	// Inject Jaeger span context
	if spans.Jaeger != nil {
		// Use existing tracing package's function for Jaeger
		req = tracing.InjectSpanContextIntoHTTPRequest(req, spans.Jaeger)
	}

	// Inject OpenTelemetry span context
	if spans.OTel != nil {
		// Get the propagator from the global tracer provider
		propagator := otel.GetTextMapPropagator()

		// Create a carrier for the headers
		carrier := propagation.HeaderCarrier(req.Header)

		// Create a context with the span context
		ctx := trace.ContextWithSpanContext(context.Background(), spans.OTel.SpanContext())

		// Inject the span context into the carrier
		propagator.Inject(ctx, carrier)
	}

	return req
}

// StartHttpServerTracerSpanWithHeader starts both Jaeger and OpenTelemetry spans for HTTP server requests
func StartHttpServerTracerSpanWithHeader(ctx context.Context, operationName string, headers http.Header) (*Spans, context.Context) {
	// Start Jaeger span using existing implementation
	ctx, jaegerSpan := tracing.StartHttpServerTracerSpanWithHeader(ctx, operationName, headers)

	// Start OpenTelemetry span
	tracer := otel.Tracer("github.com/customeros/customeros")

	// Extract context from headers for OpenTelemetry
	carrier := propagation.HeaderCarrier(headers)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	// Start the span
	ctx, otelSpan := tracer.Start(ctx, operationName)

	// Create spans struct
	spans := &Spans{
		Jaeger: jaegerSpan,
		OTel:   otelSpan,
	}

	// Log headers for both tracers
	for key, values := range headers {
		// Log to Jaeger
		if jaegerSpan != nil {
			jaegerSpan.LogFields(
				opentracinglog.String("request.header.key", key),
				opentracinglog.Object("request.header.value", values),
			)
		}
		// Log to OpenTelemetry
		if otelSpan != nil {
			otelSpan.AddEvent("request.header."+key, trace.WithAttributes(
				attribute.StringSlice("value", values),
			))
		}
	}

	// Store OpenTelemetry span in context
	ctx = context.WithValue(ctx, otelSpanKey, otelSpan)

	return spans, ctx
}

// GraphQLTracingEnhancer creates a middleware that adds tracing for GraphQL operations
func GraphQLTracingEnhancer(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		operationName := "GraphQL." + tracing.ExtractGraphQLMethodName(c.Request)

		// Start both Jaeger and OpenTelemetry spans
		spans, ctxWithSpan := StartHttpServerTracerSpanWithHeader(
			ctx,
			operationName,
			c.Request.Header,
		)
		defer spans.Finish()

		// Tag as GraphQL component
		TagComponentGraphQL(spans)
		SetSpanKindServer(spans)

		// Update request context
		c.Request = c.Request.WithContext(ctxWithSpan)
		c.Next()

		// Add response status
		if c.Writer.Status() >= 400 {
			spans.TraceError(nil)
		}
	}
}

func GetTraceIds(spans *Spans) (string, string) {
	if spans == nil {
		return "", ""
	}
	jaegerTraceId := ""
	otelTraceId := ""
	if spans.Jaeger != nil {
		tracingData := ExtractTextMapCarrier((spans.Jaeger).Context())
		jaegerTraceId = strings.Split(tracingData["uber-trace-id"], ":")[0]
	}
	if spans.OTel != nil {
		otelTraceId = spans.OTel.SpanContext().TraceID().String()
	}
	return jaegerTraceId, otelTraceId
}

func TraceIdsAsString(spans *Spans) string {
	if spans == nil {
		return ""
	}
	jaegerTraceId, otelTraceId := GetTraceIds(spans)
	return fmt.Sprintf("Jaeger Trace ID: %s, OpenTelemetry Trace ID: %s", jaegerTraceId, otelTraceId)
}

// EnrichCtxWithSpanCtxForGraphQL enriches the provided context with both Jaeger and OpenTelemetry span contexts
func EnrichCtxWithSpanCtxForGraphQL(ctx context.Context, operationContext *graphql.OperationContext) context.Context {
	// Extract and enrich Jaeger span context
	EnrichCtxWithJaegerSpanForGraphQL(ctx, operationContext)

	// Extract and enrich OpenTelemetry span context
	carrier := propagation.HeaderCarrier(operationContext.Headers)
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	return ctx
}

// TraceErrorOnActiveSpan traces an error on both active Jaeger and OpenTelemetry spans from context
func TraceErrorOnActiveSpan(ctx context.Context, err error) {
	if err == nil {
		return
	}

	// Trace error in Jaeger using the active span
	if span := opentracing.SpanFromContext(ctx); span != nil {
		tracing.TraceErr(span, err)
	}

	// Trace error in OpenTelemetry using the active span
	if span := trace.SpanFromContext(ctx); span.IsRecording() {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(
			attribute.String("event", "error"),
			attribute.String("time", time.Now().Format(time.RFC3339)),
		)
	}
}

// RecoveryWithTelemetry creates a gin middleware that recovers from panics and logs them to both Jaeger and OpenTelemetry
func RecoveryWithTelemetry(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Get the current spans from context or create new ones
				var spans *Spans

				// Try to get existing spans from context
				if existingSpan := opentracing.SpanFromContext(c.Request.Context()); existingSpan != nil {
					// Create spans struct with existing Jaeger span
					spans = &Spans{
						Jaeger: existingSpan,
					}
				} else {
					// Create new spans if none exist
					spans, _ = StartSpan(context.Background(), "panic-recovery", WithNewRoot())
				}

				// Get detailed stack trace
				stack := string(debug.Stack())

				// Log to OpenTelemetry
				if spans.OTel != nil {
					spans.OTel.RecordError(fmt.Errorf("panic: %v", r))
					spans.OTel.SetStatus(codes.Error, fmt.Sprintf("panic: %v", r))
					spans.OTel.SetAttributes(
						attribute.String("event", "panic"),
						attribute.String("error", fmt.Sprintf("%v", r)),
						attribute.String("time", time.Now().Format(time.RFC3339)),
						attribute.String("http.url", c.Request.URL.String()),
						attribute.String("http.method", c.Request.Method),
						attribute.Int("http.status_code", 500),
					)
					// Log stack trace as an event
					spans.OTel.AddEvent("panic.stack", trace.WithAttributes(
						attribute.String("stack", stack),
					))
				}

				// Log to Jaeger
				if spans.Jaeger != nil {
					spans.Jaeger.SetTag("error", true)
					spans.Jaeger.SetTag("event", "panic")
					spans.Jaeger.LogFields(
						opentracinglog.Error(fmt.Errorf("panic: %v", r)),
						opentracinglog.String("event", "panic"),
						opentracinglog.String("time", time.Now().Format(time.RFC3339)),
						opentracinglog.String("stack", stack),
						opentracinglog.String("http.url", c.Request.URL.String()),
						opentracinglog.String("http.method", c.Request.Method),
					)
					// Set HTTP tags
					ext.HTTPUrl.Set(spans.Jaeger, c.Request.URL.String())
					ext.HTTPMethod.Set(spans.Jaeger, c.Request.Method)
					ext.HTTPStatusCode.Set(spans.Jaeger, 500)
				}

				// Log to application logger
				log.Errorf("[Panic Recovery] Error: %v\nStack trace:\n%s\nPath: %s\nMethod: %s",
					r,
					stack,
					c.Request.URL.Path,
					c.Request.Method,
				)

				// Finish spans
				spans.Finish()

				// Let the chain continue to allow other recovery handlers to process the panic
				panic(r)
			}
		}()
		c.Next()
	}
}
