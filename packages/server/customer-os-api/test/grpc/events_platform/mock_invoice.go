package events_platform

import (
	"context"
	invoicepb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type MockInvoiceServiceCallbacks struct {
	NewInvoiceForContract func(context.Context, *invoicepb.NewInvoiceForContractRequest) (*invoicepb.InvoiceIdResponse, error)
	SimulateInvoice       func(context.Context, *invoicepb.SimulateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
	UpdateInvoice         func(context.Context, *invoicepb.UpdateInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
	PayInvoice            func(context.Context, *invoicepb.PayInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
	VoidInvoice           func(context.Context, *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error)
}

var invoiceCallbacks = &MockInvoiceServiceCallbacks{}

func SetInvoiceCallbacks(callbacks *MockInvoiceServiceCallbacks) {
	invoiceCallbacks = callbacks
}

type MockInvoiceService struct {
	invoicepb.UnimplementedInvoiceGrpcServiceServer
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

func (MockInvoiceService) PayInvoice(context context.Context, proto *invoicepb.PayInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.PayInvoice == nil {
		panic("invoiceCallbacks.PayInvoice is not set")
	}
	return invoiceCallbacks.PayInvoice(context, proto)
}

func (MockInvoiceService) VoidInvoice(context context.Context, proto *invoicepb.VoidInvoiceRequest) (*invoicepb.InvoiceIdResponse, error) {
	if invoiceCallbacks.VoidInvoice == nil {
		panic("invoiceCallbacks.VoidInvoice is not set")
	}
	return invoiceCallbacks.VoidInvoice(context, proto)
}
