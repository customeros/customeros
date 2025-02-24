package invoice

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	neoRepo "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/emersion/go-message/mail"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/postmark"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	webhook "github.com/customeros/customeros/packages/server/customer-os-common-module/webhook_temporal"
)

type invoiceService struct {
	log                   logger.Logger
	neo4j                 *neoRepo.Repositories
	postgresRepositories  *postgresrepository.Repositories
	events                *events.EventsService
	contractService       interfaces.ContractService
	sli                   interfaces.ServiceLineItemService
	tenantSettings        interfaces.TenantSettingsService
	postmarkService       interfaces.PostmarkService
	fileService           interfaces.FileService
	cfg                   *config.ExternalServicesConfig
	internalCfg           *config.InternalServicesConfig
	externalSystemService interfaces.ExternalSystemService
}

func NewInvoiceService(log logger.Logger,
	neo4j *neoRepo.Repositories,
	postgresRepositories *postgresrepository.Repositories,
	cfg *config.ExternalServicesConfig,
	internalCfg *config.InternalServicesConfig,
	events *events.EventsService,
	contractService interfaces.ContractService,
	sli interfaces.ServiceLineItemService,
	tenantSettings interfaces.TenantSettingsService,
	postmarkService interfaces.PostmarkService,
	fileService interfaces.FileService,
	externalSystemService interfaces.ExternalSystemService,
) interfaces.InvoiceService {
	return &invoiceService{
		log:                   log,
		neo4j:                 neo4j,
		postgresRepositories:  postgresRepositories,
		cfg:                   cfg,
		internalCfg:           internalCfg,
		events:                events,
		contractService:       contractService,
		sli:                   sli,
		tenantSettings:        tenantSettings,
		postmarkService:       postmarkService,
		fileService:           fileService,
		externalSystemService: externalSystemService,
	}
}

type InvoiceActionMetadata struct {
	Status        string  `json:"status"`
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	InvoiceNumber string  `json:"number"`
	InvoiceId     string  `json:"id"`
}

func (s *invoiceService) SetContractService(contractService interfaces.ContractService) {
	s.contractService = contractService
}

func (s *invoiceService) SetServiceLineItemService(sli interfaces.ServiceLineItemService) {
	s.sli = sli
}

func (s *invoiceService) IsInitialized() bool {
	return utils.IsInitialized(s)
}

func (s *invoiceService) InvoiceContract(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, contractId string, dataFields data_fields.InvoiceFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.InvoiceContract")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	invoiceId := ""

	contractEntity, err := s.contractService.GetById(ctx, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	// validate contract
	if contractEntity == nil {
		err := fmt.Errorf("contract not found")
		tracing.TraceErr(span, err)
		return "", err
	}

	// validate input data
	if !dataFields.DryRun && dataFields.TenantBillingProfile == nil {
		err := fmt.Errorf("missing tenant billing profile")
		tracing.TraceErr(span, err)
		return "", err
	} else if dataFields.TenantBillingProfile == nil {
		dataFields.FillWithSampleTenantBillingProfile()
	}

	offCycle := false
	dryRun := dataFields.DryRun
	preview := dataFields.Preview

	if dryRun == false {
		if contractEntity.ContractStatus != neo4jenum.ContractStatusLive && contractEntity.ContractStatus != neo4jenum.ContractStatusOutOfContract {
			err := fmt.Errorf("contract not live")
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	if preview == true && dryRun == false {
		err := fmt.Errorf("preview invoices should be dry run")
		tracing.TraceErr(span, err)
		return "", err
	}

	// prepare currency
	currency := contractEntity.Currency.String()
	if currency == "" {
		dbNode, err := s.neo4j.TenantReadRepository.GetTenantSettings(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		tenantSettings := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)
		currency = tenantSettings.BaseCurrency.String()
	}

	// prepare postpaid flag
	dbNode, err := s.neo4j.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(dbNode)
	isPostpaid := tenantSettingsEntity.InvoicingPostpaid

	// preload countries
	contractCountry := contractEntity.Country
	countryDbNode, _ := s.neo4j.CountryReadRepository.GetCountryByCodeIfExists(ctx, contractCountry)
	if countryDbNode != nil {
		countryEntity := neo4jmapper.MapDbNodeToCountryEntity(countryDbNode)
		contractCountry = countryEntity.Name
	}
	tenantBillingProfileCountry := dataFields.TenantBillingProfile.Country
	countryDbNode, _ = s.neo4j.CountryReadRepository.GetCountryByCodeIfExists(ctx, tenantBillingProfileCountry)
	if countryDbNode != nil {
		countryEntity := neo4jmapper.MapDbNodeToCountryEntity(countryDbNode)
		tenantBillingProfileCountry = countryEntity.Name
	}

	// prepare and validate dates
	var invoicePeriodStart, invoicePeriodEnd time.Time
	if preview && dataFields.InvoiceStartDate != nil {
		invoicePeriodStart = *dataFields.InvoiceStartDate
	} else {
		if contractEntity.NextInvoiceDate != nil {
			invoicePeriodStart = *contractEntity.NextInvoiceDate
		} else {
			invoicePeriodStart = *contractEntity.InvoicingStartDate
		}
	}
	if preview && dataFields.InvoiceEndDate != nil {
		invoicePeriodEnd = *dataFields.InvoiceEndDate
	} else {
		invoicePeriodEnd = s.prepareInvoiceCycleEndDate(ctx, invoicePeriodStart, tenant, contractEntity)
	}

	contractReadyForInvoicingByDates := true
	if !dryRun {
		if isPostpaid {
			contractReadyForInvoicingByDates = utils.EndOfDayInUTC(invoicePeriodEnd).Before(utils.Now())
		} else {
			contractReadyForInvoicingByDates = invoicePeriodEnd.After(invoicePeriodStart)
		}
	}
	if !contractReadyForInvoicingByDates {
		err := fmt.Errorf("contract not ready for invoicing by dates")
		tracing.TraceErr(span, err)
		return "", err
	}

	// generate invoice id
	invoiceId, err = s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelInvoice)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tracing.TagEntity(span, invoiceId)

	// prepare invoice number
	invoiceNumber := dataFields.InvoiceNumber
	if invoiceNumber == "" && !offCycle {
		filledInvoiceDbNode, err := s.neo4j.InvoiceReadRepository.GetFirstPreviewFilledInvoice(ctx, tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "GetFirstPreviewFilledInvoice"))
		}
		if filledInvoiceDbNode != nil {
			filledInvoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(filledInvoiceDbNode)
			invoiceNumber = filledInvoiceEntity.Number
		}
	}
	if invoiceNumber == "" || (dryRun && !preview) {
		invoiceNumber = s.prepareInvoiceNumber(ctx, tenant)
	}

	// prepare issue and due date
	issuedDate := utils.ToDate(utils.Now())
	if dryRun {
		issuedDate = utils.ToDate(invoicePeriodStart)
		if isPostpaid {
			issuedDate = utils.ToDate(invoicePeriodEnd).AddDate(0, 0, 1)
		}
	}
	dueDate := issuedDate.AddDate(0, 0, int(contractEntity.DueDays))

	// prepare other fielcs
	invoiceNote := contractEntity.InvoiceNote
	if preview {
		invoiceNote = ""
	}
	source := neo4jentity.DataSourceOpenline.String()
	appSource := common.GetAppSourceFromContext(ctx)

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {

		// Step 1 - create invoice node for the contract
		data := neo4jrepository.InvoiceCreateFields{
			ContractId:           contractId,
			Currency:             neo4jenum.DecodeCurrency(currency),
			DryRun:               dryRun,
			OffCycle:             offCycle,
			Postpaid:             isPostpaid,
			Preview:              preview,
			BillingCycleInMonths: contractEntity.BillingCycleInMonths,
			PeriodStartDate:      invoicePeriodStart,
			PeriodEndDate:        invoicePeriodEnd,
			CreatedAt:            utils.Now(),
			IssuedDate:           issuedDate,
			DueDate:              dueDate,
			Status:               neo4jenum.InvoiceStatusInitialized,
			Source:               source,
			AppSource:            appSource,
			Note:                 invoiceNote,
		}
		err = s.neo4j.InvoiceWriteRepository.CreateInvoiceForContract(ctx, txWithPostCommit.Tx, tenant, invoiceId, data)
		if err != nil {
			s.log.Errorf("Error while creating invoice for contract %s: %s", contractId, err.Error())
			return nil, err
		}

		// Step 2 - Remove previous initialized invoices, if any
		if dryRun && preview {
			err = s.neo4j.InvoiceWriteRepository.DeletePreviewCycleInitializedInvoices(ctx, txWithPostCommit.Tx, tenant, contractId, invoiceId)
			if err != nil {
				s.log.Errorf("Error while deleting preview invoice for contract %s: %s", contractId, err.Error())
				return nil, err
			}
		}

		// Step 3 - get temp invoice entity
		invoiceEntity, err := s.GetById(ctx, txWithPostCommit.Tx, invoiceId)
		if err != nil {
			return nil, err
		}

		// Step 4 - get service line items for the contract
		sliEntities, err := s.sli.GetServiceLineItemsForContract(ctx, contractId)
		if err != nil {
			return nil, err
		}

		// Step 5 - invoke fill invoice (calculates all data), Done only for cycle invoices
		var invoiceLines []*neo4jentity.InvoiceLineEntity
		if !offCycle {
			invoiceEntity, invoiceLines, err = s.fillCycleInvoice(ctx, invoiceEntity, *sliEntities)
			if err != nil {
				return nil, err
			}
		}

		// Step 6 - prepare invoice status
		invoiceStatus := neo4jenum.InvoiceStatusDue
		if len(invoiceLines) == 0 {
			invoiceStatus = neo4jenum.InvoiceStatusEmpty
		} else {
			if dryRun && preview {
				if contractEntity.ContractStatus == neo4jenum.ContractStatusOutOfContract {
					invoiceStatus = neo4jenum.InvoiceStatusOnHold
				} else {
					invoiceStatus = neo4jenum.InvoiceStatusScheduled
				}
			} else if invoiceEntity.TotalAmount == 0 {
				invoiceStatus = neo4jenum.InvoiceStatusPaid
			}
		}

		// Step 7 - save filled invoice and invoice lines
		fillFields := neo4jrepository.InvoiceFillFields{
			InvoiceNumber:                invoiceNumber,
			CustomerName:                 contractEntity.OrganizationLegalName,
			CustomerEmail:                contractEntity.InvoiceEmail,
			CustomerAddressLine1:         contractEntity.AddressLine1,
			CustomerAddressLine2:         contractEntity.AddressLine2,
			CustomerAddressZip:           contractEntity.Zip,
			CustomerAddressLocality:      contractEntity.Locality,
			CustomerAddressCountry:       contractCountry,
			CustomerAddressRegion:        contractEntity.Region,
			ProviderLogoRepositoryFileId: dataFields.TenantBillingProfile.LogoRepositoryFileId,
			ProviderName:                 dataFields.TenantBillingProfile.LegalName,
			ProviderAddressLine1:         dataFields.TenantBillingProfile.AddressLine1,
			ProviderAddressLine2:         dataFields.TenantBillingProfile.AddressLine2,
			ProviderAddressZip:           dataFields.TenantBillingProfile.Zip,
			ProviderAddressLocality:      dataFields.TenantBillingProfile.Locality,
			ProviderAddressCountry:       tenantBillingProfileCountry,
			ProviderAddressRegion:        dataFields.TenantBillingProfile.Region,
			Amount:                       invoiceEntity.Amount,
			VAT:                          invoiceEntity.Vat,
			TotalAmount:                  invoiceEntity.TotalAmount,
			Status:                       invoiceStatus,
			ProviderEmail:                dataFields.FromEmail,
		}
		err = s.neo4j.InvoiceWriteRepository.FillInvoice(ctx, txWithPostCommit.Tx, tenant, invoiceId, fillFields)
		if err != nil {
			s.log.Errorf("Error while filling invocie with details %s: %s", invoiceId, err.Error())
			return nil, err
		}
		for _, item := range invoiceLines {
			invoiceLineData := neo4jrepository.InvoiceLineCreateFields{
				CreatedAt:               utils.Now(),
				SkuId:                   item.SkuId,
				SkuName:                 item.SkuName,
				Name:                    utils.FirstNotEmptyString(item.SkuName, item.Name),
				Description:             item.Description,
				Price:                   item.Price,
				Quantity:                item.Quantity,
				Amount:                  item.Amount,
				VAT:                     item.Vat,
				TotalAmount:             item.TotalAmount,
				BilledType:              item.BilledType,
				Source:                  source,
				AppSource:               appSource,
				ServiceLineItemId:       item.ServiceLineItemId,
				ServiceLineItemParentId: item.ServiceLineItemParentId,
			}
			err = s.neo4j.InvoiceLineWriteRepository.CreateInvoiceLine(ctx, txWithPostCommit.Tx, tenant, invoiceId, item.Id, invoiceLineData)
			if err != nil {
				s.log.Errorf("Error while inserting invoice line %s for invoice %s: %s", item.Id, invoiceId, err.Error())
				return "", err
			}
		}

		// load invoice and invoice lines after fill
		invoiceEntityAfterFill, err := s.GetById(ctx, txWithPostCommit.Tx, invoiceId)
		if err != nil {
			return nil, err
		}

		var invoiceLineEntities []*neo4jentity.InvoiceLineEntity
		invoiceLinesNodes, err := s.neo4j.InvoiceLineReadRepository.GetAllForInvoice(ctx, txWithPostCommit.Tx, tenant, invoiceId)
		if err != nil {
			return nil, err
		}
		for _, invoiceLineNode := range invoiceLinesNodes {
			invoiceLineEntities = append(invoiceLineEntities, neo4jmapper.MapDbNodeToInvoiceLineEntity(invoiceLineNode))
		}

		// Step 8 - generate invoice PDF
		if invoiceEntityAfterFill.Status != neo4jenum.InvoiceStatusEmpty {
			fileId, err := s.generateInvoicePDF(ctx, invoiceEntityAfterFill, contractEntity, invoiceLineEntities, *dataFields.TenantBillingProfile)
			if err != nil {
				return nil, err
			}
			// save pdf file id
			err = s.neo4j.InvoiceWriteRepository.InvoicePdfGenerated(ctx, txWithPostCommit.Tx, tenant, invoiceId, fileId)
			if err != nil {
				return nil, err
			}
		}

		// Step 9 - create invoice action
		if !dryRun && invoiceStatus != neo4jenum.InvoiceStatusEmpty {
			s.createInvoiceAction(ctx, txWithPostCommit.Tx, tenant, "", *invoiceEntityAfterFill)
		}

		// Step 10 - move contract to next invoice date
		if !dryRun {
			nextInvoiceDate := invoicePeriodEnd.AddDate(0, 0, 1)
			contractDataFields := data_fields.ContractSaveFields{
				NextInvoiceDate: utils.ToPtr(nextInvoiceDate),
			}
			_, err = s.contractService.Save(ctx, txWithPostCommit, &contractEntity.Id, contractDataFields)
			if err != nil {
				return nil, err
			}
		}

		// Step 11 - invoice webhooks
		if invoiceEntityAfterFill.InvoiceInternalFields.InvoiceFinalizedWebhookProcessedAt == nil {
			// dispatch invoice finalized event
			err = s.dispatchInvoiceFinalizedEvent(ctx, tenant, invoiceEntityAfterFill, contractEntity, invoiceLineEntities)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "dispatchInvoiceFinalizedEvent"))
				s.log.Errorf("Error while dispatching invoice finalized event for invoice %s: %s", invoiceEntityAfterFill.Id, err.Error())
				// TODO: must implement retry mechanism for dispatching invoice finalized event
			}
			err = s.neo4j.CommonWriteRepository.UpdateTimePropertyInTx(ctx, txWithPostCommit.Tx, tenant, model.NodeLabelInvoice, invoiceId, string(neo4jentity.InvoicePropertyFinalizedWebhookProcessedAt), utils.NowPtr())
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error setting invoice finalized webhook processed"))
				s.log.Errorf("Error setting invoice finalized webhook processed for invoice %s: %s", invoiceEntity.Id, err.Error())
			}
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// event for storing invoicing details
			createDto := dto.CreateInvoice{
				Id:                   invoiceId,
				Currency:             currency,
				DryRun:               dryRun,
				Preview:              preview,
				Postpaid:             isPostpaid,
				Number:               invoiceNumber,
				Note:                 invoiceNote,
				BillingCycleInMonths: contractEntity.BillingCycleInMonths,
				InvoicePeriodStart:   invoicePeriodStart,
				InvoicePeriodEnd:     invoicePeriodEnd,
				OffCycle:             offCycle,
				Source:               source,
				AppSource:            appSource,
				DueDate:              dueDate,
				IssuedDate:           issuedDate,
			}
			err = s.events.Publisher.PublishFanoutEvent(ctx, invoiceId, model.INVOICE, createDto)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateInvoice"))
			}

			// finalized invoice event
			err = s.events.Publisher.PublishFanoutEvent(ctx, invoiceId, model.INVOICE, dto.InvoiceFinalized{})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message InvoiceFinalized"))
			}

			// invoice completion event for FE
			s.events.Publisher.PublishNotification(ctx, tenant, invoiceId, model.INVOICE, utils.NewEventCompletedDetails().WithCreate())

			// internal slack notification
			err = s.slackInvoiceFinalizedWebhook(ctx, tenant, invoiceEntityAfterFill)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to send slack webhook"))
			}

			if !invoiceEntityAfterFill.OffCycle && !invoiceEntityAfterFill.DryRun {
				err = s.neo4j.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, tenant, contractId, "")
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error while deleting preview invoice for contract %s: %s", contractId, err.Error())
				}

				start := utils.ToDate(invoicePeriodEnd).AddDate(0, 0, 1)
				end := s.prepareInvoiceCycleEndDate(ctx, start, tenant, contractEntity)

				previewInvoiceFields := data_fields.InvoiceFields{
					Preview:          true,
					DryRun:           true,
					InvoiceStartDate: &start,
					InvoiceEndDate:   &end,
					TenantBillingProfile: &data_fields.TenantBillingProfile{
						LogoRepositoryFileId:       dataFields.TenantBillingProfile.LogoRepositoryFileId,
						Country:                    dataFields.TenantBillingProfile.Country,
						LegalName:                  dataFields.TenantBillingProfile.LegalName,
						AddressLine1:               dataFields.TenantBillingProfile.AddressLine1,
						AddressLine2:               dataFields.TenantBillingProfile.AddressLine2,
						Zip:                        dataFields.TenantBillingProfile.Zip,
						Locality:                   dataFields.TenantBillingProfile.Locality,
						Region:                     dataFields.TenantBillingProfile.Region,
						IncludeBankTransferDetails: dataFields.TenantBillingProfile.IncludeBankTransferDetails,
						BankName:                   dataFields.TenantBillingProfile.BankName,
						AccountNumber:              dataFields.TenantBillingProfile.AccountNumber,
						IBAN:                       dataFields.TenantBillingProfile.IBAN,
						BIC:                        dataFields.TenantBillingProfile.BIC,
						SortCode:                   dataFields.TenantBillingProfile.SortCode,
						RoutingNumber:              dataFields.TenantBillingProfile.RoutingNumber,
						OtherDetails:               dataFields.TenantBillingProfile.OtherDetails,
					},
				}
				_, err = s.InvoiceContract(ctx, nil, contractId, previewInvoiceFields)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error while creating preview invoice for contract %s: %s", contractId, err.Error())
				}
			} else if preview && dryRun {
				err = s.neo4j.InvoiceWriteRepository.DeletePreviewCycleInvoices(ctx, tenant, contractId, invoiceId)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error while deleting preview invoice for contract %s: %s", contractId, err.Error())
				}
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invoicing failed"))
		return "", err
	}

	return invoiceId, nil
}

