package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"reflect"
)

const (
	SortContractName               = "CONTRACT_NAME"
	SearchSortContractBillingCycle = "CONTRACT_BILLING_CYCLE"
	SearchSortContractEnded        = "CONTRACT_ENDED"
	SearchInvoiceDryRunDeprecated  = "DRY_RUN"
	SearchInvoicePreview           = "INVOICE_PREVIEW"
	SearchInvoiceDryRun            = "INVOICE_DRY_RUN"
	SearchSortInvoiceStatus        = "INVOICE_STATUS"
	SearchInvoiceNumberDeprecated  = "NUMBER"
	SearchInvoiceNumber            = "INVOICE_NUMBER"
	SearchInvoiceIssueDate         = "INVOICE_ISSUED_DATE"
)

type InvoiceService interface {
	CountInvoices(ctx context.Context, tenant, organizationId string, where *model.Filter) (int64, error)
	GetInvoices(ctx context.Context, organizationId string, page, limit int, where *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error)
	UpdateInvoice(ctx context.Context, input model.InvoiceUpdateInput) error
}
type invoiceService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewInvoiceService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) InvoiceService {
	return &invoiceService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *invoiceService) CountInvoices(ctx context.Context, tenant, organizationId string, where *model.Filter) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.CountInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.Object("where", where))

	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	invoiceFilterCypher, invoiceFilterParams := "", make(map[string]interface{})

	organizationFilter := new(utils.CypherFilter)
	organizationFilter.Negate = false
	organizationFilter.LogicalOperator = utils.AND
	organizationFilter.Filters = make([]*utils.CypherFilter, 0)

	invoiceFilter := new(utils.CypherFilter)
	invoiceFilter.Negate = false
	invoiceFilter.LogicalOperator = utils.AND
	invoiceFilter.Filters = make([]*utils.CypherFilter, 0)

	if organizationId != "" {
		organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("id", organizationId, utils.EQUALS))
		organizationFilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
	}

	if where != nil {

		for _, f := range where.And {
			if f.Filter.Property == SearchInvoiceDryRunDeprecated {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("dryRun", *f.Filter.Value.Bool))
			}
			if f.Filter.Property == SearchInvoiceDryRun {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("dryRun", *f.Filter.Value.Bool))
			}
			if f.Filter.Property == SearchInvoicePreview {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("preview", *f.Filter.Value.Bool))
			}
		}

	}

	invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterNotEq("status", neo4jenum.InvoiceStatusInitialized.String()))
	invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterNotEq("status", neo4jenum.InvoiceStatusEmpty.String()))

	if len(invoiceFilter.Filters) > 0 {
		invoiceFilterCypher, invoiceFilterParams = invoiceFilter.BuildCypherFilterFragmentWithParamName("i", "i_param_")
	}

	filter := ""
	params := map[string]any{}

	utils.MergeMapToMap(organizationFilterParams, params)
	utils.MergeMapToMap(invoiceFilterParams, params)

	if organizationFilterCypher != "" {
		filter += organizationFilterCypher
	}
	if invoiceFilterCypher != "" {
		if filter != "" {
			filter += " AND "
		}
		filter += invoiceFilterCypher
	}

	if filter != "" {
		filter = " WHERE " + filter
	}

	span.LogFields(log.String("filter", filter))
	span.LogFields(log.Object("params", params))

	return s.repositories.Neo4jRepositories.InvoiceReadRepository.CountInvoices(ctx, tenant, filter, params)
}

