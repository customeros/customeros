package neo4j_repository

import (
	"fmt"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type InvoiceReadRepository interface {
	GetInvoiceById(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, invoiceId string) (*dbtype.Node, error)
	GetInvoicesByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error)
	GetAllNonDryRunInvoices(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetInvoiceByIdAcrossAllTenants(ctx context.Context, invoiceId string) (*dbtype.Node, string, error)
	GetInvoiceByNumber(ctx context.Context, tenant, invoiceNumber string) (*dbtype.Node, error)
	CountInvoices(ctx context.Context, tenant, filterString string, filterParams map[string]interface{}) (int64, error)
	GetPaginatedInvoices(ctx context.Context, tenant string, skip, limit int, filterCypher string, filterParams map[string]interface{}, sorting *utils.Cypher) (*utils.DbNodesWithTotalCount, error)
	GetInvoicesForPayNotifications(ctx context.Context, minutesFromCreate, minutesFromLastAttempt, lookbackWindow, limit int) ([]*utils.DbNodeAndTenant, error)
	GetInvoicesForPastDueNotifications(ctx context.Context, tenant string, referenceTime time.Time, overdueDays, limit int) ([]*utils.DbNodeAndTenant, error)
	CountNonDryRunInvoicesForContract(ctx context.Context, tenant, contractId string) (int, error)
	GetPreviousCycleInvoice(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetLastIssuedOnCycleInvoiceForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetLastIssuedInvoiceForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetFirstPreviewFilledInvoice(ctx context.Context, tenant, contractId string) (*dbtype.Node, error)
	GetExpiredDryRunInvoices(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error)
	GetPreviewInvoicesForEndedContracts(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error)
	GetAllForContracts(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetAllForServiceLineItems(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetInvoicesForOverdue(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error)
	GetInvoicesForOnHold(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error)
	GetInvoicesForScheduled(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error)
	GetExpiredPaymentProcessingInvoices(ctx context.Context, paymentProcessingMaxDays, limit int) ([]*utils.DbNodeAndTenant, error)
	GetNonDryRunInvoicesForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error)
	GetReadyInvoicesForFinalizedWebhook(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error)
	GetUpcomingInvoices(ctx context.Context, tenant string) ([]*dbtype.Node, error)
}

type invoiceReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInvoiceReadRepository(driver *neo4j.DriverWithContext, database string) InvoiceReadRepository {
	return &invoiceReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *invoiceReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *invoiceReadRepository) CountInvoices(ctx context.Context, tenant, filterString string, filterParams map[string]interface{}) (int64, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.CountInvoices")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_INVOICE]->(i:Invoice_%s) 
			%s
			RETURN count(i)`, tenant, tenant, tenant, filterString)
	params := map[string]any{
		"tenant": tenant,
	}
	utils.MergeMapToMap(filterParams, params)

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return 0, err
	}
	count := dbRecord.(*db.Record).Values[0].(int64)
	spans.LogKV("result.count", count)
	return count, nil
}

func (r *invoiceReadRepository) GetPaginatedInvoices(ctx context.Context, tenant string, skip, limit int, filterCypher string, filterParams map[string]interface{}, sorting *utils.Cypher) (*utils.DbNodesWithTotalCount, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetPaginatedInvoices")
	defer spans.Finish()

	spans.LogKV("skip", skip)
	spans.LogKV("limit", limit)
	spans.LogKV("filterCypher", filterCypher)
	spans.LogObjectAsJson("filterParams", filterParams)
	spans.LogObjectAsJson("sorting", sorting)

	sortFragment := string(*sorting)
	if sortFragment == "" {
		sortFragment = "ORDER BY i.number ASC"
	} else {
		sortFragment += ", i.number ASC"
	}

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	dbNodesWithTotalCount := new(utils.DbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		countParams := map[string]any{
			"tenant": tenant,
		}
		queryParams := map[string]any{
			"tenant": tenant,
			"skip":   skip,
			"limit":  limit,
		}

		utils.MergeMapToMap(filterParams, countParams)

		countCypher := fmt.Sprintf(` MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_INVOICE]->(i:Invoice_%s) 
				 %s 
				 RETURN count(i) as count`, tenant, tenant, tenant, filterCypher)

		spans.LogKV("countCypher", countCypher)
		spans.LogObjectAsJson("countParams", countParams)

		queryResult, err := tx.Run(ctx, countCypher, countParams)
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		dbNodesWithTotalCount.Count = count.Values[0].(int64)

		utils.MergeMapToMap(filterParams, queryParams)

		cypher := fmt.Sprintf(` MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization_%s)-[:HAS_CONTRACT]->(c:Contract_%s)-[:HAS_INVOICE]->(i:Invoice_%s) 
				 %s 
				 WITH c, i 
				 %s 
				 RETURN i
				 SKIP $skip LIMIT $limit`, tenant, tenant, tenant, filterCypher, sortFragment)

		spans.LogKV("cypher", cypher)
		spans.LogObjectAsJson("queryParams", queryParams)

		queryResult, err = tx.Run(ctx, cypher, queryParams)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		dbNodesWithTotalCount.Nodes = append(dbNodesWithTotalCount.Nodes, utils.NodePtr(v.Values[0].(neo4j.Node)))
	}
	return dbNodesWithTotalCount, nil
}

func (r *invoiceReadRepository) GetInvoiceById(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, invoiceId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoiceById")
	defer spans.Finish()

	spans.LogKV("invoiceId", invoiceId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) RETURN i`
	params := map[string]any{
		"tenant": tenant,
		"id":     invoiceId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetInvoicesByIds(ctx context.Context, tenant string, ids []string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoicesByIds")
	defer spans.Finish()

	spans.LogObjectAsJson("ids", ids)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice) WHERE i.id IN $ids RETURN i`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetAllNonDryRunInvoices(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetAllNonDryRunInvoices")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {dryRun:false}) RETURN i`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetInvoiceByIdAcrossAllTenants(ctx context.Context, invoiceId string) (*dbtype.Node, string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoiceByIdAcrossAllTenants")
	defer spans.Finish()

	spans.LogKV("invoiceId", invoiceId)

	cypher := `MATCH (t:Tenant)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {id:$id}) RETURN i, t.name`
	params := map[string]any{
		"id": invoiceId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)
	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, "", err
	}

	convertedResult, _ := result.([]*utils.DbNodeAndTenant)
	if len(convertedResult) == 0 {
		spans.LogKV("result.found", false)
		return nil, "", nil
	}
	spans.LogKV("result.found", true)
	spans.LogKV("result.tenant", convertedResult[0].Tenant)
	return convertedResult[0].Node, convertedResult[0].Tenant, err
}

func (r *invoiceReadRepository) GetInvoiceByNumber(ctx context.Context, tenant, invoiceNumber string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoiceByNumber")
	defer spans.Finish()

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice {number:$number}) RETURN i limit 1`
	params := map[string]any{
		"tenant": tenant,
		"number": invoiceNumber,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetInvoicesForPayNotifications(ctx context.Context, minutesFromCreate, minutesFromLastAttempt, lookbackWindow, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoicesForPayNotifications")
	defer spans.Finish()

	spans.LogKV("minutesFromCreate", minutesFromCreate)
	spans.LogKV("minutesFromLastAttempt", minutesFromLastAttempt)
	spans.LogKV("lookbackWindow", lookbackWindow)
	spans.LogKV("limit", limit)

	cypher := `MATCH (i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE 
				i.dryRun = false AND
				i.status IN $statuses AND
				(i.techPayNotificationRequestedAt IS NULL OR i.techPayNotificationRequestedAt + duration({minutes: $minutesFromLastAttempt}) < datetime()) AND
				i.customerEmail IS NOT NULL AND
				i.customerEmail <> '' AND	
				i.techPayInvoiceNotificationSentAt IS NULL AND
				i.createdAt+duration({days: $lookbackWindow}) > datetime() AND
				(i.updatedAt + duration({minutes: $minutesFromCreate}) < datetime())
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"minutesFromLastAttempt": minutesFromLastAttempt,
		"minutesFromCreate":      minutesFromCreate,
		"lookbackWindow":         lookbackWindow,
		"limit":                  limit,
		"statuses": []string{
			neo4jenum.InvoiceStatusDue.String(), neo4jenum.InvoiceStatusOverdue.String(), neo4jenum.InvoiceStatusPaymentProcessing.String(),
		},
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetInvoicesForPastDueNotifications(ctx context.Context, tenant string, referenceTime time.Time, overdueDays, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoicesForPastDueNotifications")
	defer spans.Finish()

	spans.LogObjectAsJson("referenceTime", referenceTime)
	spans.LogKV("overdueDays", overdueDays)
	spans.LogKV("limit", limit)

	cypher := `MATCH (i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant {name:$tenant})
			WHERE 
				i.dryRun = false AND
				i.totalAmount > 0 AND
				i.status IN $acceptedStatuses AND
				i.customerEmail IS NOT NULL AND
				i.customerEmail <> '' AND	
				date(i.dueDate) < date($referenceTime)-duration({days:$overdueDays}) AND
				i.lastRemindInvoiceNotificationSentAt IS NULL AND
				(i.techRemindInvoiceNotificationRequestedAt IS NULL OR i.techRemindInvoiceNotificationRequestedAt + duration({hours: 12}) < $referenceTime)
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"tenant":        tenant,
		"referenceTime": referenceTime,
		"overdueDays":   overdueDays,
		"limit":         limit,
		"acceptedStatuses": []string{
			neo4jenum.InvoiceStatusOverdue.String(),
		},
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) CountNonDryRunInvoicesForContract(ctx context.Context, tenant, contractId string) (int, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.CountNonDryRunInvoicesForContract")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	count, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization)-[:HAS_CONTRACT]->(c:Contract {id:$contractId})-[:HAS_INVOICE]->(i:Invoice {dryRun:false}) RETURN count(i) as count`
		params := map[string]any{
			"tenant":     tenant,
			"contractId": contractId,
		}
		spans.LogKV("cypher", cypher)
		spans.LogObjectAsJson("params", params)

		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, queryResult, err)
	})
	if err != nil {
		spans.LogKV("result.found", false)
		return 0, err
	}
	spans.LogKV("result.count", count.(int64))
	return int(count.(int64)), nil
}