func (s *invoiceService) slackInvoiceFinalizedWebhook(ctx context.Context, tenant string, invoice *neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.slackInvoiceFinalizedWebhook")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	if invoice.DryRun {
		return nil
	}

	if s.cfg.SlackConfig.InternalAlertsRegisteredWebhook == "" {
		return nil
	}

	// get organization linked to invoice
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// Create a struct to hold the JSON data
	type SlackMessage struct {
		Text string `json:"text"`
	}
	message := SlackMessage{Text: fmt.Sprintf("Tenant %s, Invoice %s has been finalized for customer %s", tenant, invoice.Number, organizationEntity.Name)}
	// Convert struct to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	// Send POST request
	resp, err := http.Post(s.cfg.SlackConfig.InternalAlertsRegisteredWebhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return err
	}
	defer resp.Body.Close()

	span.LogFields(log.String("result.status", resp.Status))

	return nil
}

func (s *invoiceService) prepareInvoiceCycleEndDate(ctx context.Context, start time.Time, tenant string, contractEntity *neo4jentity.ContractEntity) time.Time {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.prepareInvoiceCycleEndDate")
	defer span.Finish()

	nextStart := start.AddDate(0, int(contractEntity.BillingCycleInMonths), 0)
	if start.Day() == 1 {
		// if previous invoice was generated end of month, we need to substract extra 1 day
		previousCycleInvoiceDbNode, err := s.neo4j.InvoiceReadRepository.GetPreviousCycleInvoice(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(nil, errors.Wrap(err, "Error getting previous cycle invoice"))
		}
		if previousCycleInvoiceDbNode != nil {
			previousInvoice := neo4jmapper.MapDbNodeToInvoiceEntity(previousCycleInvoiceDbNode)
			if previousInvoice.PeriodStartDate.Day() != 1 {
				nextStart = nextStart.AddDate(0, -1, 0)
				nextStart = time.Date(nextStart.Year(), nextStart.Month(), previousInvoice.PeriodStartDate.Day(), 0, 0, 0, 0, nextStart.Location())
			}
		}
	}
	invoiceCycleEnd := nextStart.AddDate(0, 0, -1)
	span.LogFields(log.Object("result.invoiceCycleEnd", invoiceCycleEnd))
	return invoiceCycleEnd
}

func (s *invoiceService) prepareInvoiceNumber(ctx context.Context, tenant string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.prepareInvoiceNumber")
	defer span.Finish()

	maxAttempts := 20
	var invoiceNumber string
	for attempt := 1; attempt < maxAttempts+1; attempt++ {
		invoiceNumber = s.generateNewRandomInvoiceNumber()
		invoiceNumberEntity := postgresentity.InvoiceNumberEntity{
			InvoiceNumber: invoiceNumber,
			Tenant:        tenant,
			Attempts:      attempt,
		}
		innerErr := s.postgresRepositories.InvoiceRepository.Reserve(ctx, invoiceNumberEntity)
		if innerErr == nil {
			break
		}
	}

	span.LogFields(log.String("invoiceNumber", invoiceNumber))
	return invoiceNumber
}

func (s *invoiceService) generateNewRandomInvoiceNumber() string {
	digits := "0123456789"
	consonants := "BCDFGHJKLMNPQRSTVWXYZ"
	invoiceNumber := utils.GenerateRandomStringFromCharset(3, consonants) + "-" + utils.GenerateRandomStringFromCharset(5, digits)
	return invoiceNumber
}

func (s *invoiceService) GetById(ctx context.Context, tx *neo4j.ManagedTransaction, invoiceId string) (*neo4jentity.InvoiceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	if invoiceDbNode, err := s.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, tx, common.GetTenantFromContext(ctx), invoiceId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with id {%s} not found", invoiceId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), nil
	}
}

func (s *invoiceService) GetByIdAcrossAllTenants(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetByIdAcrossAllTenants")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("invoiceId", invoiceId)

	invoiceDbNode, tenant, err := s.neo4j.InvoiceReadRepository.GetInvoiceByIdAcrossAllTenants(ctx, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting invoice by id"))
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with id {%s} not found", invoiceId))
		return nil, "", wrappedErr
	}
	if invoiceDbNode == nil {
		return nil, "", nil
	} else {
		return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), tenant, nil
	}
}

func (s *invoiceService) GetByNumber(ctx context.Context, number string) (*neo4jentity.InvoiceEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetByNumber")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("number", number))

	if invoiceDbNode, err := s.neo4j.InvoiceReadRepository.GetInvoiceByNumber(ctx, common.GetTenantFromContext(ctx), number); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Invoice with number {%s} not found", number))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode), nil
	}
}

func (s *invoiceService) GetInvoiceLinesForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.InvoiceLineEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoiceLinesForInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("invoiceIds", invoiceIds))

	invoiceLines, err := s.neo4j.InvoiceLineReadRepository.GetAllForInvoices(ctx, common.GetTenantFromContext(ctx), invoiceIds)
	if err != nil {
		return nil, err
	}
	invoiceLineEntities := make(neo4jentity.InvoiceLineEntities, 0, len(invoiceLines))
	for _, v := range invoiceLines {
		invoiceLineEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(v.Node)
		invoiceLineEntity.DataloaderKey = v.LinkedNodeId
		invoiceLineEntities = append(invoiceLineEntities, *invoiceLineEntity)
	}
	return &invoiceLineEntities, nil
}

