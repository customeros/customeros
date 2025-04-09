package event_logger

import (
	"context"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/mailstack/internal/enum"
	"github.com/customeros/mailstack/internal/telemetry"
	pb_mappers "github.com/customeros/mailstack/proto/mappers"
	"github.com/customeros/mailstack/proto/pb"
)

func (s *EventLoggerService) processErrorMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processErrorMessage")
	defer spans.Finish()

	message := &pb.ErrorEvent{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EmailEvent(msg.Subject)
	event.Publisher = pb_mappers.ServiceNameToMailstackService(message.Publisher)
	event.EmailID = message.EmailId
	event.HasError = true
	event.ErrorMessage = message.ErrorMessage
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processSkipInboundProcessing(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processReceivedIMAPMessage")
	defer spans.Finish()

	message := &pb.SkipInboundProcessing{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundClassifiedSkip
	event.Publisher = enum.MailstackContentService
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data
	event.EmailID = message.EmailId
	event.MailboxID = message.MailboxId
	event.Classification = pb_mappers.PbToEmailClassification(message.Classification)

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processReceivedIMAPMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processReceivedIMAPMessage")
	defer spans.Finish()

	message := &pb.EmailReceivedIMAP{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundStored
	event.Publisher = enum.MailstackIMAPService
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processStoredMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processStoredMessage")
	defer spans.Finish()

	message := &pb.EmailStored{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundStored
	event.Publisher = enum.MailstackStorageService
	event.EmailID = message.EmailId
	event.MailboxID = message.MailboxId
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processClassificationMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processClassificationMessage")
	defer spans.Finish()

	message := &pb.EmailClassificationRequest{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundClassify
	event.Publisher = enum.MailstackContentService
	event.EmailID = message.EmailId
	event.MailboxID = message.MailboxId
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processAnalysisMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processAnalysisMessage")
	defer spans.Finish()

	message := &pb.AnalyzeEmailRequest{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundAnalysis
	event.Publisher = enum.MailstackContentService
	event.EmailID = message.EmailId
	event.MailboxID = message.MailboxId
	event.Classification = pb_mappers.PbToEmailClassification(message.Classification)
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processAttachmentsMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processAttachmentsMessage")
	defer spans.Finish()

	message := &pb.ProcessAttachmentRequest{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundAttachments
	event.Publisher = enum.MailstackContentService
	event.EmailID = message.EmailId
	event.MailboxID = message.MailboxId
	event.Classification = pb_mappers.PbToEmailClassification(message.Classification)
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processThreadMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processThreadMessage")
	defer spans.Finish()

	message := &pb.AttachToThreadRequest{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundThread
	event.Publisher = enum.MailstackContentService
	event.EmailID = message.EmailId
	event.MailboxID = message.MailboxId
	event.Classification = pb_mappers.PbToEmailClassification(message.Classification)
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}

func (s *EventLoggerService) processInboundCompletedMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventLoggerService.processInboundCompletedMessage")
	defer spans.Finish()

	message := &pb.InboundEmailProcessingCompleted{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	event := s.NewEmailEventRecord(ctx)
	event.Event = enum.EventEmailInboundCompleted
	event.Publisher = enum.MailstackContentService
	event.EmailID = message.EmailId
	event.MailboxID = message.MailboxId
	event.Classification = pb_mappers.PbToEmailClassification(message.Classification)
	event.ThreadID = message.ThreadId
	event.Direction = enum.EmailDirectionInbound
	event.Payload = msg.Data

	err = s.repositories.EmailEventRepository.Create(ctx, event)
	if err != nil {
		spans.TraceError(err)
	}
	return
}
