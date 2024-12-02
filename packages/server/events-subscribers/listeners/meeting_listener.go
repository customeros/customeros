package listeners

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func OnMeetingSummaryCreated(ctx context.Context, services *service.Services, input any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.OnMeetingSummaryCreated")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	message := input.(*dto.Event)
	// check message data type before conversion
	if message.Event.Data == nil {
		err := errors.New("message data is nil")
		tracing.TraceErr(span, err)
		return nil
	}
	messageData := message.Event.Data.(*dto.MeetingSummaryCreated)

	// do smth with message data
	fmt.Println(messageData)

	return nil
}