func (s *invoiceService) GetInvoicesForContracts(ctx context.Context, contractIds []string) (*neo4jentity.InvoiceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoicesForContracts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractIds", contractIds))

	invoices, err := s.neo4j.InvoiceReadRepository.GetAllForContracts(ctx, common.GetTenantFromContext(ctx), contractIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	invoiceEntities := make(neo4jentity.InvoiceEntities, 0, len(invoices))
	for _, v := range invoices {
		invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(v.Node)
		invoiceEntity.DataloaderKey = v.LinkedNodeId
		invoiceEntities = append(invoiceEntities, *invoiceEntity)
	}
	return &invoiceEntities, nil
}

func (s *invoiceService) GetInvoicesForServiceLineItems(ctx context.Context, sliIds []string) (*neo4jentity.InvoiceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoicesForServiceLineItems")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("sliIds", sliIds))

	invoices, err := s.neo4j.InvoiceReadRepository.GetAllForServiceLineItems(ctx, common.GetTenantFromContext(ctx), sliIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	invoiceEntities := make(neo4jentity.InvoiceEntities, 0, len(invoices))
	for _, v := range invoices {
		invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(v.Node)
		invoiceEntity.DataloaderKey = v.LinkedNodeId
		invoiceEntities = append(invoiceEntities, *invoiceEntity)
	}
	return &invoiceEntities, nil
}

type SimulateInvoices struct {
	ContractId         string
	IssueDate          time.Time
	DueDate            time.Time
	InvoicePeriodStart time.Time
	InvoicePeriodEnd   time.Time
	InvoiceNumber      string
	InvoiceLines       []*SimulateInvoiceLine
}

type SimulateInvoiceLine struct {
	Key               string
	ServiceLineItemID string
	ParentID          string
	Name              string
	Price             float64
	Quantity          int64
	Amount            float64
	TotalAmount       float64
}

func (s *invoiceService) SimulateInvoice(ctx context.Context, simulateInvoicesWithChanges *interfaces.SimulateInvoiceRequestData) ([]*interfaces.SimulateInvoiceResponseData, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.SimulateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("simulateInvoicesWithChanges", simulateInvoicesWithChanges))

	if len(simulateInvoicesWithChanges.ServiceLines) == 0 {
		err := fmt.Errorf("no invoice lines to simulate")
		tracing.TraceErr(span, err)
		return nil, err
	}

	var response []*interfaces.SimulateInvoiceResponseData

	tenantSettings, err := s.tenantSettings.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// fetch existing contract and set the next invoice date
	contract, err := s.contractService.GetById(ctx, simulateInvoicesWithChanges.ContractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoiceDate := utils.Today()
	if contract.NextInvoiceDate != nil {
		invoiceDate = *contract.NextInvoiceDate
	} else if contract.InvoicingStartDate != nil {
		invoiceDate = *contract.InvoicingStartDate
	}

	existingSlis, err := s.sli.GetServiceLineItemsForContract(ctx, contract.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	sliEntities := *existingSlis

	// TODO temporary disable prorated off-cycle simulated invoices
	allowProratedOffCycleInvoices := false
	if !tenantSettings.InvoicingPostpaid && allowProratedOffCycleInvoices {
		//determine the interval to compute invoices
		//[ invoicing starts, max service start date ]
		invoicePeriodStartGeneration := invoiceDate
		if contract.InvoicingStartDate != nil {
			invoicePeriodStartGeneration = *contract.InvoicingStartDate
		}
		invoicePeriodEndGeneration := time.Time{}

		for _, sliData := range simulateInvoicesWithChanges.ServiceLines {
			if sliData.ServiceStarted.After(invoicePeriodEndGeneration) {
				invoicePeriodEndGeneration = sliData.ServiceStarted
			}
		}
		invoicePeriodEndGeneration = calculateInvoiceCycleEnd(invoicePeriodEndGeneration, contract.BillingCycleInMonths)

		for true {
			if invoicePeriodStartGeneration.After(invoicePeriodEndGeneration) {
				break
			}

			// identify first SLI change in the simualtion and replace it in the sliEntities
			for _, sliData := range simulateInvoicesWithChanges.ServiceLines {

				prorationNeeded := false

				if sliData.ServiceStarted.Before(invoicePeriodStartGeneration.Add(1)) {
					continue
				}

				if sliData.ServiceStarted.After(calculateInvoiceCycleEnd(invoicePeriodStartGeneration, contract.BillingCycleInMonths)) {
					continue
				}

				if sliData.ServiceLineItemID == "" {
					// new sli item - adding it to the sli entities and trigger proration
					sliEntity := neo4jentity.ServiceLineItemEntity{
						ID:          sliData.Key,
						SkuId:       sliData.SkuID,
						Comments:    sliData.Comments,
						Billed:      sliData.BillingCycle,
						Price:       sliData.Price,
						Quantity:    sliData.Quantity,
						StartedAt:   sliData.ServiceStarted,
						Description: sliData.Description,
						EndedAt:     nil,
					}
					sliEntities = append(sliEntities, sliEntity)
					prorationNeeded = true
				} else {
					// existing sli item - to check if there is any change in the sli item to decide if proration is needed
					existingSli, err := s.sli.GetById(ctx, sliData.ServiceLineItemID)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					if sliData.ServiceStarted.After(utils.Today()) {
						// new version
					} else {
						// update existing version

						if (existingSli.Billed != sliData.BillingCycle || existingSli.Price != sliData.Price || existingSli.Quantity != sliData.Quantity || existingSli.StartedAt != sliData.ServiceStarted || existingSli.VatRate != utils.IfNotNilFloat64(sliData.TaxRate)) && (existingSli.Price*float64(existingSli.Quantity) < sliData.Price*float64(sliData.Quantity)) {
							// existing sli item - new version - proration needed
							for i, sliEntity := range sliEntities {
								if sliEntity.ID == sliData.ServiceLineItemID {
									sliEntities[i].Billed = sliData.BillingCycle
									sliEntities[i].Price = sliData.Price
									sliEntities[i].Quantity = sliData.Quantity
									sliEntities[i].StartedAt = sliData.ServiceStarted
									sliEntities[i].VatRate = utils.IfNotNilFloat64(sliData.TaxRate)
									break
								}
							}
							prorationNeeded = true
						}
					}
				}

				if prorationNeeded {

					proratedInvoice, err := s.SimulateOffCycleInvoice(ctx, contract, &sliEntities, span)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					response = append(response, proratedInvoice)

					contract.NextInvoiceDate = utils.Ptr(proratedInvoice.Invoice.PeriodEndDate.AddDate(0, 0, 1))
					onCycleInvoice, err := s.SimulateOnCycleInvoice(ctx, contract, &sliEntities, span)
					if err != nil {
						tracing.TraceErr(span, err)
						return nil, err
					}

					response = append(response, onCycleInvoice)
				}

			}

			nextInvoiceDate := invoicePeriodStartGeneration.AddDate(0, int(contract.BillingCycleInMonths), 0)
			invoicePeriodStartGeneration = nextInvoiceDate

			nextOnCycleDate := invoiceDate.AddDate(0, int(contract.BillingCycleInMonths), 0)
			contract.NextInvoiceDate = &nextOnCycleDate
		}
	}

	// no proration needed - only on cycle invoice
	if len(response) == 0 {

		contract.NextInvoiceDate = &invoiceDate

		// build sli entities to reflect the changes for the period
		onCycleSliEntities := neo4jentity.ServiceLineItemEntities{}

		// prepare simulation invoice lines grouped by parent id
		invoiceLinesGroupedByParentId := map[string][]interfaces.SimulateInvoiceRequestServiceLineData{}
		for i := range simulateInvoicesWithChanges.ServiceLines {
			sliData := &simulateInvoicesWithChanges.ServiceLines[i]
			if sliData.ServiceLineItemID == "" {
				sliData.ServiceLineItemID = sliData.Key
			}
			if sliData.ServiceLineItemID == "" {
				sliData.ServiceLineItemID = uuid.New().String()
			}
			if sliData.ParentID == "" {
				sliData.ParentID = sliData.ServiceLineItemID
			}
			invoiceLinesGroupedByParentId[sliData.ParentID] = append(invoiceLinesGroupedByParentId[sliData.ParentID], *sliData)
		}
		// per parent group sort by service started date and set end date for each sli as previous sli start date
		for _, sliDataGroup := range invoiceLinesGroupedByParentId {
			sort.Slice(sliDataGroup, func(i, j int) bool {
				return sliDataGroup[i].ServiceStarted.Before(sliDataGroup[j].ServiceStarted)
			})
			for i, sliData := range sliDataGroup {
				if i > 0 {
					sliDataGroup[i-1].ServiceEnded = &sliData.ServiceStarted
				}
			}
		}

		for _, sliData := range simulateInvoicesWithChanges.ServiceLines {
			sliEntity := neo4jentity.ServiceLineItemEntity{
				ID:          utils.IfNotNilString(sliData.ServiceLineItemID),
				ParentID:    utils.IfNotNilString(sliData.ParentID),
				SkuId:       sliData.SkuID,
				Comments:    sliData.Comments,
				Billed:      sliData.BillingCycle,
				Price:       sliData.Price,
				Quantity:    sliData.Quantity,
				StartedAt:   sliData.ServiceStarted,
				EndedAt:     sliData.ServiceEnded,
				VatRate:     utils.IfNotNilFloat64(sliData.TaxRate),
				Canceled:    sliData.Canceled,
				Description: sliData.Description,
			}

			onCycleSliEntities = append(onCycleSliEntities, sliEntity)
		}

		onCycleInvoice, err := s.SimulateOnCycleInvoice(ctx, contract, &onCycleSliEntities, span)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		response = append(response, onCycleInvoice)
	}

	return response, nil
}

func (s *invoiceService) SimulateOnCycleInvoice(ctx context.Context, contract *neo4jentity.ContractEntity, sliEntities *neo4jentity.ServiceLineItemEntities, span opentracing.Span) (*interfaces.SimulateInvoiceResponseData, error) {
	invoiceEntity := &neo4jentity.InvoiceEntity{}
	var invoiceLines []*neo4jentity.InvoiceLineEntity

	tenantSettings, err := s.tenantSettings.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var invoicePeriodStart, invoicePeriodEnd time.Time
	if contract.NextInvoiceDate != nil {
		invoicePeriodStart = *contract.NextInvoiceDate
	} else {
		invoicePeriodStart = *contract.InvoicingStartDate
	}
	invoicePeriodEnd = calculateInvoiceCycleEnd(invoicePeriodStart, contract.BillingCycleInMonths)

	invoiceEntity.Number = s.generateNewRandomInvoiceNumber()
	invoiceEntity.OffCycle = false
	invoiceEntity.Postpaid = tenantSettings.InvoicingPostpaid
	invoiceEntity.BillingCycleInMonths = contract.BillingCycleInMonths
	invoiceEntity.PeriodStartDate = invoicePeriodStart
	invoiceEntity.PeriodEndDate = invoicePeriodEnd
	invoiceEntity.Note = contract.InvoiceNote
	invoiceEntity.IssuedDate = invoicePeriodStart
	invoiceEntity.DueDate = invoicePeriodEnd.AddDate(0, 0, int(contract.DueDays))

	if contract.Currency != "" {
		invoiceEntity.Currency = contract.Currency
	} else {
		invoiceEntity.Currency = tenantSettings.BaseCurrency
	}

	invoiceEntity, invoiceLines, err = s.fillCycleInvoice(ctx, invoiceEntity, *sliEntities)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	onCycleInvoice := &interfaces.SimulateInvoiceResponseData{
		Invoice: invoiceEntity,
		Lines:   []*neo4jentity.InvoiceLineEntity{},
	}
	for _, line := range invoiceLines {
		invoiceLineEntity := &neo4jentity.InvoiceLineEntity{
			Id:                      line.ServiceLineItemId,
			ServiceLineItemParentId: line.ServiceLineItemParentId,
			SkuId:                   line.SkuId,
			SkuName:                 line.SkuName,
			Name:                    utils.FirstNotEmptyString(line.SkuName, line.Name),
			Price:                   line.Price,
			Quantity:                line.Quantity,
			Amount:                  line.Amount,
			TotalAmount:             line.TotalAmount,
			Vat:                     line.Vat,
		}
		onCycleInvoice.Lines = append(onCycleInvoice.Lines, invoiceLineEntity)
	}

	return onCycleInvoice, nil
}

func (s *invoiceService) SimulateOffCycleInvoice(ctx context.Context, contract *neo4jentity.ContractEntity, sliEntities *neo4jentity.ServiceLineItemEntities, span opentracing.Span) (*interfaces.SimulateInvoiceResponseData, error) {
	invoiceEntity := &neo4jentity.InvoiceEntity{}
	var invoiceLines []*neo4jentity.InvoiceLineEntity

	tenantSettings, err := s.tenantSettings.GetTenantSettings(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoicePeriodEnd := contract.InvoicingStartDate
	invoicePeriodStart := utils.Ptr(utils.FirstTimeOfMonth(9999, 12))
	for _, sliData := range *sliEntities {
		if sliData.StartedAt.Before(*invoicePeriodStart) {
			invoicePeriodStart = utils.Ptr(sliData.StartedAt.Add(time.Hour * 24))
		}
	}

	for true {
		invoicePeriodEnd = utils.Ptr(calculateInvoiceCycleEnd(*invoicePeriodEnd, contract.BillingCycleInMonths))
		if invoicePeriodEnd.After(*invoicePeriodStart) {
			break
		} else {
			invoicePeriodEnd = utils.Ptr(invoicePeriodEnd.AddDate(0, 0, 1))
		}
	}

	invoiceEntity.Number = s.generateNewRandomInvoiceNumber()
	invoiceEntity.OffCycle = true
	invoiceEntity.Postpaid = false
	invoiceEntity.BillingCycleInMonths = contract.BillingCycleInMonths
	invoiceEntity.PeriodStartDate = *invoicePeriodStart
	invoiceEntity.PeriodEndDate = *invoicePeriodEnd
	invoiceEntity.IssuedDate = *invoicePeriodStart
	invoiceEntity.DueDate = invoicePeriodEnd.AddDate(0, 0, int(contract.DueDays))

	if contract.Currency != "" {
		invoiceEntity.Currency = contract.Currency
	} else {
		invoiceEntity.Currency = tenantSettings.BaseCurrency
	}

	sliEntitiesForProration := neo4jentity.ServiceLineItemEntities{}
	for _, sliData := range *sliEntities {
		sliEntity := neo4jentity.ServiceLineItemEntity{
			SkuId:       sliData.SkuId,
			Comments:    sliData.Comments,
			Billed:      sliData.Billed,
			Price:       sliData.Price,
			Quantity:    sliData.Quantity,
			StartedAt:   sliData.StartedAt,
			EndedAt:     nil,
			VatRate:     sliData.VatRate,
			Description: sliData.Description,
		}
		sliEntitiesForProration = append(sliEntitiesForProration, sliEntity)
	}

	invoiceEntity, invoiceLines, err = s.FillOffCyclePrepaidInvoice(ctx, invoiceEntity, sliEntitiesForProration)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	onCycleInvoice := &interfaces.SimulateInvoiceResponseData{
		Invoice: invoiceEntity,
		Lines:   []*neo4jentity.InvoiceLineEntity{},
	}
	for _, line := range invoiceLines {
		invoiceLineEntity := &neo4jentity.InvoiceLineEntity{
			Id:                      line.ServiceLineItemId,
			ServiceLineItemParentId: line.ServiceLineItemParentId,
			SkuId:                   line.SkuId,
			SkuName:                 utils.FirstNotEmptyString(line.SkuName, line.Name),
			Name:                    line.Name,
			Price:                   line.Price,
			Quantity:                line.Quantity,
			Amount:                  line.Amount,
			TotalAmount:             line.TotalAmount,
			Vat:                     line.Vat,
		}
		onCycleInvoice.Lines = append(onCycleInvoice.Lines, invoiceLineEntity)
	}

	return onCycleInvoice, nil
}

func (s *invoiceService) NextInvoiceDryRun(ctx context.Context, contractId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.NextInvoiceDryRun")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractId", contractId))

	invoiceFields := data_fields.InvoiceFields{
		DryRun:  true,
		Preview: false,
	}
	invoiceId, err := s.InvoiceContract(ctx, nil, contractId, invoiceFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	return invoiceId, nil
}

func (s *invoiceService) PayInvoice(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.PayInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("invoiceId", invoiceId))

	err := s.UpdateInvoice(ctx, nil, invoiceId, neo4jrepository.InvoiceUpdateFields{
		Status:       neo4jenum.InvoiceStatusPaid,
		UpdateStatus: true,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *invoiceService) VoidInvoice(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.VoidInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)

	err := s.UpdateInvoice(ctx, nil, invoiceId, neo4jrepository.InvoiceUpdateFields{
		Status:       neo4jenum.InvoiceStatusVoid,
		UpdateStatus: true,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *invoiceService) fillCycleInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*neo4jentity.InvoiceLineEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.fillCycleInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagEntityId, invoiceEntity.Id)

	tenant := common.GetTenantFromContext(ctx)

	amount, vat := float64(0), float64(0)
	var invoiceLines []*neo4jentity.InvoiceLineEntity

	referenceTime := invoiceEntity.PeriodStartDate
	periodEndTime := utils.EndOfDayInUTC(invoiceEntity.PeriodEndDate)
	if invoiceEntity.Postpaid {
		referenceTime = periodEndTime
	}

	cancelledSliParentIds := []string{}
	for _, sliEntity := range sliEntities {
		if sliEntity.Canceled && sliEntity.ParentID != "" {
			cancelledSliParentIds = append(cancelledSliParentIds, sliEntity.ParentID)
		}
	}

	reasonForSliExcludedFromInvoicing := map[string]string{}

	for _, sliEntity := range sliEntities {
		// skip paused subscriptions
		if sliEntity.Paused {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is paused"
			continue
		}
		// skip for now usage SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeUsage {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Billed type is Usage"
			continue
		}
		// skip SLI if of None type
		if sliEntity.Billed == neo4jenum.BilledTypeNone {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Billed type is None"
			continue
		}
		// skip SLI if ended on the reference time
		if sliEntity.EndedAt != nil && sliEntity.EndedAt.Before(referenceTime) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI ended before reference time"
			continue
		}
		// skip SLI if not active on the reference time
		if sliEntity.IsRecurrent() && !sliEntity.IsActiveAt(referenceTime) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is not active at reference time"
			continue
		}
		// skip ONE TIME SLI if started after the end period
		if sliEntity.IsOneTime() && sliEntity.StartedAt.After(periodEndTime) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "One time SLI started after the end period"
			continue
		}

		// Any Cancelled SLI should not be invoiced
		if sliEntity.Canceled {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is cancelled"
			continue
		}

		// Any SLI that has any version of same group cancelled should not be invoiced
		if sliEntity.ParentID != "" && utils.Contains(cancelledSliParentIds, sliEntity.ParentID) {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI is part of a cancelled group"
			continue
		}

		// skip SLI if quantity or price is negative
		if sliEntity.Quantity < 0 {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Quantity is negative"
			continue
		}

		// skip SLI if price is negative for non one times
		if sliEntity.Price < 0 && !sliEntity.IsOneTime() {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Price is negative"
			continue
		}

		// skip SLI if quantity is zero
		if sliEntity.Quantity == 0 {
			reasonForSliExcludedFromInvoicing[sliEntity.ID] = "Quantity is zero"
			continue
		}

		calculatedSLIAmount, calculatedSLIVat := float64(0), float64(0)
		invoiceLineCalculationsReady := false
		// process one time SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeOnce {
			// Get all invoices for one time
			invoices, err := s.GetInvoicesForServiceLineItems(ctx, []string{sliEntity.ID})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
			if invoices != nil && len(*invoices) > 0 {
				alreadyInvoiced := false
				for _, invoice := range *invoices {
					if !invoice.DryRun && invoice.Status != neo4jenum.InvoiceStatusInitialized {
						alreadyInvoiced = true
						break
					}
				}
				if alreadyInvoiced {
					// SLI already invoiced
					reasonForSliExcludedFromInvoicing[sliEntity.ID] = "SLI already invoiced"
					continue
				}
			}

			quantity := sliEntity.Quantity
			calculatedSLIAmount = utils.RoundHalfUpFloat64(float64(quantity)*sliEntity.Price, 2)
			calculatedSLIVat = utils.RoundHalfUpFloat64(calculatedSLIAmount*sliEntity.VatRate/100, 2)
			invoiceLineCalculationsReady = true
		}
		// process monthly, quarterly and annually SLIs
		if sliEntity.Billed == neo4jenum.BilledTypeMonthly || sliEntity.Billed == neo4jenum.BilledTypeQuarterly || sliEntity.Billed == neo4jenum.BilledTypeAnnually {
			calculatedSLIAmount = calculateSLIAmountForCycleInvoicing(sliEntity.Quantity, sliEntity.Price, sliEntity.Billed, invoiceEntity.BillingCycleInMonths)
			calculatedSLIAmount = utils.RoundHalfUpFloat64(calculatedSLIAmount, 2)
			calculatedSLIVat = utils.RoundHalfUpFloat64(calculatedSLIAmount*sliEntity.VatRate/100, 2)
			invoiceLineCalculationsReady = true
		}
		if invoiceLineCalculationsReady {
			amount += calculatedSLIAmount
			vat += calculatedSLIVat

			sliName, err := s.sli.GetServiceLineItemName(ctx, sliEntity.ID)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			invoiceLine := neo4jentity.InvoiceLineEntity{
				Name:                    sliName,
				Description:             sliEntity.Description,
				Price:                   utils.RoundHalfUpFloat64(calculatePriceForBilledType(sliEntity.Price, sliEntity.Billed, invoiceEntity.BillingCycleInMonths), 2),
				Quantity:                sliEntity.Quantity,
				Amount:                  calculatedSLIAmount,
				TotalAmount:             utils.RoundHalfUpFloat64(calculatedSLIAmount+calculatedSLIVat, 2),
				Vat:                     calculatedSLIVat,
				ServiceLineItemId:       sliEntity.ID,
				ServiceLineItemParentId: sliEntity.ParentID,
				BilledType:              sliEntity.Billed,
			}
			if sliEntity.SkuId != "" {
				sku, err := s.postgresRepositories.SkuRepository.Get(ctx, tenant, sliEntity.SkuId)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, nil, err
				}
				if sku != nil {
					invoiceLine.SkuId = sku.ID
					invoiceLine.SkuName = sku.Name
				}
			}

			invoiceLines = append(invoiceLines, &invoiceLine)
			continue
		}
		// if remained any unprocessed SLI log an error
		err := errors.Errorf("Unprocessed SLI %s", sliEntity.ID)
		tracing.TraceErr(span, err)
		s.log.Errorf("Error processing SLI during invoicing %s: %s", sliEntity.ID, err.Error())
	}

	span.LogFields(log.Object("result.ignored_SLIs", reasonForSliExcludedFromInvoicing))

	invoiceEntity.Amount = utils.RoundHalfUpFloat64(amount, 2)
	invoiceEntity.Vat = utils.RoundHalfUpFloat64(vat, 2)
	invoiceEntity.TotalAmount = utils.RoundHalfUpFloat64(amount+vat, 2)

	return invoiceEntity, invoiceLines, nil
}

