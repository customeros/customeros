package events_platform

import (
	"context"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type MockInvoiceServiceCallbacks struct {
	NextPreviewInvoiceForContract func(context.Context, *invoicepb.NextPreviewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error)
	NewInvoiceForContract         func(context.Context, *invoicepb.NewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error)
	SimulateInvoice               func(context.Context, *invoicepb.SimulateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
	UpdateInvoice                 func(context.Context, *invoicepb.UpdateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
	VoidInvoice                   func(context.Context, *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
}

var invoiceCallbacks = &MockInvoiceServiceCallbacks{}

func SetInvoiceCallbacks(callbacks *MockInvoiceServiceCallbacks) {
	invoiceCallbacks = callbacks
}

type MockInvoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
}

func (MockInvoiceService) NextPreviewInvoiceForContract(context context.Context, proto *invoicepb.NextPreviewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.NextPreviewInvoiceForContract == nil {
		panic("invoiceCallbacks.NextPreviewInvoiceForContract is not set")
	}
	return invoiceCallbacks.NextPreviewInvoiceForContract(context, proto)
}

func (MockInvoiceService) NewInvoiceForContract(context context.Context, proto *invoicepb.NewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.NewInvoiceForContract == nil {
		panic("invoiceCallbacks.NewInvoiceForContract is not set")
	}
	return invoiceCallbacks.NewInvoiceForContract(context, proto)
}

func (MockInvoiceService) SimulateInvoice(context context.Context, proto *invoicepb.SimulateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.SimulateInvoice == nil {
		panic("invoiceCallbacks.SimulateInvoice is not set")
	}
	return invoiceCallbacks.SimulateInvoice(context, proto)
}

func (MockInvoiceService) UpdateInvoice(context context.Context, proto *invoicepb.UpdateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.UpdateInvoice == nil {
		panic("invoiceCallbacks.UpdateInvoice is not set")
	}
	return invoiceCallbacks.UpdateInvoice(context, proto)
}

func (MockInvoiceService) VoidInvoice(context context.Context, proto *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.VoidInvoice == nil {
		panic("invoiceCallbacks.VoidInvoice is not set")
	}
	return invoiceCallbacks.VoidInvoice(context, proto)
}