func (r *invoiceReadRepository) GetPreviousCycleInvoice(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetPreviousCycleInvoice")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := `MATCH (c:Contract {id:$contractId})-[:HAS_INVOICE]->(i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			WHERE i.dryRun = false AND i.offCycle = false
			RETURN i ORDER BY i.createdAt DESC LIMIT 1`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)

	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetLastIssuedOnCycleInvoiceForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetLastIssuedOnCycleInvoiceForContract")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := `MATCH (c:Contract {id:$contractId})-[:HAS_INVOICE]->(i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			WHERE i.dryRun = false AND i.offCycle = false
			RETURN i ORDER BY i.periodEndDate DESC LIMIT 1`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)

	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetLastIssuedInvoiceForContract(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetLastIssuedInvoiceForContract")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := `MATCH (c:Contract {id:$contractId})-[:HAS_INVOICE]->(i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			WHERE i.dryRun = false
			RETURN i ORDER BY i.periodEndDate DESC LIMIT 1`
	params := map[string]any{
		"tenant":     tenant,
		"contractId": contractId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)

	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetFirstPreviewFilledInvoice(ctx context.Context, tenant, contractId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetFirstPreviewFilledInvoice")
	defer spans.Finish()

	spans.LogKV("contractId", contractId)

	cypher := `MATCH (c:Contract {id:$contractId})-[:HAS_INVOICE]->(i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			WHERE i.dryRun = true AND i.preview = true AND i.status <> $statusInitialized AND i.number IS NOT NULL AND i.number <> ''
			RETURN i ORDER BY i.createdAt DESC LIMIT 1`
	params := map[string]any{
		"tenant":            tenant,
		"contractId":        contractId,
		"statusInitialized": neo4jenum.InvoiceStatusInitialized.String(),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractFirstRecordFirstValueAsDbNodePtr(ctx, queryResult, err)

	})
	if err != nil {
		spans.LogKV("result.found", false)
		spans.TraceError(err)
		return nil, err
	}
	if result == nil {
		spans.LogKV("result.found", false)
		return nil, nil
	}
	spans.LogKV("result.found", result != nil)
	return result.(*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetExpiredDryRunInvoices(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetExpiredDryRunInvoices")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE 
				i.dryRun = true AND
				(i.preview = false OR i.preview IS NULL) AND
				i.createdAt + duration({days: 7}) < $now AND
				date(i.periodEndDate + duration({days: 7})) < date($now)
			RETURN distinct(i), t.name limit 100`
	params := map[string]any{
		"now":   utils.Now(),
		"limit": limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetPreviewInvoicesForEndedContracts(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetPreviewInvoicesForEndedContracts")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (c:Contract)--(i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE 
				i.dryRun = true AND i.preview = true AND
				c.status = $ended
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"ended": neo4jenum.ContractStatusEnded.String(),
		"limit": limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetAllForContracts(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetAllForContracts")
	defer spans.Finish()

	spans.LogObjectAsJson("contractIds", ids)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)<-[:HAS_INVOICE]-(c:Contract) 
			WHERE c.id IN $ids
			RETURN i, c.id`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *invoiceReadRepository) GetAllForServiceLineItems(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetAllForServiceLineItems")
	defer spans.Finish()

	spans.LogObjectAsJson("sliIds", ids)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)-[:HAS_INVOICE_LINE]->(:InvoiceLine)-[:INVOICED]->(sli:ServiceLineItem)
			WHERE sli.id IN $ids
			RETURN i, sli.id`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndId)))
	return result.([]*utils.DbNodeAndId), err
}

func (r *invoiceReadRepository) GetInvoicesForOverdue(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoicesForOverdue")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE 
				i.dryRun = false AND
				i.status = $dueStatus AND
				date(i.dueDate) < date(datetime())
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"dueStatus": neo4jenum.InvoiceStatusDue.String(),
		"limit":     limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetInvoicesForOnHold(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoicesForOnHold")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)<-[:HAS_INVOICE]-(c:Contract)
			WHERE 
				i.dryRun = true AND
				i.preview = true AND
				i.status = $scheduledStatus AND
				c.status = $outOfContractStatus
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"now":                 utils.Now(),
		"scheduledStatus":     neo4jenum.InvoiceStatusScheduled.String(),
		"outOfContractStatus": neo4jenum.ContractStatusOutOfContract.String(),
		"limit":               limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetInvoicesForScheduled(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetInvoicesForScheduled")
	defer spans.Finish()

	spans.LogKV("limit", limit)

	cypher := `MATCH (t:Tenant)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)<-[:HAS_INVOICE]-(c:Contract)
			WHERE 
				i.dryRun = true AND
				i.preview = true AND
				i.status = $onHoldStatus AND
				c.status <> $outOfContractStatus
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"now":                 utils.Now(),
		"onHoldStatus":        neo4jenum.InvoiceStatusOnHold.String(),
		"outOfContractStatus": neo4jenum.ContractStatusOutOfContract.String(),
		"limit":               limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetExpiredPaymentProcessingInvoices(ctx context.Context, paymentProcessingMaxDays, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetExpiredPaymentProcessingInvoices")
	defer spans.Finish()

	spans.LogKV("limit", limit)
	spans.LogKV("paymentProcessingMaxDays", paymentProcessingMaxDays)

	cypher := `MATCH (t:Tenant)<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)
			WHERE 
				i.dryRun = false AND
				i.status = $paymentProcessingStatus AND
				i.techPaymentProcessingAt IS NOT NULL AND
				i.techPaymentProcessingAt < datetime() - duration({days: $paymentProcessingMaxDays})
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"paymentProcessingStatus":  neo4jenum.InvoiceStatusPaymentProcessing.String(),
		"paymentProcessingMaxDays": paymentProcessingMaxDays,
		"limit":                    limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetReadyInvoicesForFinalizedWebhook(ctx context.Context, limit int) ([]*utils.DbNodeAndTenant, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetReadyInvoicesForFinalizedWebhook")
	defer spans.Finish()

	cypher := `MATCH (c:Contract)-[:HAS_INVOICE]->(i:Invoice)-[:INVOICE_BELONGS_TO_TENANT]->(t:Tenant)
			WHERE 
				i.dryRun = false AND
				i.amount > 0 AND
				i.techInvoiceFinalizedWebhookProcessedAt IS NULL AND
				i.updatedAt + duration({minutes: $minutesFromLastUpdate}) < datetime() AND
				i.createdAt > datetime('2025-01-01')
			RETURN distinct(i), t.name limit $limit`
	params := map[string]any{
		"minutesFromLastUpdate": 15,
		"limit":                 limit,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsAsDbNodeAndTenant(ctx, queryResult, err)

	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*utils.DbNodeAndTenant)))
	return result.([]*utils.DbNodeAndTenant), err
}

func (r *invoiceReadRepository) GetNonDryRunInvoicesForOrganization(ctx context.Context, tenant, organizationId string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetNonDryRunInvoicesForOrganization")
	defer spans.Finish()

	spans.LogKV("organizationId", organizationId)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId})-[:HAS_CONTRACT]->(c:Contract)-[:HAS_INVOICE]->(i:Invoice)
			WHERE i.dryRun = false AND i.status IN $acceptedStatuses
			RETURN i`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
		"acceptedStatuses": []string{
			neo4jenum.InvoiceStatusDue.String(),
			neo4jenum.InvoiceStatusOverdue.String(),
			neo4jenum.InvoiceStatusPaid.String(),
			neo4jenum.InvoiceStatusVoid.String(),
			neo4jenum.InvoiceStatusOnHold.String(),
			neo4jenum.InvoiceStatusPaymentProcessing.String(),
		},
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), nil
}

func (r *invoiceReadRepository) GetUpcomingInvoices(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InvoiceReadRepository.GetUpcomingInvoices")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:INVOICE_BELONGS_TO_TENANT]-(i:Invoice)
			WHERE i.preview = true AND i.dryRun = true
			RETURN i`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), nil
}