// Deprecated
func (s *invoiceService) FillOffCyclePrepaidInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*neo4jentity.InvoiceLineEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.FillOffCyclePrepaidInvoice")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
	span.SetTag(tracing.SpanTagEntityId, invoiceEntity.Id)

	tenant := common.GetTenantFromContext(ctx)

	// filter out not applicable SLIs
	referenceDate := invoiceEntity.PeriodStartDate
	filteredSliEntities := neo4jentity.ServiceLineItemEntities{}
	for _, sliEntity := range sliEntities {
		// process only monthly, quarterly and annually SLIs
		if sliEntity.Billed != neo4jenum.BilledTypeMonthly &&
			sliEntity.Billed != neo4jenum.BilledTypeQuarterly &&
			sliEntity.Billed != neo4jenum.BilledTypeAnnually &&
			sliEntity.Billed != neo4jenum.BilledTypeOnce {
			continue
		}
		// SLIs that started on or after reference date are not applicable
		if sliEntity.StartedAt.After(referenceDate) || sliEntity.StartedAt.Equal(referenceDate) {
			continue
		}
		// One time invoiced and cancelled SLIs are not applicable
		if sliEntity.Billed == neo4jenum.BilledTypeOnce {
			if sliEntity.Quantity <= 0 || sliEntity.Price == 0 {
				continue
			}
			if sliEntity.Canceled {
				continue
			}
			ilDbNodeAndInvoiceId, err := s.neo4j.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), sliEntity.ParentID)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", sliEntity.ParentID, err.Error())
				return nil, nil, err
			}
			if ilDbNodeAndInvoiceId != nil {
				continue
			}
		}
		filteredSliEntities = append(filteredSliEntities, sliEntity)
	}
	// sort SLIs by startedAt
	sort.Slice(filteredSliEntities, func(i, j int) bool {
		return filteredSliEntities[i].StartedAt.Before(filteredSliEntities[j].StartedAt)
	})
	// group SLIs by parent id
	sliByParentID := map[string][]neo4jentity.ServiceLineItemEntity{}
	for _, sliEntity := range filteredSliEntities {
		sliByParentID[sliEntity.ParentID] = append(sliByParentID[sliEntity.ParentID], sliEntity)
	}

	span.LogFields(log.Int("result - amount of SLIs to process", len(filteredSliEntities)))

	amount, vat := float64(0), float64(0)
	var invoiceLines []*neo4jentity.InvoiceLineEntity

	proratedSliFound := false
	// iterate SLIs by parent id
	for parentId, slis := range sliByParentID {
		// get latest SLI that is active on reference date
		var sliEntityToInvoice *neo4jentity.ServiceLineItemEntity
		for _, sliEntity := range slis {
			if sliEntity.IsActiveAt(invoiceEntity.PeriodStartDate) {
				sliEntityToInvoice = &sliEntity
			}
		}
		// if no SLI is active on reference date, skip
		if sliEntityToInvoice == nil {
			span.LogFields(log.String("result - no active SLI for parent id", parentId))
			continue
		}
		// get invoice line for latest invoiced SLI per parent
		ilDbNodeAndInvoiceId, err := s.neo4j.InvoiceLineReadRepository.GetLatestInvoiceLineWithInvoiceIdByServiceLineItemParentId(ctx, common.GetTenantFromContext(ctx), parentId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error getting latest invoice line for sli parent id {%s}: {%s}", parentId, err.Error())
			return nil, nil, err
		}
		finalSLIAmount, calculatedSLIVat := float64(0), float64(0)
		if sliEntityToInvoice.Billed == neo4jenum.BilledTypeOnce {
			quantity := sliEntityToInvoice.Quantity
			if quantity <= 0 {
				continue
			}
			finalSLIAmount = utils.RoundHalfUpFloat64(float64(quantity)*sliEntityToInvoice.Price, 2)
			if finalSLIAmount == 0 {
				continue
			}
			calculatedSLIVat = utils.RoundHalfUpFloat64(finalSLIAmount*sliEntityToInvoice.VatRate/100, 2)
		} else {
			proratedInvoicedSLIAmount := float64(0)
			if ilDbNodeAndInvoiceId != nil {
				previousInvoiceDbNode, err := s.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, nil, common.GetTenantFromContext(ctx), ilDbNodeAndInvoiceId.LinkedNodeId)
				if err != nil {
					tracing.TraceErr(span, err)
					s.log.Errorf("Error getting invoice {%s}: {%s}", ilDbNodeAndInvoiceId.LinkedNodeId, err.Error())
					return nil, nil, err
				}
				previousInvoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(previousInvoiceDbNode)
				// if previous invoice is for different cycle, charge full amount
				if !previousInvoiceEntity.PeriodEndDate.Before(invoiceEntity.PeriodEndDate) {
					// calculate already invoiced amount, prorated for the period
					invoiceLineEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNodeAndInvoiceId.Node)
					calculatedInvoicedSLIAmountFor1Year := calculateSLIAmountForCycleInvoicing(invoiceLineEntity.Quantity, invoiceLineEntity.Price, invoiceLineEntity.BilledType, 12)
					proratedInvoicedSLIAmount = prorateAnnualSLIAmount(sliEntityToInvoice.StartedAt, invoiceEntity.PeriodEndDate, calculatedInvoicedSLIAmountFor1Year)
					proratedInvoicedSLIAmount = utils.RoundHalfUpFloat64(proratedInvoicedSLIAmount, 2)
				}
			}

			calculatedSLIAmountFor1Year := calculateSLIAmountForCycleInvoicing(sliEntityToInvoice.Quantity, sliEntityToInvoice.Price, sliEntityToInvoice.Billed, 12)
			proratedSLIAmount := prorateAnnualSLIAmount(sliEntityToInvoice.StartedAt, invoiceEntity.PeriodEndDate, calculatedSLIAmountFor1Year)
			proratedSLIAmount = utils.RoundHalfUpFloat64(proratedSLIAmount, 2)
			finalSLIAmount = utils.RoundHalfUpFloat64(proratedSLIAmount-proratedInvoicedSLIAmount, 2)
			span.LogFields(log.Float64(fmt.Sprintf("result - final amount for SLI with parent id %s", parentId), finalSLIAmount))
			if finalSLIAmount <= 0 {
				continue
			}
			calculatedSLIVat = utils.RoundHalfUpFloat64(finalSLIAmount*sliEntityToInvoice.VatRate/100, 2)
			proratedSliFound = true
		}
		amount += finalSLIAmount
		vat += calculatedSLIVat
		invoiceLine := neo4jentity.InvoiceLineEntity{
			Name:                    sliEntityToInvoice.Description,
			Price:                   utils.RoundHalfUpFloat64(calculatePriceForBilledType(sliEntityToInvoice.Price, sliEntityToInvoice.Billed, invoiceEntity.BillingCycleInMonths), 2),
			Quantity:                sliEntityToInvoice.Quantity,
			Amount:                  finalSLIAmount,
			TotalAmount:             utils.RoundHalfUpFloat64(finalSLIAmount+calculatedSLIVat, 2),
			Vat:                     calculatedSLIVat,
			ServiceLineItemId:       sliEntityToInvoice.ID,
			ServiceLineItemParentId: sliEntityToInvoice.ParentID,
		}

		if sliEntityToInvoice.SkuId != "" {
			sku, err := s.postgresRepositories.SkuRepository.Get(ctx, tenant, sliEntityToInvoice.SkuId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
			if sku != nil {
				invoiceLine.SkuId = sku.ID
				invoiceLine.SkuName = sku.Name
				invoiceLine.Name = sku.Name
			}
		}
		invoiceLine.BilledType = sliEntityToInvoice.Billed
		invoiceLines = append(invoiceLines, &invoiceLine)
	}

	if !proratedSliFound && len(invoiceLines) > 0 {
		// if no prorated SLI found, then invoice contains only once billed SLIs
		// accept the invoice if today is monthly anniversary of the contract invoicing start date

		// UPDATE: The rule is on hold, invoice will be issued even if contains only one time SLIs

		//if !isMonthlyAnniversary(invoiceEntity.PeriodEndDate.AddDate(0, 0, 1)) {
		//	invoiceLines = []*invoicepb.InvoiceLine{}
		//}
	}

	invoiceEntity.Amount = utils.RoundHalfUpFloat64(amount, 2)
	invoiceEntity.Vat = utils.RoundHalfUpFloat64(vat, 2)
	invoiceEntity.TotalAmount = utils.RoundHalfUpFloat64(amount+vat, 2)

	return invoiceEntity, invoiceLines, nil
}