func (s *invoiceService) GetInvoices(ctx context.Context, organizationId string, page, limit int, where *model.Filter, sortBy []*model.SortBy) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceService.GetInvoices")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId))
	span.LogFields(log.Object("page", page))
	span.LogFields(log.Object("limit", limit))
	span.LogFields(log.Object("where", where))
	span.LogFields(log.Object("sortBy", sortBy))

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	if len(sortBy) == 0 {
		sortBy = []*model.SortBy{
			{
				By:        "INVOICE_DUE_DATE",
				Direction: model.SortingDirectionDesc,
			},
		}
	}

	cypherSort, err := buildSortMultipleEntities(sortBy, []SortMultipleEntitiesDefinition{
		{
			EntityPrefix:  "CONTRACT",
			EntityMapping: reflect.TypeOf(neo4jentity.ContractEntity{}),
			EntityAlias:   "c",
			EntityDefaults: []SortMultipleEntitiesDefinitionDefault{
				{
					PropertyName: "ENDED_AT",
					AscDefault:   "date('2100-01-01')",
					DescDefault:  "date('1900-01-01')",
				},
			},
		},
		{
			EntityPrefix:  "INVOICE",
			EntityMapping: reflect.TypeOf(neo4jentity.InvoiceEntity{}),
			EntityAlias:   "i",
		},
	})
	if err != nil {
		return nil, err
	}

	organizationFilterCypher, organizationFilterParams := "", make(map[string]interface{})
	contractFilterCypher, contractFilterParams := "", make(map[string]interface{})
	invoiceFilterCypher, invoiceFilterParams := "", make(map[string]interface{})

	organizationFilter := new(utils.CypherFilter)
	organizationFilter.Negate = false
	organizationFilter.LogicalOperator = utils.AND
	organizationFilter.Filters = make([]*utils.CypherFilter, 0)

	contractFilter := new(utils.CypherFilter)
	contractFilter.Negate = false
	contractFilter.LogicalOperator = utils.AND
	contractFilter.Filters = make([]*utils.CypherFilter, 0)

	invoiceFilter := new(utils.CypherFilter)
	invoiceFilter.Negate = false
	invoiceFilter.LogicalOperator = utils.AND
	invoiceFilter.Filters = make([]*utils.CypherFilter, 0)

	if organizationId != "" {
		organizationFilter.Filters = append(organizationFilter.Filters, utils.CreateStringCypherFilter("id", organizationId, utils.EQUALS))
		organizationFilterCypher, organizationFilterParams = organizationFilter.BuildCypherFilterFragmentWithParamName("o", "o_param_")
	}

	if where != nil {

		for _, f := range where.And {
			if f.Filter.Property == SortContractName {
				contractFilter.Filters = append(contractFilter.Filters, utils.CreateStringCypherFilter("name", *f.Filter.Value.Str, utils.CONTAINS))
			}
			if f.Filter.Property == SearchSortContractBillingCycle {
				arrayInt := []int64{}
				for _, v := range *f.Filter.Value.ArrayStr {
					if v == "MONTHLY" {
						arrayInt = append(arrayInt, 1)
					} else if v == "QUARTERLY" {
						arrayInt = append(arrayInt, 3)
					} else if v == "ANNUALLY" {
						arrayInt = append(arrayInt, 12)
					} else if v == "NONE" {
						arrayInt = append(arrayInt, 0)
					}
				}
				contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilterIn("billingCycleInMonths", arrayInt))
			}
			if f.Filter.Property == SearchSortContractEnded {
				if f.Filter.Value.Bool != nil && *f.Filter.Value.Bool {
					contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilterIsNotNull("endedAt"))
				} else {
					contractFilter.Filters = append(contractFilter.Filters, utils.CreateCypherFilterIsNull("endedAt"))
				}
			}
			if f.Filter.Property == SearchSortInvoiceStatus {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterIn("status", *f.Filter.Value.ArrayStr))
			}
			if f.Filter.Property == SearchInvoiceDryRunDeprecated {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("dryRun", *f.Filter.Value.Bool))
			}
			if f.Filter.Property == SearchInvoiceDryRun {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("dryRun", *f.Filter.Value.Bool))
			}
			if f.Filter.Property == SearchInvoicePreview {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("preview", *f.Filter.Value.Bool))
			}
			if f.Filter.Property == SearchInvoiceNumberDeprecated {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("number", *f.Filter.Value.Str))
			}
			if f.Filter.Property == SearchInvoiceNumber {
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterEq("number", *f.Filter.Value.Str))
			}
			if f.Filter.Property == SearchInvoiceIssueDate && f.Filter.Value.ArrayTime != nil && len(*f.Filter.Value.ArrayTime) == 2 {
				times := *f.Filter.Value.ArrayTime
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilter("issuedDate", times[0], utils.GTE))
				invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilter("issuedDate", times[1], utils.LTE))
			}
		}

		if len(contractFilter.Filters) > 0 {
			contractFilterCypher, contractFilterParams = contractFilter.BuildCypherFilterFragmentWithParamName("c", "c_param_")
		}
	}

	invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterNotEq("status", neo4jenum.InvoiceStatusInitialized.String()))
	invoiceFilter.Filters = append(invoiceFilter.Filters, utils.CreateCypherFilterNotEq("status", neo4jenum.InvoiceStatusEmpty.String()))
	if len(invoiceFilter.Filters) > 0 {
		invoiceFilterCypher, invoiceFilterParams = invoiceFilter.BuildCypherFilterFragmentWithParamName("i", "i_param_")
	}

	filter := ""
	params := map[string]any{}

	utils.MergeMapToMap(organizationFilterParams, params)
	utils.MergeMapToMap(contractFilterParams, params)
	utils.MergeMapToMap(invoiceFilterParams, params)

	if organizationFilterCypher != "" {
		filter += organizationFilterCypher
	}
	if contractFilterCypher != "" {
		if filter != "" {
			filter += " AND "
		}
		filter += contractFilterCypher
	}
	if invoiceFilterCypher != "" {
		if filter != "" {
			filter += " AND "
		}
		filter += invoiceFilterCypher
	}

	if filter != "" {
		filter = " WHERE " + filter
	}

	span.LogFields(log.String("filter", filter))
	span.LogFields(log.Object("params", params))

	dbNodesWithTotalCount, err := s.repositories.Neo4jRepositories.InvoiceReadRepository.GetPaginatedInvoices(ctx, common.GetTenantFromContext(ctx),
		paginatedResult.GetSkip(),
		paginatedResult.GetLimit(),
		filter,
		params,
		cypherSort)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	var invoices neo4jentity.InvoiceEntities

	for _, v := range dbNodesWithTotalCount.Nodes {
		invoices = append(invoices, *mapper.MapDbNodeToInvoiceEntity(v))
	}
	paginatedResult.SetRows(&invoices)
	return &paginatedResult, nil
}

