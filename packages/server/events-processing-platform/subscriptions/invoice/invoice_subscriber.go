package invoice

import (
	"context"
	fsc "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/file_store_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoice"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
	"github.com/opentracing/opentracing-go"
	"io/ioutil"
	"os"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/subscriptions"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"golang.org/x/sync/errgroup"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type InvoiceSubscriber struct {
	log                 logger.Logger
	db                  *esdb.Client
	cfg                 *config.Config
	grpcClients         *grpc_client.Clients
	fsc                 fsc.FileStoreApiService
	invoiceEventHandler *InvoiceEventHandler
}

func NewInvoiceSubscriber(log logger.Logger, db *esdb.Client, cfg *config.Config, repositories *repository.Repositories, grpcClients *grpc_client.Clients, fsc fsc.FileStoreApiService) *InvoiceSubscriber {
	return &InvoiceSubscriber{
		log:                 log,
		db:                  db,
		cfg:                 cfg,
		grpcClients:         grpcClients,
		fsc:                 fsc,
		invoiceEventHandler: NewInvoiceEventHandler(log, repositories, &cfg.EventNotifications, grpcClients),
	}
}

func (s *InvoiceSubscriber) Connect(ctx context.Context, worker subscriptions.Worker) error {
	group, ctx := errgroup.WithContext(ctx)
	for i := 1; i <= s.cfg.Subscriptions.InvoiceSubscription.PoolSize; i++ {
		sub, err := s.db.SubscribeToPersistentSubscriptionToAll(
			ctx,
			s.cfg.Subscriptions.InvoiceSubscription.GroupName,
			esdb.SubscribeToPersistentSubscriptionOptions{
				BufferSize: s.cfg.Subscriptions.InvoiceSubscription.BufferSizeClient,
			},
		)
		if err != nil {
			return err
		}
		defer sub.Close()

		group.Go(s.runWorker(ctx, worker, sub, i))
	}
	return group.Wait()
}

func (consumer *InvoiceSubscriber) runWorker(ctx context.Context, worker subscriptions.Worker, stream *esdb.PersistentSubscription, i int) func() error {
	return func() error {
		return worker(ctx, stream, i)
	}
}

func (s *InvoiceSubscriber) ProcessEvents(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error {

	for {
		event := stream.Recv()
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if event.SubscriptionDropped != nil {
			s.log.Errorf("(SubscriptionDropped) err: {%v}", event.SubscriptionDropped.Error)
			return errors.Wrap(event.SubscriptionDropped.Error, "Subscription Dropped")
		}

		if event.EventAppeared != nil {
			s.log.EventAppeared(s.cfg.Subscriptions.InvoiceSubscription.GroupName, event.EventAppeared.Event, workerID)

			if event.EventAppeared.Event.Event == nil {
				s.log.Errorf("(InvoiceSubscriber) event.EventAppeared.Event.Event is nil")
			} else {
				err := s.When(ctx, eventstore.NewEventFromRecorded(event.EventAppeared.Event.Event))
				if err != nil {
					s.log.Errorf("(InvoiceSubscription.when) err: {%v}", err)

					if err := stream.Nack(err.Error(), esdb.NackActionPark, event.EventAppeared.Event); err != nil {
						s.log.Errorf("(stream.Nack) err: {%v}", err)
						return errors.Wrap(err, "stream.Nack")
					}
				}
			}

			err := stream.Ack(event.EventAppeared.Event)
			if err != nil {
				s.log.Errorf("(stream.Ack) err: {%v}", err)
				return errors.Wrap(err, "stream.Ack")
			}
			s.log.Debugf("(ACK) event: {%+v}", eventstore.NewRecordedBaseEventFromRecorded(event.EventAppeared.Event.Event))
		}
	}
}

func (s *InvoiceSubscriber) When(ctx context.Context, evt eventstore.Event) error {
	ctx, span := tracing.StartProjectionTracerSpan(ctx, "InvoiceSubscriber.When", evt)
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()), log.String("EventType", evt.GetEventType()))

	if strings.HasPrefix(evt.GetAggregateID(), "$") {
		return nil
	}

	if s.cfg.Subscriptions.InvoiceSubscription.IgnoreEvents {
		return nil
	}

	switch evt.GetEventType() {
	case invoice.InvoiceNewV1:
		return s.invoiceEventHandler.onInvoiceNewV1(ctx, evt)
	case invoice.InvoiceFillV1:
		return s.onInvoiceFillV1(ctx, evt)
	case invoice.InvoicePdfGeneratedV1:
		return s.invoiceEventHandler.onInvoicePdfGeneratedV1(ctx, evt)
	default:
		return nil
	}

	return nil
}