func calculateInvoiceCycleEnd(start time.Time, billingCycleInMonths int64) time.Time {
	end := start.AddDate(0, int(billingCycleInMonths), 0)
	previousDay := end.AddDate(0, 0, -1)
	return previousDay
}

func calculateSLIAmountForCycleInvoicing(quantity int64, price float64, billed neo4jenum.BilledType, billingCycleInMonths int64) float64 {
	if quantity == 0 || price == 0 {
		return 0
	}
	unitAmount := calculatePriceForBilledType(price, billed, billingCycleInMonths)
	unitAmount = utils.RoundHalfUpFloat64(unitAmount, 2)
	return float64(quantity) * unitAmount
}

func calculatePriceForBilledType(price float64, billed neo4jenum.BilledType, billingCycleInMonths int64) float64 {
	if billed == neo4jenum.BilledTypeOnce {
		return price
	}

	if billingCycleInMonths == 0 || billed.InMonths() == 0 {
		return 0
	}

	return price * float64(billingCycleInMonths) / float64(billed.InMonths())
}

func prorateAnnualSLIAmount(startDate, endDate time.Time, amount float64) float64 {
	start := utils.ToDate(startDate)
	end := utils.ToDate(endDate)
	days := end.Sub(start).Hours() / 24
	proratedAmount := amount * (days / 365)
	if proratedAmount <= 0 {
		return 0
	}
	return proratedAmount
}

func (s *invoiceService) GetNonDryRunInvoicesForOrganization(ctx context.Context, tenant, organizationId string) (*neo4jentity.InvoiceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetNonDryRunInvoicesForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	invoiceDbNodes, err := s.neo4j.InvoiceReadRepository.GetNonDryRunInvoicesForOrganization(ctx, tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	invoiceEntities := make(neo4jentity.InvoiceEntities, 0, len(invoiceDbNodes))
	for _, v := range invoiceDbNodes {
		invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(v)
		invoiceEntities = append(invoiceEntities, *invoiceEntity)
	}
	return &invoiceEntities, nil
}

func (s *invoiceService) UpdateInvoice(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, invoiceId string, data neo4jrepository.InvoiceUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.UpdateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)
	span.LogFields(log.String("invoiceId", invoiceId))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	invoiceEntityBeforeUpdate, err := s.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
		return err
	}

	err = s.validateStatusTransition(ctx, invoiceEntityBeforeUpdate.Status, data.Status)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	_, err = utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		err = s.neo4j.InvoiceWriteRepository.UpdateInvoice(ctx, txWithPostCommit.Tx, tenant, invoiceId, data)
		if err != nil {
			return nil, err
		}

		invoiceEntityAfterUpdate, err := s.GetById(ctx, txWithPostCommit.Tx, invoiceId)
		if err != nil {
			s.log.Errorf("Error while getting invoice %s: %s", invoiceId, err.Error())
			return nil, err
		}

		s.createInvoiceAction(ctx, txWithPostCommit.Tx, tenant, invoiceEntityBeforeUpdate.Status, *invoiceEntityAfterUpdate)

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// status changed
			if invoiceEntityBeforeUpdate.Status != invoiceEntityAfterUpdate.Status {
				if invoiceEntityAfterUpdate.Status == neo4jenum.InvoiceStatusVoid {
					err = s.events.Publisher.PublishFanoutEvent(ctx, invoiceId, model.INVOICE, dto.InvoiceVoided{
						InvoiceID: invoiceId,
						DryRun:    invoiceEntityAfterUpdate.DryRun,
					})
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "PublishEvent"))
						s.log.Errorf("Error from events processing: %s", err.Error())
					}
				} else if invoiceEntityAfterUpdate.Status == neo4jenum.InvoiceStatusPaid {
					err = s.events.Publisher.PublishFanoutEvent(ctx, invoiceId, model.INVOICE, dto.InvoicePaid{
						InvoiceID: invoiceId,
						DryRun:    invoiceEntityAfterUpdate.DryRun,
					})
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "PublishEvent"))
						s.log.Errorf("Error from events processing: %s", err.Error())
					}
				}
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			// status changed
			if invoiceEntityBeforeUpdate.Status != invoiceEntityAfterUpdate.Status {
				if invoiceEntityAfterUpdate.Status == neo4jenum.InvoiceStatusPaid {
					// do not dispatch invoice paid event if it was already dispatched
					if invoiceEntityAfterUpdate.InvoiceInternalFields.InvoicePaidWebhookProcessedAt == nil {
						// dispatch invoice paid event
						err = s.dispatchInvoicePaidEvent(ctx, tenant, *invoiceEntityAfterUpdate)
						if err != nil {
							tracing.TraceErr(span, errors.Wrap(err, "dispatchInvoicePaidEvent"))
							s.log.Errorf("Error dispatching invoice paid event for invoice %s: %s", invoiceId, err.Error())
						} else {
							err = s.neo4j.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelInvoice, invoiceEntityAfterUpdate.Id, string(neo4jentity.InvoicePropertyPaidWebhookProcessedAt), utils.NowPtr())
							if err != nil {
								tracing.TraceErr(span, errors.Wrap(err, "UpdateTimeProperty"))
								s.log.Errorf("Error setting invoice paid webhook processed for invoice %s: %s", invoiceEntityAfterUpdate.Id, err.Error())
							}
						}
					}
				}
			}

			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			updateDto := dto.UpdateInvoice{}
			if data.UpdateStatus {
				updateDto.Status = utils.StringPtr(data.Status.String())
			}
			if data.UpdatePaymentLink {
				updateDto.PaymentLink = utils.StringPtr(data.PaymentLink)
				updateDto.PaymentLinkValidUntil = data.PaymentLinkValidUntil
			}

			err = s.events.Publisher.PublishFanoutEvent(ctx, invoiceId, model.INVOICE, updateDto)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateInvoice"))
			}

			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *invoiceService) createInvoiceAction(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, previousStatus neo4jenum.InvoiceStatus, invoiceEntity neo4jentity.InvoiceEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.createInvoiceAction")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("invoiceId", invoiceEntity.Id))
	span.LogFields(log.String("previousStatus", previousStatus.String()))
	span.LogFields(log.String("newStatus", invoiceEntity.Status.String()))
	span.LogFields(log.Bool("dryRun", invoiceEntity.DryRun))
	span.LogFields(log.Float64("totalAmount", invoiceEntity.TotalAmount))

	if previousStatus == invoiceEntity.Status {
		return
	}
	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		span.LogFields(log.String("result", "dry run or total amount is 0"))
		return
	}

	metadata, err := utils.ToJson(InvoiceActionMetadata{
		Status:        invoiceEntity.Status.String(),
		Currency:      invoiceEntity.Currency.String(),
		Amount:        invoiceEntity.TotalAmount,
		InvoiceNumber: invoiceEntity.Number,
		InvoiceId:     invoiceEntity.Id,
	})

	actionType := enum.ActionNA
	message := ""
	switch invoiceEntity.Status {
	case neo4jenum.InvoiceStatusDue:
		message = "Invoice N° " + invoiceEntity.Number + " issued with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = enum.ActionInvoiceIssued
	case neo4jenum.InvoiceStatusPaid:
		message = "Invoice N° " + invoiceEntity.Number + " paid in full: " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)
		actionType = enum.ActionInvoicePaid
	case neo4jenum.InvoiceStatusVoid:
		message = "Invoice N° " + invoiceEntity.Number + " voided"
		actionType = enum.ActionInvoiceVoided
	case neo4jenum.InvoiceStatusOverdue:
		message = "Invoice N° " + invoiceEntity.Number + " overdue"
		actionType = enum.ActionInvoiceOverdue
	}
	if actionType == enum.ActionNA {
		span.LogFields(log.String("result", "status not supported"))
		return
	}
	if invoiceEntity.Status == neo4jenum.InvoiceStatusDue {
		_, err = s.neo4j.ActionWriteRepository.MergeByActionType(ctx, tx, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	} else {
		actionId, err := s.neo4j.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelAction)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}
		actionFields := data_fields.ActionFields{
			ActionType: utils.ToPtr(actionType),
			AppSource:  utils.StringPtr(common.GetAppSourceFromContext(ctx)),
			Content:    utils.StringPtr(message),
			Metadata:   utils.StringPtr(metadata),
		}
		err = s.neo4j.ActionWriteRepository.CreateV2(ctx, tx, tenant, actionId, invoiceEntity.Id, model.INVOICE, actionFields)
	}
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}