func (s *invoiceService) UpdateInvoice(ctx context.Context, input model.InvoiceUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.UpdateInvoice")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	if input.ID == "" {
		err := fmt.Errorf("invoice id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	invoiceExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.ID, model2.NodeLabelInvoice)
	if !invoiceExists {
		err := fmt.Errorf("invoice with id {%s} not found", input.ID)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	data := neo4jrepository.InvoiceUpdateFields{}

	// prepare invoice status
	if input.Status != nil {
		switch *input.Status {
		case model.InvoiceStatusInitialized:
			data.Status = neo4jenum.InvoiceStatusInitialized
		case model.InvoiceStatusEmpty:
			data.Status = neo4jenum.InvoiceStatusEmpty
		case model.InvoiceStatusDraft:
			data.Status = neo4jenum.InvoiceStatusInitialized
		case model.InvoiceStatusDue:
			data.Status = neo4jenum.InvoiceStatusDue
		case model.InvoiceStatusPaid:
			data.Status = neo4jenum.InvoiceStatusPaid
		case model.InvoiceStatusVoid:
			data.Status = neo4jenum.InvoiceStatusVoid
		default:
			data.Status = neo4jenum.InvoiceStatusNone
		}
	}

	if input.Patch {
		if input.Status != nil {
			data.UpdateStatus = true
		} else {
			span.LogFields(log.String("result", "No fields to update"))
			return nil
		}
	}

	err := s.services.CommonServices.InvoiceService.UpdateInvoice(ctx, input.ID, data)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(err.Error())
		return err
	}

	return nil
}
