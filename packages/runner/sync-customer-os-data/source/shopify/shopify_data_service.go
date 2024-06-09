package shopify

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	source_entity "github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

const (
	CustomersTableSuffix = "customers"
	OrdersTableSuffix    = "orders"
)

var sourceTableSuffixByDataType = map[string][]string{
	string(common.ORGANIZATIONS): {CustomersTableSuffix},
	string(common.ORDERS):        {OrdersTableSuffix},
}

type shopifyDataService struct {
	airbyteStoreDb *config.RawDataStoreDB
	tenant         string
	instance       string
	processingIds  map[string]source.ProcessingEntity
	dataFuncs      map[common.SyncedEntityType]func(context.Context, int, string) []any
	log            logger.Logger
}

func NewShopifyDataService(airbyteStoreDb *config.RawDataStoreDB, tenant string, log logger.Logger) source.SourceDataService {
	dataService := shopifyDataService{
		airbyteStoreDb: airbyteStoreDb,
		tenant:         tenant,
		processingIds:  map[string]source.ProcessingEntity{},
		log:            log,
	}
	dataService.dataFuncs = map[common.SyncedEntityType]func(context.Context, int, string) []any{}
	dataService.dataFuncs[common.ORGANIZATIONS] = dataService.GetOrganizationsForSync
	dataService.dataFuncs[common.ORDERS] = dataService.GetOrdersForSync
	return &dataService
}

func (s *shopifyDataService) GetOrganizationsForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.ORGANIZATIONS)

	var organizations []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRawRecords(ctx, s.getDb(), 100, runId, currentEntity, sourceTableSuffix, s.tenant, s.SourceId())
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range airbyteRecords {
			if len(organizations) >= batchSize {
				break
			}
			outputJSON, err := MapOrganization(v.AirbyteData)
			organization, err := source.MapJsonToOrganization(outputJSON, v.AirbyteRawId, s.SourceId())
			if err != nil {
				organization = entity.OrganizationData{
					BaseData: entity.BaseData{
						SyncId: v.AirbyteRawId,
					},
				}
			}

			s.processingIds[v.AirbyteRawId] = source.ProcessingEntity{
				ExternalId:  organization.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}

			if organization.IsCustomer {
				organizations = append(organizations, organization)
			} else {
				err := repository.MarkAirbyteRawRecordProcessed(ctx, s.getDb(), s.tenant, currentEntity, sourceTableSuffix, v.AirbyteRawId, true, false, runId, organization.ExternalId, "Organization is not a customer")
				if err != nil {
					s.log.Errorf("error while marking %s with external reference %s as synced for %s", currentEntity, organization.ExternalId, s.SourceId())
					return nil
				}
			}
		}
	}
	return organizations
}

func (s *shopifyDataService) GetOrdersForSync(ctx context.Context, batchSize int, runId string) []any {
	s.processingIds = make(map[string]source.ProcessingEntity)
	currentEntity := string(common.ORDERS)

	var organizations []any
	for _, sourceTableSuffix := range sourceTableSuffixByDataType[currentEntity] {
		airbyteRecords, err := repository.GetAirbyteUnprocessedRawRecords(ctx, s.getDb(), batchSize, runId, currentEntity, sourceTableSuffix, s.tenant, s.SourceId())
		if err != nil {
			s.log.Error(err)
			return nil
		}
		for _, v := range airbyteRecords {
			if len(organizations) >= batchSize {
				break
			}
			outputJSON, err := MapOrder(v.AirbyteData)
			order, err := source.MapJsonToOrder(outputJSON, v.AirbyteRawId, s.SourceId())
			if err != nil {
				order = entity.OrderData{
					BaseData: entity.BaseData{
						SyncId: v.AirbyteRawId,
					},
				}
			}

			s.processingIds[v.AirbyteRawId] = source.ProcessingEntity{
				ExternalId:  order.ExternalId,
				Entity:      currentEntity,
				TableSuffix: sourceTableSuffix,
			}
			organizations = append(organizations, order)
		}
	}
	return organizations
}

func (s *shopifyDataService) GetDataForSync(ctx context.Context, dataType common.SyncedEntityType, batchSize int, runId string) []interface{} {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ShopifyDataService.GetDataForSync")
	defer span.Finish()
	tracing.SetDefaultSyncServiceSpanTags(ctx, span)
	span.LogFields(log.String("dataType", string(dataType)), log.Int("batchSize", batchSize))

	if ok := s.dataFuncs[dataType]; ok != nil {
		return s.dataFuncs[dataType](ctx, batchSize, runId)
	} else {
		s.log.Infof("No %s data function for %s", s.SourceId(), dataType)
		return nil
	}
}

func (s *shopifyDataService) Init() {
	err := s.getDb().AutoMigrate(&source_entity.SyncStatusForAirbyte{})
	if err != nil {
		s.log.Error(err)
	}
}

func (s *shopifyDataService) getDb() *gorm.DB {
	//schemaName := s.SourceId()
	//
	//if len(s.instance) > 0 {
	//	schemaName = schemaName + "_" + s.instance
	//}
	//schemaName = schemaName + "_" + s.tenant
	return s.airbyteStoreDb.GetDBHandler(&config.Context{
		Schema: "airbyte_internal",
	})
}

func (s *shopifyDataService) SourceId() string {
	return string(entity.AirbyteSourceShopify)
}

func (s *shopifyDataService) Close() {
	s.processingIds = make(map[string]source.ProcessingEntity)
}

func (s *shopifyDataService) MarkProcessed(ctx context.Context, syncId, runId string, synced, skipped bool, reason string) error {
	v, ok := s.processingIds[syncId]
	if ok {
		err := repository.MarkAirbyteRawRecordProcessed(ctx, s.getDb(), s.tenant, v.Entity, v.TableSuffix, syncId, synced, skipped, runId, v.ExternalId, reason)
		if err != nil {
			s.log.Errorf("error while marking %s with external reference %s as synced for %s", v.Entity, v.ExternalId, s.SourceId())
		}
		return err
	}
	return nil
}