func (s *invoiceService) SendPaidInvoiceNotification(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.sendPaidInvoiceNotification")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)

	tenant := common.GetTenantFromContext(ctx)
	eventTriggeredByUser := common.GetUserIdFromContext(ctx) != ""

	var invoiceEntity *neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	// load invoice details
	invoiceEntity, err := s.GetById(ctx, nil, invoiceId)
	if err != nil {
		return nil
	}

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		return nil
	}

	// paid notification already sent, skip
	if invoiceEntity.InvoiceInternalFields.PaidInvoiceNotificationSentAt != nil {
		return nil
	}

	// load contract
	contractNode, err := s.neo4j.ContractReadRepository.GetContractForInvoice(ctx, tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePaidV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return nil
	}

	// load tenant billing profile from neo4j
	tenantSettingsDbNode, err := s.neo4j.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)

	cc := contractEntity.InvoiceEmailCC
	cc = append(cc, invoiceEntity.Provider.CC...)
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := contractEntity.InvoiceEmailBCC
	bcc = append(bcc, invoiceEntity.Provider.BCC...)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	postmarkEmail := interfaces.PostmarkEmail{
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentDate}}":    utils.Now().Format("02 Jan 2006"),
		},
		Attachments: []interfaces.PostmarkEmailAttachment{},
	}
	if tenantSettingsEntity.StripeCustomerPortalLink != "" {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = fmt.Sprintf(`PS: If you pay by card you can manage your billing details <a href="%s">here</a>.`, tenantSettingsEntity.StripeCustomerPortalLink)
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = `PS: If you pay by card you can manage your billing details here.`
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = tenantSettingsEntity.StripeCustomerPortalLink
	} else {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = ""
	}

	if eventTriggeredByUser {
		postmarkEmail.WorkflowId = postmark.WorkflowInvoicePaymentReceived
		postmarkEmail.Subject = fmt.Sprintf(postmark.WorkflowInvoicePaymentReceivedSubject, invoiceEntity.Number, invoiceEntity.Provider.Name)
	} else {
		postmarkEmail.WorkflowId = postmark.WorkflowInvoicePaid
		postmarkEmail.Subject = fmt.Sprintf(postmark.WorkflowInvoicePaidSubject, invoiceEntity.Number, invoiceEntity.Provider.Name)
	}

	err = s.appendInvoiceFileToEmailAsAttachment(ctx, tenant, *invoiceEntity, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendInvoiceFileToEmailAsAttachment"))
		s.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.appendProviderLogoToEmail(ctx, tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendProviderLogoToEmail"))
		s.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendCustomerOSLogoToEmail"))
		s.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.postmarkService.SendNotification(ctx, postmarkEmail, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SendNotification"))
		s.log.Errorf("Error sending invoice paid notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = s.neo4j.InvoiceWriteRepository.SetPaidInvoiceNotificationSentAt(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetPaidInvoiceNotificationSentAt"))
		s.log.Errorf("Error setting invoice paid notification sent at for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	return nil
}

func (s *invoiceService) SendVoidedInvoiceNotification(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.sendPaidInvoiceNotification")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)

	tenant := common.GetTenantFromContext(ctx)

	var invoiceEntity *neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	// load invoice details
	invoiceEntity, err := s.GetById(ctx, nil, invoiceId)
	if err != nil {
		return err
	}

	// void notification already sent, skip
	if invoiceEntity.InvoiceInternalFields.VoidInvoiceNotificationSentAt != nil {
		return nil
	}

	// load invoice lines
	invoiceLines, err := s.GetInvoiceLinesForInvoices(ctx, []string{invoiceId})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoiceLinesForInvoices"))
		return err
	}

	allInvoiceLinesEmpty := true
	for _, invoiceLine := range *invoiceLines {
		if invoiceLine.Amount != float64(0) {
			allInvoiceLinesEmpty = false
			break
		}
	}

	if invoiceEntity.DryRun || allInvoiceLinesEmpty {
		return nil
	}

	contractNode, err := s.neo4j.ContractReadRepository.GetContractForInvoice(ctx, tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePaidV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	// load tenant billing profile from neo4j
	cc := contractEntity.InvoiceEmailCC
	cc = append(cc, invoiceEntity.Provider.CC...)
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := contractEntity.InvoiceEmailBCC
	bcc = append(bcc, invoiceEntity.Provider.BCC...)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	postmarkEmail := interfaces.PostmarkEmail{
		WorkflowId:    postmark.WorkflowInvoiceVoided,
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            invoiceEntity.Customer.Email,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(postmark.WorkflowInvoiceVoidedSubject, invoiceEntity.Number), // "Voided invoice " + invoiceEntity.Number,
		TemplateData: map[string]string{
			"{{userFirstName}}":  invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":  invoiceEntity.Number,
			"{{currencySymbol}}": invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":         fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{issueDate}}":      invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		},
		Attachments: []interfaces.PostmarkEmailAttachment{},
	}

	err = s.appendProviderLogoToEmail(ctx, tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendProviderLogoToEmail"))
		s.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "appendCustomerOSLogoToEmail"))
		s.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.postmarkService.SendNotification(ctx, postmarkEmail, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SendNotification"))
		s.log.Errorf("Error sending invoice voided notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = s.neo4j.InvoiceWriteRepository.SetVoidInvoiceNotificationSentAt(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetVoidInvoiceNotificationSentAt"))
		s.log.Errorf("Error setting invoice void notification sent at for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	return nil
}

func (s *invoiceService) appendCustomerOSLogoToEmail(ctx context.Context, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.appendCustomerOSLogoToEmail")
	defer span.Finish()

	currentDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "os.Getwd")
	}

	file, err := utils.GetFileByName(filepath.Join(currentDir, "/static", "customer-os.png"))

	b, err := ioutil.ReadAll(file)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ioutil.ReadAll"))
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, interfaces.PostmarkEmailAttachment{
		Filename:       "customer-os-encoded",
		ContentEncoded: base64.StdEncoding.EncodeToString(b),
		ContentType:    "image/png",
		ContentID:      "cid:customer-os-encoded",
	})

	return nil
}

func (s *invoiceService) appendProviderLogoToEmail(ctx context.Context, tenant, logoFileId string, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.appendProviderLogoToEmail")
	defer span.Finish()

	if logoFileId == "" {
		return nil
	}

	fileInfo, err := s.fileService.GetById(ctx, logoFileId)
	if err != nil {
		return err
	}
	if fileInfo == nil {
		return nil
	}

	fileBytes, err := s.fileService.GetFileBytes(ctx, logoFileId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, interfaces.PostmarkEmailAttachment{
		Filename:       "provider-logo-file-encoded",
		ContentEncoded: base64.StdEncoding.EncodeToString(*fileBytes),
		ContentType:    fileInfo.MimeType,
		ContentID:      "cid:provider-logo-file-encoded",
	})

	return nil
}

func isValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *invoiceService) appendInvoiceFileToEmailAsAttachment(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity, postmarkEmail *interfaces.PostmarkEmail) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.appendInvoiceFileToEmailAsAttachment")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	fileInfo, err := s.fileService.GetById(ctx, invoice.RepositoryFileId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if fileInfo == nil {
		return nil
	}

	invoiceFileBytes, err := s.fileService.GetFileBytes(ctx, invoice.RepositoryFileId)
	if err != nil {
		return err
	}

	postmarkEmail.Attachments = append(postmarkEmail.Attachments, interfaces.PostmarkEmailAttachment{
		Filename:       "Invoice " + invoice.Number + ".pdf",
		ContentEncoded: base64.StdEncoding.EncodeToString(*invoiceFileBytes),
		ContentType:    "application/pdf",
	})

	return nil
}

func (s *invoiceService) dispatchInvoicePaidEvent(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceService.dispatchInvoicePaidEvent")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	tracing.LogObjectAsJson(span, "invoice", invoice)

	// get organization linked to invoice to build payload for webhook
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetOrganizationByInvoiceId"))
		s.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get contract linked to invoice to build payload for webhook
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{organizationEntity.ID})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractsForOrganizations"))
		s.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	contractEntity := neo4jentity.ContractEntity{}
	if len(contractDbNode) > 0 && contractDbNode[0] != nil {
		node := contractDbNode[0].Node
		if node != nil {
			contractEntity = *neo4jmapper.MapDbNodeToContractEntity(node)
		}
	}

	// get invoice line items linked to invoice to build payload for webhook
	invoiceLineDbNodes, err := s.neo4j.InvoiceLineReadRepository.GetAllForInvoice(ctx, nil, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetAllForInvoice"))
		s.log.Errorf("Error getting invoice line items for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	ilEntities := []*neo4jentity.InvoiceLineEntity{}
	for _, ilDbNode := range invoiceLineDbNodes {
		ilEntity := neo4jmapper.MapDbNodeToInvoiceLineEntity(ilDbNode)
		ilEntities = append(ilEntities, ilEntity)
	}

	webhookPayload := webhook.PopulateInvoicePayload(&invoice, &organizationEntity, &contractEntity, ilEntities)
	// dispatch the event
	err = webhook.DispatchWebhook(
		ctx,
		tenant,
		webhook.WebhookEventInvoiceStatusPaid,
		webhookPayload,
		s.postgresRepositories,
		*s.cfg,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "DispatchWebhook"))
		s.log.Errorf("Error dispatching invoice paid event for invoice %s: %s", invoice.Id, err.Error())
		return err
	}

	return nil
}

func (s *invoiceService) dispatchInvoiceFinalizedEvent(ctx context.Context, tenant string,
	invoiceEntity *neo4jentity.InvoiceEntity,
	contractEntity *neo4jentity.ContractEntity,
	invoiceLineEntities []*neo4jentity.InvoiceLineEntity) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "InvoiceService.dispatchInvoiceFinalizedEvent")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, invoiceEntity.Id)

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		return nil
	}

	// get organization linked to invoice to build payload for webhook
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetOrganizationByInvoiceId"))
		s.log.Errorf("Error getting organization for invoice %s: %s", invoiceEntity.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	webhookPayload := webhook.PopulateInvoicePayload(invoiceEntity, &organizationEntity, contractEntity, invoiceLineEntities)
	// dispatch the event
	err = webhook.DispatchWebhook(
		ctx,
		tenant,
		webhook.WebhookEventInvoiceFinalized,
		webhookPayload,
		s.postgresRepositories,
		*s.cfg,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "DispatchWebhook"))
		s.log.Errorf("Error dispatching invoice finalized event for invoice %s: %s", invoiceEntity.Id, err.Error())
		return err
	}

	return nil
}