type InvoiceData struct {
	CustomerName  string
	ProviderName  string
	InvoiceNumber string
	// Add more fields as needed
}

func (s *InvoiceSubscriber) onInvoiceFillV1(ctx context.Context, evt eventstore.Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceSubscriber.onInvoiceFillV1")
	defer span.Finish()
	span.LogFields(log.String("AggregateID", evt.GetAggregateID()))

	var eventData invoice.InvoiceFillEvent
	if err := evt.GetJsonData(&eventData); err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "evt.GetJsonData")
	}

	//TODO build invoice PDF
	//load tenant billing details ( logo, address, etc ) from neo4j
	//load billing profile for organization from neo4j
	//load invoice from neo4j for invoice date and due date
	//take the data from InvoiceFillEvent to fill the invoice

	// Define your invoice data
	data := map[string]interface{}{
		"CustomerName":          "John Doe",
		"CustomerAddress":       "12 Stardust Street Skyline, SK 98765, nited States, invoices@robertinc.com",
		"ProviderLogoUrl":       "https://static.wikia.nocookie.net/fictionalcompanies/images/c/c2/ACME_Corporation.png",
		"ProviderLogoExtension": "",
		"ProviderName":          "XYZ Company",
		"ProviderAddress":       "29 Maple Lane, Springfield, Haven County, San Francisco • 89302, United States",
		"InvoiceNumber":         "123456",
		"InvoiceIssueDate":      "22.01.2024",
		"InvoiceDueDate":        "22.01.2024",
		"InvoiceCurrency":       "US$",
		"InvoiceSubtotal":       "123.00",
		"InvoiceTax":            "0.00",
		"InvoiceTotal":          "123.00",
		"InvoicePaid":           "50.00",
		"InvoiceAmountDue":      "73.00",
	}
	data["ProviderLogoExtension"] = GetFileExtensionFromUrl(data["ProviderLogoUrl"].(string))

	//prepare the temp html file
	tmpInvoiceFile, err := os.CreateTemp("", "invoice_*.html")
	if err != nil {
		return errors.Wrap(err, "ioutil.TempFile")
	}
	defer os.Remove(tmpInvoiceFile.Name()) // Delete the temporary HTML file when done
	defer tmpInvoiceFile.Close()

	//fill the template with data and store it in temp
	err = FillInvoiceHtmlTemplate(tmpInvoiceFile, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "FillInvoiceHtmlTemplate")
	}

	//convert the temp to pdf
	pdfBytes, err := ConvertInvoiceHtmlToPdf(s.cfg.Subscriptions.InvoiceSubscription.PdfConverterUrl, tmpInvoiceFile, data)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "ConvertInvoiceHtmlToPdf")
	}

	// Save the PDF file to disk
	err = ioutil.WriteFile("output.pdf", *pdfBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "ioutil.WriteFile")
	}

	fileDTO, err := s.fsc.UploadSingleFileBytes(eventData.Tenant, *pdfBytes)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.UploadSingleFileBytes")
	}

	if fileDTO.Id == "" {
		return errors.New("fileDTO.Id is empty")
	}

	err = s.CallPdfGeneratedInvoice(ctx, eventData.Tenant, evt.GetAggregateID(), fileDTO.Id, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceFillV1.CallPdfGeneratedInvoice")
	}

	return nil
}

func (s *InvoiceSubscriber) CallPdfGeneratedInvoice(ctx context.Context, tenant, invoiceId, repositoryFileId string, span opentracing.Span) error {
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := s.grpcClients.InvoiceClient.PdfGeneratedInvoice(ctx, &invoicepb.PdfGeneratedInvoiceRequest{
		Tenant:           tenant,
		InvoiceId:        invoiceId,
		RepositoryFileId: repositoryFileId,
		AppSource:        constants.AppSourceEventProcessingPlatform,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error sending the pdf generated request for invoice %s: %s", invoiceId, err.Error())
		return err
	}
	return nil
}
