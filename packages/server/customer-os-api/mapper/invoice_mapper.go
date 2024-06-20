package mapper

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

func MapEntityToInvoice(entity *neo4jentity.InvoiceEntity) *model.Invoice {
	if entity == nil {
		return nil
	}
	invoice := model.Invoice{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
			Version:       entity.AggregateVersion,
		},
		DryRun:             entity.DryRun,
		Postpaid:           entity.Postpaid,
		OffCycle:           entity.OffCycle,
		Preview:            entity.Preview,
		InvoiceNumber:      entity.Number,
		InvoicePeriodStart: entity.PeriodStartDate,
		InvoicePeriodEnd:   entity.PeriodEndDate,
		Due:                entity.DueDate,
		Issued:             entity.IssuedDate,
		AmountDue:          entity.TotalAmount,
		TaxDue:             entity.Vat,
		Subtotal:           entity.Amount,
		Currency:           entity.Currency.String(),
		RepositoryFileID:   entity.RepositoryFileId,
		InvoiceURL:         fmt.Sprintf(constants.UrlFileStoreFileDownloadUrlTemplate, entity.RepositoryFileId),
		Status:             utils.ToPtr(mapper.MapInvoiceStatusToModel(entity.Status)),
		Note:               utils.StringPtrNillable(entity.Note),
		PaymentLink:        utils.StringPtrNillable(entity.PaymentDetails.PaymentLink),
		Customer: &model.InvoiceCustomer{
			Name:            utils.StringPtrNillable(entity.Customer.Name),
			Email:           utils.StringPtrNillable(entity.Customer.Email),
			AddressLine1:    utils.StringPtrNillable(entity.Customer.AddressLine1),
			AddressLine2:    utils.StringPtrNillable(entity.Customer.AddressLine2),
			AddressZip:      utils.StringPtrNillable(entity.Customer.Zip),
			AddressLocality: utils.StringPtrNillable(entity.Customer.Locality),
			AddressCountry:  utils.StringPtrNillable(entity.Customer.Country),
			AddressRegion:   utils.StringPtrNillable(entity.Customer.Region),
		},
		Provider: &model.InvoiceProvider{
			LogoRepositoryFileID: utils.StringPtrNillable(entity.Provider.LogoRepositoryFileId),
			Name:                 utils.StringPtrNillable(entity.Provider.Name),
			AddressLine1:         utils.StringPtrNillable(entity.Provider.AddressLine1),
			AddressLine2:         utils.StringPtrNillable(entity.Provider.AddressLine2),
			AddressZip:           utils.StringPtrNillable(entity.Provider.Zip),
			AddressLocality:      utils.StringPtrNillable(entity.Provider.Locality),
			AddressCountry:       utils.StringPtrNillable(entity.Provider.Country),
			AddressRegion:        utils.StringPtrNillable(entity.Provider.Region),
		},
	}
	if entity.Status == neo4jenum.InvoiceStatusPaid {
		invoice.Paid = true
		invoice.AmountRemaining = float64(0)
		invoice.AmountPaid = entity.TotalAmount
	} else {
		invoice.Paid = false
		invoice.AmountRemaining = entity.TotalAmount
		invoice.AmountPaid = float64(0)
	}
	return &invoice
}

func MapEntityToInvoiceLine(entity *neo4jentity.InvoiceLineEntity) *model.InvoiceLine {
	if entity == nil {
		return nil
	}
	return &model.InvoiceLine{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
		},
		Description: entity.Name,
		Price:       entity.Price,
		Quantity:    entity.Quantity,
		Total:       entity.TotalAmount,
		Subtotal:    entity.Amount,
		TaxDue:      entity.Vat,
	}
}

func MapEntitiesToInvoices(entities *neo4jentity.InvoiceEntities) []*model.Invoice {
	var output []*model.Invoice
	for _, v := range *entities {
		output = append(output, MapEntityToInvoice(&v))
	}
	return output
}

func MapEntitiesToInvoiceLines(entities *neo4jentity.InvoiceLineEntities) []*model.InvoiceLine {
	var output []*model.InvoiceLine
	for _, v := range *entities {
		output = append(output, MapEntityToInvoiceLine(&v))
	}
	return output
}