func (s *invoiceService) generateInvoicePDF(ctx context.Context,
	invoiceEntity *neo4jentity.InvoiceEntity,
	contractEntity *neo4jentity.ContractEntity,
	invoiceLineEntities []*neo4jentity.InvoiceLineEntity,
	tenantBillingProfile data_fields.TenantBillingProfile) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.generateInvoicePDF")
	defer span.Finish()
	tenant := common.GetTenantFromContext(ctx)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, invoiceEntity.Id)

	invoiceHasVat := false

	if invoiceEntity.Vat > 0 {
		invoiceHasVat = true
	}

	dataForPdf := map[string]interface{}{
		"Tenant":                       tenant,
		"CustomerName":                 invoiceEntity.Customer.Name,
		"CustomerEmail":                invoiceEntity.Customer.Email,
		"CustomerAddressLine1":         invoiceEntity.Customer.AddressLine1,
		"CustomerAddressLine2":         invoiceEntity.Customer.AddressLine2,
		"CustomerAddressLine3":         utils.JoinNonEmpty(", ", invoiceEntity.Customer.Locality, invoiceEntity.Customer.Zip),
		"CustomerAddressLine4":         utils.JoinNonEmpty(", ", invoiceEntity.Customer.Region, invoiceEntity.Customer.Country),
		"ProviderLogoExtension":        "",
		"ProviderLogoRepositoryFileId": invoiceEntity.Provider.LogoRepositoryFileId,
		"ProviderName":                 invoiceEntity.Provider.Name,
		"ProviderEmail":                invoiceEntity.Provider.Email,
		"ProviderAddressLine1":         invoiceEntity.Provider.AddressLine1,
		"ProviderAddressLine2":         invoiceEntity.Provider.AddressLine2,
		"ProviderAddressLine3":         utils.JoinNonEmpty(", ", invoiceEntity.Provider.Locality, invoiceEntity.Provider.Zip),
		"ProviderAddressLine4":         utils.JoinNonEmpty(", ", invoiceEntity.Provider.Region, invoiceEntity.Provider.Country),
		"InvoiceNumber":                invoiceEntity.Number,
		"InvoiceIssueDate":             invoiceEntity.CreatedAt.Format("02 Jan 2006"),
		"InvoiceDueDate":               invoiceEntity.DueDate.Format("02 Jan 2006"),
		"InvoiceCurrency":              invoiceEntity.Currency.String() + "" + invoiceEntity.Currency.Symbol(),
		"InvoiceSubtotal":              utils.FormatAmount(invoiceEntity.Amount, 2),
		"InvoiceTotal":                 utils.FormatAmount(invoiceEntity.TotalAmount, 2),
		"InvoiceAmountDue":             utils.FormatAmount(invoiceEntity.TotalAmount, 2),
		"InvoiceLineItems":             []map[string]string{},
		"Note":                         invoiceEntity.Note,
		"CanPayByCheck":                contractEntity.Check,
		"DryRun":                       invoiceEntity.DryRun,
	}

	// Include bank details
	if contractEntity.CanPayWithBankTransfer {
		if tenantBillingProfile.IncludeBankTransferDetails {
			dataForPdf["BankDetailsAvailable"] = true
			dataForPdf["BankAccountName"] = tenantBillingProfile.BankName
			dataForPdf["BankAccountNumber"] = tenantBillingProfile.AccountNumber
			dataForPdf["BankAccountIBAN"] = tenantBillingProfile.IBAN
			dataForPdf["BankAccountBIC"] = tenantBillingProfile.BIC
			dataForPdf["BankAccountSortCode"] = tenantBillingProfile.SortCode
			dataForPdf["BankAccountRoutingNumber"] = tenantBillingProfile.RoutingNumber
			dataForPdf["BankAccountOtherDetails"] = tenantBillingProfile.OtherDetails
		}
	}

	if invoiceHasVat {
		dataForPdf["InvoiceVat"] = fmt.Sprintf("%.2f", invoiceEntity.Vat)
	}

	for _, invoiceLine := range invoiceLineEntities {
		invoiceLineItem := map[string]string{
			"Name":      utils.FirstNotEmptyString(invoiceLine.SkuName, invoiceLine.Name),
			"Quantity":  fmt.Sprintf("%d", invoiceLine.Quantity),
			"UnitPrice": invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceLine.Price, 2),
			"Amount":    invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceLine.Amount, 2),
			"Vat":       invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceLine.Vat, 2),
		}
		if invoiceLine.Description != "" {
			invoiceLineItem["InvoiceLineDescription"] = invoiceLine.Description
		}
		sliDbNode, _ := s.neo4j.ServiceLineItemReadRepository.GetServiceLineItemById(ctx, tenant, invoiceLine.ServiceLineItemId)
		sliEntity := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)

		if invoiceLine.BilledType == neo4jenum.BilledTypeOnce {
			invoiceLineItem["InvoiceLineSubtitle"] = sliEntity.StartedAt.Format("02 Jan 2006")
		}
		// if the invoice line item does not have a subtitle, we will use the period start and end date
		if _, ok := invoiceLineItem["InvoiceLineSubtitle"]; !ok {
			invoiceLineSubtitle := fmt.Sprintf("%s - %s", invoiceEntity.PeriodStartDate.Format("02 Jan 2006"), invoiceEntity.PeriodEndDate.Format("02 Jan 2006"))
			if sliEntity.Billed.IsRecurrent() && sliEntity.Billed.InMonths() != invoiceEntity.BillingCycleInMonths {
				invoiceLineSubtitle += ". "
				invoiceLineSubtitle += invoiceEntity.Currency.Symbol()
				invoiceLineSubtitle += utils.FormatAmount(sliEntity.Price, 2)
				invoiceLineSubtitle += "/"
				switch sliEntity.Billed {
				case neo4jenum.BilledTypeMonthly:
					invoiceLineSubtitle += "month"
				case neo4jenum.BilledTypeQuarterly:
					invoiceLineSubtitle += "quarter"
				case neo4jenum.BilledTypeAnnually:
					invoiceLineSubtitle += "year"
				}
			}
			invoiceLineItem["InvoiceLineSubtitle"] = invoiceLineSubtitle
		}

		if invoiceHasVat {
			invoiceLineItem["InvoiceHasVat"] = "true"
		}

		dataForPdf["InvoiceLineItems"] = append(dataForPdf["InvoiceLineItems"].([]map[string]string), invoiceLineItem)
	}

	// prepare the temp html file
	tmpInvoiceFile, err := os.CreateTemp("", "invoice_*.html")
	if err != nil {
		return "", errors.Wrap(err, "os.TempFile")
	}
	defer os.Remove(tmpInvoiceFile.Name()) // Delete the temporary HTML file when done
	defer tmpInvoiceFile.Close()

	if invoiceEntity.Provider.LogoRepositoryFileId != "" {
		fileMetadata, err := s.fileService.GetById(ctx, invoiceEntity.Provider.LogoRepositoryFileId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "GetFileMetadata"))
			s.log.Errorf("Error getting file metadata for file %s: %s", invoiceEntity.Provider.LogoRepositoryFileId, err.Error())
		} else {
			dataForPdf["ProviderLogoExtension"] = GetFileExtensionFromMetadata(fileMetadata)
		}
	}

	// fill the template with data and store it in temp
	err = FillInvoiceHtmlTemplate(ctx, tmpInvoiceFile, dataForPdf)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "FillInvoiceHtmlTemplate"))
		return "", errors.Wrap(err, "FillInvoiceHtmlTemplate")
	}

	// convert the temp to pdf
	// Max attempts
	maxAttempts := 3
	var pdfBytes *[]byte

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Try to convert the temp to pdf
		pdfBytes, err = ConvertInvoiceHtmlToPdf(ctx, s.fileService, s.internalCfg.PdfConverterConfig.PdfConverterUrl, tmpInvoiceFile, dataForPdf)
		if err == nil {
			// Success, no need to retry
			break
		}
		// Log the error and trace on failure
		tracing.TraceErr(span, errors.Wrap(err, "ConvertInvoiceHtmlToPdf"))
		if attempt == maxAttempts {
			return "", err
		}
	}

	if pdfBytes == nil {
		return "", errors.New("pdfBytes is nil")
	}

	// TODO remove this at some point when we are sure that the pdf is generated correctly
	// Save the PDF file to disk
	os.WriteFile("output.pdf", *pdfBytes, 0644)

	basePath := fmt.Sprintf("/INVOICE/%d/%s", invoiceEntity.CreatedAt.Year(), invoiceEntity.CreatedAt.Format("01"))

	if invoiceEntity.DryRun {
		basePath = basePath + "/DRY_RUN"
	}

	fileDTO, err := s.fileService.UploadSingleFileBytesDirect(ctx, basePath, invoiceEntity.Id, "Invoice - "+invoiceEntity.Number+".pdf", pdfBytes, true)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "UploadSingleFileBytes"))
		return "", err
	}

	if fileDTO.ID == "" {
		return "", errors.New("fileDTO.Id is empty")
	}

	return fileDTO.ID, nil
}

func (s *invoiceService) SendPayInvoiceNotification(ctx context.Context, invoiceId string, allowPayLinkInEmail bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.SendPayInvoiceNotification")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)
	span.LogFields(log.Bool("allowPayLinkInEmail", allowPayLinkInEmail))

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	// load invoice entity
	invoiceNode, err := s.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, nil, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoice"))
		return nil
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		tracing.TraceErr(span, errors.New("invoiceNode is nil"))
		return nil
	}

	// Do not send email if invoice is dry run or total amount is 0 or invoice is not due or overdue
	invoiceStatusAllowedForPayNotification := invoiceEntity.IsDue() || invoiceEntity.IsOverdue()
	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) || !invoiceStatusAllowedForPayNotification {
		span.LogFields(log.String("result", "skipped pay notification"))
		return nil
	}

	if invoiceEntity.Provider.Email == "" {
		s.log.Warnf("Provider email address is empty for invoice %s", invoiceId)
		return nil
	}

	// load contract entity
	contractNode, err := s.neo4j.ContractReadRepository.GetContractForInvoice(ctx, tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return errors.New("contractEntity.InvoiceEmail is empty or invalid")
	}

	// Mark notification requested, to avoid double notifications
	err = s.neo4j.InvoiceWriteRepository.MarkPayNotificationRequested(ctx, tenant, invoiceId, utils.Now())
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error marking pay notification requested for invoice %s: %s", invoiceId, err.Error())
	}

	// prepare email
	workflowId := ""
	if allowPayLinkInEmail && (contractEntity.PayOnline || contractEntity.PayAutomatically) {
		workflowId = postmark.WorkflowInvoiceReadyWithPaymentLink
	} else {
		workflowId = postmark.WorkflowInvoiceReadyNoPaymentLink
	}

	cc := contractEntity.InvoiceEmailCC
	cc = append(cc, invoiceEntity.Provider.CC...)
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := contractEntity.InvoiceEmailBCC
	bcc = append(bcc, invoiceEntity.Provider.BCC...)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	paymentLink := ""
	// prepare payment link for email only if invoice payment link was generated
	if contractEntity.PayOnline || contractEntity.PayAutomatically {
		paymentLink = s.internalCfg.CustomerOsApi.ApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
	}

	postmarkEmail := interfaces.PostmarkEmail{
		WorkflowId:    workflowId,
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(postmark.WorkflowInvoiceReadySubject, invoiceEntity.Number),
		TemplateData: map[string]string{
			"{{organizationName}}": invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":    invoiceEntity.Number,
			"{{currencySymbol}}":   invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":           fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentLink}}":      paymentLink,
		},
		Attachments: []interfaces.PostmarkEmailAttachment{},
	}

	err = s.appendInvoiceFileToEmailAsAttachment(ctx, tenant, invoiceEntity, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendInvoiceFileToEmailAsAttachment")
		tracing.TraceErr(span, wrappedErr)
		s.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = s.appendProviderLogoToEmail(ctx, tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendProviderLogoToEmail")
		tracing.TraceErr(span, wrappedErr)
		s.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = s.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.AppendCustomerOSLogoToEmail"))
		s.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.postmarkService.SendNotification(ctx, postmarkEmail, tenant)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoicePayNotificationV1.SendNotification")
		tracing.TraceErr(span, wrappedErr)
		s.log.Errorf("Error sending invoice pay request notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	s.createPayNotificationInvoiceAction(ctx, tenant, invoiceEntity)

	// Request was successful
	err = s.neo4j.InvoiceWriteRepository.SetPayInvoiceNotificationSentAt(ctx, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "SetPayInvoiceNotificationSentAt"))
		s.log.Errorf("Error setting invoice pay notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (s *invoiceService) createPayNotificationInvoiceAction(ctx context.Context, tenant string, invoiceEntity neo4jentity.InvoiceEntity) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.createPayNotificationInvoiceAction")
	defer span.Finish()
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, invoiceEntity.Id)

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) {
		return
	}

	metadata, err := utils.ToJson(InvoiceActionMetadata{
		Status:        invoiceEntity.Status.String(),
		Currency:      invoiceEntity.Currency.String(),
		Amount:        invoiceEntity.TotalAmount,
		InvoiceNumber: invoiceEntity.Number,
		InvoiceId:     invoiceEntity.Id,
	})

	actionType := enum.ActionInvoiceSent
	message := "Sent invoice N° " + invoiceEntity.Number + " with an amount of " + invoiceEntity.Currency.Symbol() + utils.FormatAmount(invoiceEntity.TotalAmount, 2)

	_, err = s.neo4j.ActionWriteRepository.MergeByActionType(ctx, nil, tenant, invoiceEntity.Id, model.INVOICE, actionType, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "ActionWriteRepository.MergeByActionType"))
		s.log.Errorf("Failed creating invoice action for invoice %s: %s", invoiceEntity.Id, err.Error())
	}
}

func (s *invoiceService) SendPayReminderInvoiceNotification(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.SendPayReminderInvoiceNotification")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	var invoiceEntity neo4jentity.InvoiceEntity
	var contractEntity neo4jentity.ContractEntity

	// load invoice
	invoiceNode, err := s.neo4j.InvoiceReadRepository.GetInvoiceById(ctx, nil, tenant, invoiceId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetInvoice"))
		return nil
	}
	if invoiceNode != nil {
		invoiceEntity = *neo4jmapper.MapDbNodeToInvoiceEntity(invoiceNode)
	} else {
		tracing.TraceErr(span, errors.New("invoiceNode is nil"))
		return nil
	}

	if invoiceEntity.DryRun || invoiceEntity.TotalAmount == float64(0) || !invoiceEntity.IsOverdue() {
		tracing.TraceErr(span, errors.New("remind invoice notification requested for not applicable invoice"))
		return nil
	}

	if invoiceEntity.Provider.Email == "" {
		s.log.Warnf("Provider email address is empty for invoice %s", invoiceId)
		return nil
	}

	// load contract
	contractNode, err := s.neo4j.ContractReadRepository.GetContractForInvoice(ctx, tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetContractForInvoice"))
		return errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.GetContractForInvoice")
	}
	if contractNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractNode)
	} else {
		tracing.TraceErr(span, errors.New("contractNode is nil"))
		return errors.New("contractNode is nil")
	}

	if contractEntity.InvoiceEmail == "" || !isValidEmailSyntax(contractEntity.InvoiceEmail) {
		tracing.TraceErr(span, errors.New("contractEntity.InvoiceEmail is empty or invalid"))
		return errors.New("contractEntity.InvoiceEmail is empty or invalid")
	}

	// load tenant billing profile from neo4j
	tenantSettingsDbNode, err := s.neo4j.TenantReadRepository.GetTenantSettings(ctx, tenant)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "GetTenantSettings"))
		return err
	}
	tenantSettingsEntity := neo4jmapper.MapDbNodeToTenantSettingsEntity(tenantSettingsDbNode)

	// prepare email
	workflowId := ""
	if invoiceEntity.PaymentDetails.PaymentLink == "" {
		workflowId = postmark.WorkflowInvoiceRemindNoPaymentLink
	} else {
		workflowId = postmark.WorkflowInvoiceRemindWithPaymentLink
	}

	cc := contractEntity.InvoiceEmailCC
	cc = append(cc, invoiceEntity.Provider.CC...)
	cc = utils.RemoveEmpties(cc)
	cc = utils.RemoveDuplicates(cc)

	bcc := contractEntity.InvoiceEmailBCC
	bcc = append(bcc, invoiceEntity.Provider.BCC...)
	bcc = utils.RemoveEmpties(bcc)
	bcc = utils.RemoveDuplicates(bcc)

	paymentLink := ""
	// prepare payment link for email only if invoice payment link was generated
	if invoiceEntity.PaymentDetails.PaymentLink != "" {
		paymentLink = s.internalCfg.CustomerOsApi.ApiUrl + "/invoice/" + invoiceEntity.Id + "/pay"
	}

	postmarkEmail := interfaces.PostmarkEmail{
		WorkflowId:    workflowId,
		MessageStream: postmark.PostmarkMessageStreamInvoice,
		From:          invoiceEntity.Provider.Email,
		To:            contractEntity.InvoiceEmail,
		CC:            cc,
		BCC:           bcc,
		Subject:       fmt.Sprintf(postmark.WorkflowInvoiceRemindSubject, invoiceEntity.Number),
		TemplateData: map[string]string{
			"{{organizationName}}": invoiceEntity.Customer.Name,
			"{{invoiceNumber}}":    invoiceEntity.Number,
			"{{currencySymbol}}":   invoiceEntity.Currency.Symbol(),
			"{{amtDue}}":           fmt.Sprintf("%.2f", invoiceEntity.TotalAmount),
			"{{paymentLink}}":      paymentLink,
		},
		Attachments: []interfaces.PostmarkEmailAttachment{},
	}

	if tenantSettingsEntity.StripeCustomerPortalLink != "" {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = fmt.Sprintf(`PS: If you pay by card you can manage your billing details <a href="%s">here</a>.`, tenantSettingsEntity.StripeCustomerPortalLink)
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = `PS: If you pay by card you can manage your billing details here.`
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = tenantSettingsEntity.StripeCustomerPortalLink
	} else {
		postmarkEmail.TemplateData["{{stripeFooterHtml}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterTxt}}"] = ""
		postmarkEmail.TemplateData["{{stripeFooterLink}}"] = ""
	}

	err = s.appendInvoiceFileToEmailAsAttachment(ctx, tenant, invoiceEntity, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendInvoiceFileToEmailAsAttachment")
		tracing.TraceErr(span, wrappedErr)
		s.log.Errorf("Error appending invoice file to email attachment for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = s.appendProviderLogoToEmail(ctx, tenant, invoiceEntity.Provider.LogoRepositoryFileId, &postmarkEmail)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendProviderLogoToEmail")
		tracing.TraceErr(span, wrappedErr)
		s.log.Errorf("Error appending provider logo to email for invoice %s: %s", invoiceId, err.Error())
		return wrappedErr
	}

	err = s.appendCustomerOSLogoToEmail(ctx, &postmarkEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.AppendCustomerOSLogoToEmail"))
		s.log.Errorf("Error appending customeros logo to email for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	err = s.postmarkService.SendNotification(ctx, postmarkEmail, tenant)
	if err != nil {
		wrappedErr := errors.Wrap(err, "InvoiceSubscriber.onInvoiceRemindNotificationV1.SendNotification")
		tracing.TraceErr(span, wrappedErr)
		s.log.Errorf("Error sending invoice remind request notification for invoice %s: %s", invoiceId, err.Error())
		return nil
	}

	// Request was successful
	err = s.neo4j.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelInvoice, invoiceEntity.Id, string(neo4jentity.InvoicePropertyLastRemindInvoiceNotificationSentAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "UpdateTimeProperty"))
		s.log.Errorf("Error setting invoice remind notification sent at for invoice %s: %s", invoiceId, err.Error())
		return err
	}

	return nil
}

