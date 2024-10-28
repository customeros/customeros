package tracing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
)

const (
	SpanTagTenant                = "tenant"
	SpanTagUserId                = "user-id"
	SpanTagUserEmail             = "user-email"
	SpanTagEntityId              = "entity-id"
	SpanTagComponent             = "component"
	SpanTagExternalSystem        = "external-system"
	SpanTagAggregateId           = "aggregateID"
	SpanTagRedundantEventSkipped = "redundantEventSkipped"
)

const (
	SpanTagComponentPostgresRepository = "postgresRepository"
	SpanTagComponentNeo4jRepository    = "neo4jRepository"
	SpanTagComponentRest               = "rest"
	SpanTagComponentCronJob            = "cronJob"
)

func GraphQlTracingEnhancer(ctx context.Context) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxWithSpan, span := StartHttpServerTracerSpanWithHeader(ctx, ExtractGraphQLMethodName(c.Request), c.Request.Header)
		for k, v := range c.Request.Header {
			span.LogFields(log.String("request.header.key", k), log.Object("request.header.value", v))
		}
		defer span.Finish()
		TagComponentRest(span)
		c.Request = c.Request.WithContext(ctxWithSpan)
		c.Next()
	}
}

func TracingEnhancer(ctx context.Context, endpoint string) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctxWithSpan, span := StartHttpServerTracerSpanWithHeader(ctx, endpoint, c.Request.Header)
		for k, v := range c.Request.Header {
			span.LogFields(log.String("request.header.key", k), log.Object("request.header.value", v))
		}
		defer span.Finish()
		TagComponentRest(span)
		c.Request = c.Request.WithContext(ctxWithSpan)
		c.Next()
	}
}

func StartHttpServerTracerSpanWithHeader(ctx context.Context, operationName string, headers http.Header) (context.Context, opentracing.Span) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))

	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		opentracing.GlobalTracer().Inject(serverSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
		return opentracing.ContextWithSpan(ctx, serverSpan), serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	return opentracing.ContextWithSpan(ctx, serverSpan), serverSpan
}

func StartRabbitMQMessageTracerSpanWithHeader(ctx context.Context, operationName string, uberTraceId string) (context.Context, opentracing.Span) {
	textMapCarrierFromMetaData := make(opentracing.TextMapCarrier)
	textMapCarrierFromMetaData.Set("uber-trace-id", uberTraceId)

	span, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, textMapCarrierFromMetaData)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(span))
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)
	return ctx, serverSpan

}

func StartTracerSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
	return serverSpan, opentracing.ContextWithSpan(ctx, serverSpan)
}

func InjectSpanContextIntoGrpcMetadata(ctx context.Context, span opentracing.Span) context.Context {
	if span != nil {
		// Inject the span context into the gRPC request metadata.
		textMapCarrier := make(opentracing.TextMapCarrier)
		err := span.Tracer().Inject(span.Context(), opentracing.TextMap, textMapCarrier)
		if err == nil {
			// Add the injected metadata to the gRPC context.
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.New(nil)
			}
			for key, val := range textMapCarrier {
				md.Set(key, val)
			}
			ctx = metadata.NewOutgoingContext(ctx, md)
			return ctx
		}
	}
	return ctx
}

func InjectSpanContextIntoHTTPRequest(req *http.Request, span opentracing.Span) *http.Request {
	if span != nil {
		// Prepare to inject span context into HTTP headers
		tracer := span.Tracer()
		textMapCarrier := opentracing.HTTPHeadersCarrier(req.Header)

		// Inject the span context into the HTTP headers
		err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, textMapCarrier)
		if err != nil {
			// Log error or handle it as per the application's error handling strategy
			fmt.Println("Error injecting span context into headers:", err)
		}
	}
	return req
}

func InjectSpanContextIntoGraphQLRequest(req *graphql.Request, span opentracing.Span) {
	if span != nil {
		carrier := make(opentracing.TextMapCarrier)
		err := span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier)
		if err != nil {
			// Log error or handle it as per your application's error handling strategy
			fmt.Println("Error injecting span context into GraphQL request:", err)
			return
		}

		for k, v := range carrier {
			req.Header.Set(k, v)
		}
	}
}

func ExtractGraphQLMethodName(req *http.Request) string {
	// Read the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		// Handle error
		return ""
	}

	// Restore the request body
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	// Parse the request body as JSON
	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		// Handle error
		return ""
	}

	// Extract the method name from the GraphQL request
	if operationName, ok := requestBody["operationName"].(string); ok {
		return operationName
	}

	// If the method name is not found, you can add additional logic here to extract it from the request body or headers if applicable
	// ...
	return ""
}

func setDefaultSpanTags(ctx context.Context, span opentracing.Span) {
	tenant := common.GetTenantFromContext(ctx)
	loggedInUserId := common.GetUserIdFromContext(ctx)
	loggedInUserEmail := common.GetUserEmailFromContext(ctx)
	if tenant != "" {
		span.SetTag(SpanTagTenant, tenant)
	}
	if loggedInUserId != "" {
		span.SetTag(SpanTagUserId, loggedInUserId)
	}
	if loggedInUserEmail != "" {
		span.SetTag(SpanTagUserEmail, loggedInUserEmail)
	}
}

func SetDefaultServiceSpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	span.SetTag(SpanTagComponent, constants.ComponentService)
}
func SetDefaultNeo4jRepositorySpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	TagComponentNeo4jRepository(span)
}
func SetDefaultPostgresRepositorySpanTags(ctx context.Context, span opentracing.Span) {
	setDefaultSpanTags(ctx, span)
	TagComponentNeo4jRepository(span)
}

func TraceErr(span opentracing.Span, err error, fields ...log.Field) {
	// Log the error with the fields
	if span == nil {
		return
	}
	ext.LogError(span, err, fields...)
}

func LogObjectAsJson(span opentracing.Span, name string, object any) {
	if object == nil {
		span.LogFields(log.String(name, "nil"))
	}
	jsonObject, err := json.Marshal(object)
	if err == nil {
		span.LogFields(log.String(name, string(jsonObject)))
	} else {
		span.LogFields(log.Object(name, object))
	}
}

func InjectTextMapCarrier(spanCtx opentracing.SpanContext) (opentracing.TextMapCarrier, error) {
	m := make(opentracing.TextMapCarrier)
	if err := opentracing.GlobalTracer().Inject(spanCtx, opentracing.TextMap, m); err != nil {
		return nil, err
	}
	return m, nil
}

func ExtractTextMapCarrier(spanCtx opentracing.SpanContext) opentracing.TextMapCarrier {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return make(opentracing.TextMapCarrier)
	}
	return textMapCarrier
}

func TagComponentPostgresRepository(span opentracing.Span) {
	span.SetTag(SpanTagComponent, SpanTagComponentPostgresRepository)
}

func TagComponentNeo4jRepository(span opentracing.Span) {
	span.SetTag(SpanTagComponent, SpanTagComponentNeo4jRepository)
}

func TagTenant(span opentracing.Span, tenant string) {
	if tenant != "" {
		span.SetTag(SpanTagTenant, tenant)
	}
}

func TagEntity(span opentracing.Span, entityId string) {
	if entityId != "" {
		span.SetTag(SpanTagEntityId, entityId)
	}
}

func TagComponentCronJob(span opentracing.Span) {
	span.SetTag(SpanTagComponent, SpanTagComponentCronJob)
}

func TagComponentRest(span opentracing.Span) {
	span.SetTag(SpanTagComponent, SpanTagComponentRest)
}
