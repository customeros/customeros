package mocked_grpc

import (
	"context"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type MockInvoiceServiceCallbacks struct {
	GenerateInvoicePdf func(ctx context.Context, proto *invoicepb.GenerateInvoicePdfRequest) (*invoicepb.InvoiceIdResponse, error)
	RequestFillInvoice func(ctx context.Context, proto *invoicepb.RequestFillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
	FillInvoice        func(ctx context.Context, proto *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
}

var InvoiceCallbacks = &MockInvoiceServiceCallbacks{}

func SetInvoiceCallbacks(callbacks *MockInvoiceServiceCallbacks) {
	InvoiceCallbacks = callbacks
}

type MockInvoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
}

func (MockInvoiceService) GenerateInvoicePdf(ctx context.Context, proto *invoicepb.GenerateInvoicePdfRequest) (*invoicepb.InvoiceIdResponse, error) {
	if InvoiceCallbacks.GenerateInvoicePdf == nil {
		panic("InvoiceCallbacks.GenerateInvoicePdf is not set")
	}
	return InvoiceCallbacks.GenerateInvoicePdf(ctx, proto)
}

func (MockInvoiceService) RequestFillInvoice(ctx context.Context, proto *invoicepb.RequestFillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if InvoiceCallbacks.RequestFillInvoice == nil {
		panic("InvoiceCallbacks.RequestFillInvoice is not set")
	}
	return InvoiceCallbacks.RequestFillInvoice(ctx, proto)
}

func (MockInvoiceService) FillInvoice(ctx context.Context, proto *invoicepb.FillInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if InvoiceCallbacks.FillInvoice == nil {
		panic("InvoiceCallbacks.FillInvoice is not set")
	}
	return InvoiceCallbacks.FillInvoice(ctx, proto)
}