func (s *invoiceService) AutopayInvoice(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.AutopayInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if s.cfg.IntegrationAppConfig.IntegrationAppEventWebhookUrls.InvoiceFinalizedUrl == "" {
		err := errors.New("InvoiceFinalizedUrl is not configured")
		tracing.TraceErr(span, err)
		return err
	}

	// load invoice
	invoiceEntity, err := s.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// check if invoice can be auto-paid
	if invoiceEntity.DryRun {
		span.LogFields(log.String("result", "skipped autopay for dry run invoice"))
		return nil
	} else if invoiceEntity.TotalAmount == 0 {
		span.LogFields(log.String("result", "skipped autopay for invoice with total amount of 0"))
		return nil
	} else if !invoiceEntity.IsDue() && !invoiceEntity.IsOverdue() {
		span.LogFields(log.String("result", "skipped autopay for invoice not due or overdue"))
		return nil
	} else if invoiceEntity.InvoiceInternalFields.InvoiceFinalizedSentAt != nil {
		span.LogFields(log.String("result", "skipped autopay for invoice already finalized"))
		return nil
	}

	err = s.integrationAppInvoiceFinalizedWebhook(ctx, tenant, *invoiceEntity)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error invoking invoice finalized webhook"))
		s.log.Errorf("error invoking invoice ready webhook for invoice %s: %s", invoiceEntity.Id, err.Error())
	}

	err = s.neo4j.CommonWriteRepository.UpdateTimeProperty(ctx, tenant, model.NodeLabelInvoice, invoiceEntity.Id, string(neo4jentity.InvoicePropertyInvoiceFinalizedEventSentAt), utils.NowPtr())
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error updating invoice finalized sent at"))
		s.log.Errorf("Error updating invoice finalized sent at for invoice %s: %s", invoiceEntity.Id, err.Error)
		return err
	}

	return nil
}

type IntegrationAppInvoiceFinalizedEventBody struct {
	Tenant                       string `json:"tenant"`
	Currency                     string `json:"currency"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	InvoiceId                    string `json:"invoiceId"`
	InvoiceDescription           string `json:"invoiceDescription"`
	CustomerOsId                 string `json:"customerOsId"`
	Status                       string `json:"status"`
	CustomerEmail                string `json:"customerEmail"`
	CustomerName                 string `json:"customerName"`
	PrimaryStripeCustomerId      string `json:"stripeCustomerId"`
	Pay                          struct {
		PayAutomatically      bool `json:"payAutomatically"`
		CanPayWithCard        bool `json:"canPayWithCard"`
		CanPayWithDirectDebit bool `json:"canPayWithDirectDebit"`
	} `json:"pay"`
}

func (s *invoiceService) integrationAppInvoiceFinalizedWebhook(ctx context.Context, tenant string, invoice neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.integrationAppInvoiceFinalizedWebhook")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	// get organization linked to invoice
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting organization for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	// get contract linked to invoice
	contractDbNode, err := s.neo4j.ContractReadRepository.GetContractForInvoice(ctx, tenant, invoice.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting contract for invoice %s: %s", invoice.Id, err.Error())
		return err
	}
	contractEntity := neo4jentity.ContractEntity{}
	if contractDbNode != nil {
		contractEntity = *neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	}

	// convert amount to the smallest currency unit
	amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.TotalAmount)
	if err != nil {
		return fmt.Errorf("error converting amount to smallest currency unit: %v", err.Error())
	}

	primaryStripeCustomerId, err := s.externalSystemService.GetPrimaryExternalId(ctx, enum.SourceStripe.String(), organizationEntity.ID, model.ORGANIZATION)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting primary stripe customer id for contract %s: %s", contractEntity.Id, err.Error())
	}

	requestBody := IntegrationAppInvoiceFinalizedEventBody{
		Tenant:                       tenant,
		Currency:                     invoice.Currency.String(),
		AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
		InvoiceId:                    invoice.Id,
		InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
		CustomerOsId:                 organizationEntity.CustomerOsId,
		Status:                       invoice.Status.String(),
		CustomerEmail:                contractEntity.InvoiceEmail,
		CustomerName:                 utils.GetReadableNameFromEmail(contractEntity.InvoiceEmail),
		PrimaryStripeCustomerId:      primaryStripeCustomerId,
		Pay: struct {
			PayAutomatically      bool `json:"payAutomatically"`
			CanPayWithCard        bool `json:"canPayWithCard"`
			CanPayWithDirectDebit bool `json:"canPayWithDirectDebit"`
		}{
			PayAutomatically:      contractEntity.PayAutomatically && (invoice.Status == neo4jenum.InvoiceStatusDue || invoice.Status == neo4jenum.InvoiceStatusOverdue),
			CanPayWithCard:        true,
			CanPayWithDirectDebit: true,
		},
	}

	// Convert the request body to JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	// Create an HTTP client
	client := &http.Client{}

	// Create a POST request with headers and body
	req, err := http.NewRequest("POST", s.cfg.IntegrationAppConfig.IntegrationAppEventWebhookUrls.InvoiceFinalizedUrl, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send the POST request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %s", resp.Status)
	}

	return nil
}

func (s *invoiceService) GenerateNewPaymentLink(ctx context.Context, invoiceId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GenerateNewPaymentLink")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, invoiceId)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	tenant := common.GetTenantFromContext(ctx)

	if s.cfg.IntegrationAppConfig.IntegrationAppEventWebhookUrls.GeneratePaymentLinkUrl == "" {
		err := errors.New("GeneratePaymentLinkUrl is not configured")
		tracing.TraceErr(span, err)
		return err
	}

	// load invoice
	invoiceEntity, err := s.GetById(ctx, nil, invoiceId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// check if invoice can be auto-paid
	if invoiceEntity.DryRun {
		span.LogFields(log.String("result", "skipped payment link generation, invoice is dry run"))
		return nil
	} else if invoiceEntity.TotalAmount == 0 {
		span.LogFields(log.String("result", "skipped payment link generation, invoice total amount is 0"))
		return nil
	} else if !invoiceEntity.IsDue() && !invoiceEntity.IsOverdue() {
		span.LogFields(log.String("result", "skipped payment link generation, invoice is not due or overdue"))
		return nil
	}

	// get organization linked to invoice
	organizationDbNode, err := s.neo4j.OrganizationReadRepository.GetOrganizationByInvoiceId(ctx, tenant, invoiceEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error getting organization for invoice %s: %s", invoiceEntity.Id, err.Error())
		return err
	}
	organizationEntity := neo4jentity.OrganizationEntity{}
	if organizationDbNode != nil {
		organizationEntity = *neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
	}

	primaryStripeCustomerId, err := s.externalSystemService.GetPrimaryExternalId(ctx, enum.SourceStripe.String(), organizationEntity.ID, model.ORGANIZATION)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error fetching primary stripe customer ID"))
	}

	err = callIntegrationAppWithApiRequestForNewPaymentLink(ctx, s.cfg.IntegrationAppConfig.WorkspaceKey,
		s.cfg.IntegrationAppConfig.WorkspaceSecret, tenant,
		s.cfg.IntegrationAppConfig.ApiTriggerUrlCreatePaymentLinks, primaryStripeCustomerId, invoiceEntity)

	return nil
}

func (s *invoiceService) validateStatusTransition(ctx context.Context, fromStatus neo4jenum.InvoiceStatus, toStatus neo4jenum.InvoiceStatus) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.validateStatusTransition")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("fromStatus", fromStatus.String()), log.String("toStatus", toStatus.String()))

	if fromStatus == toStatus {
		return nil
	}

	transitionAllowed := true
	if toStatus == neo4jenum.InvoiceStatusPaymentProcessing {
		if fromStatus == neo4jenum.InvoiceStatusPaid {
			transitionAllowed = false
		}
	}

	if !transitionAllowed {
		err := errors.New("Invalid status transition")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

type ApiRequestCreatePaymentLinks struct {
	Input ApiRequestCreatePaymentLinksInput `json:"input"`
}

type ApiRequestCreatePaymentLinksInput struct {
	InvoiceId                    string `json:"invoiceId"`
	AmountInSmallestCurrencyUnit int64  `json:"amountInSmallestCurrencyUnit"`
	Currency                     string `json:"currency"`
	InvoiceDescription           string `json:"invoiceDescription"`
	CustomerEmail                string `json:"customerEmail"`
	PrimaryStripeCustomerId      string `json:"stripeCustomerId"`
}

func callIntegrationAppWithApiRequestForNewPaymentLink(ctx context.Context, key, secret, tenant, url, primaryStripeCustomerId string, invoice *neo4jentity.InvoiceEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "callIntegrationAppWithApiRequestForNewPaymentLink")
	defer span.Finish()
	span.LogKV("url", url)
	span.LogKV("primaryStripeCustomerId", primaryStripeCustomerId)
	span.SetTag(tracing.SpanTagTenant, tenant)

	SigningKey := []byte(secret)

	claims := jwt.MapClaims{
		"id":   tenant,
		"name": tenant,
		// To prevent token from being used for too long
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iss": key,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		return errors.Wrap(err, "Error signing JWT token")
	}

	amountInSmallestCurrencyUnit, err := data.InSmallestCurrencyUnit(invoice.Currency.String(), invoice.TotalAmount)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error converting amount to smallest currency unit"))
		return err
	}

	input := ApiRequestCreatePaymentLinks{
		Input: ApiRequestCreatePaymentLinksInput{
			InvoiceId:                    invoice.Id,
			AmountInSmallestCurrencyUnit: amountInSmallestCurrencyUnit,
			Currency:                     invoice.Currency.String(),
			InvoiceDescription:           fmt.Sprintf("Invoice %s", invoice.Number),
			CustomerEmail:                invoice.Customer.Email,
			PrimaryStripeCustomerId:      primaryStripeCustomerId,
		},
	}
	payload, err := json.Marshal(input)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error marshalling input"))
		return err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error creating HTTP request"))
		return err
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)
	req.Header.Add("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error calling integration app"))
		return err
	}
	return nil
}

func (s *invoiceService) GetInvoicesByIds(ctx context.Context, ids []string) (*neo4jentity.InvoiceEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoicesByIds")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate tenant
	if err := common.ValidateTenant(ctx); err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tenant := common.GetTenantFromContext(ctx)
	invoiceEntities := neo4jentity.InvoiceEntities{}

	// Define chunk size
	const chunkSize = 1000

	// Process IDs in chunks
	for i := 0; i < len(ids); i += chunkSize {
		end := i + chunkSize
		if end > len(ids) {
			end = len(ids)
		}
		chunk := ids[i:end]
		invoiceDbNodes, err := s.neo4j.InvoiceReadRepository.GetInvoicesByIds(ctx, tenant, chunk)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		// Map each DB node to an invoice entity and append to the final result
		for _, invoiceDbNode := range invoiceDbNodes {
			invoiceEntity := neo4jmapper.MapDbNodeToInvoiceEntity(invoiceDbNode)
			invoiceEntities = append(invoiceEntities, *invoiceEntity)
		}
	}

	return &invoiceEntities, nil
}
