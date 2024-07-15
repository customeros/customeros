package eventstore

import (
	"context"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/openline-ai/openline-customer-os/packages/server/events/utils"
	"strings"
	"time"
)

type CommonTenantIdTempAggregate struct {
	*AggregateBase
	when *func(event Event) error
}

func (a CommonTenantIdTempAggregate) NotFound() bool {
	return a.GetVersion() < 0
}

func NewCommonTempAggregateWithTenantAndId(aggregateType AggregateType, tenant, id string) *CommonTenantIdTempAggregate {
	if id == "" {
		return nil
	}
	aggregate := NewCommonTempAggregate(aggregateType)
	aggregate.SetID(fmt.Sprintf("%s-%s-%s", utils.StreamTempPrefix, tenant, id))
	return aggregate
}

func (ca *CommonTenantIdTempAggregate) setWhen(when func(event Event) error) {
	ca.when = &when
}

func GetTempAggregateWithTenantAndIdObjectID(aggregateID string, aggregateType AggregateType, tenant string) string {
	return strings.ReplaceAll(aggregateID, string(aggregateType)+"-"+tenant+"-", "")
}

func LoadCommonTempAggregateWithTenantAndId(ctx context.Context, eventStore AggregateStore, aggregateType AggregateType, tenant, objectID string) (*CommonTenantIdTempAggregate, error) {
	aggregate := NewCommonTempAggregateWithTenantAndId(aggregateType, tenant, objectID)
	err := eventStore.Load(ctx, aggregate)
	if err != nil {
		return nil, err
	}
	return aggregate, nil
}

func NewCommonTempAggregate(aggregateType AggregateType) *CommonTenantIdTempAggregate {
	commonTempAggregate := &CommonTenantIdTempAggregate{}
	base := NewAggregateBase(commonTempAggregate.When)
	base.SetType(aggregateType)
	commonTempAggregate.AggregateBase = base
	return commonTempAggregate
}

func (a *CommonTenantIdTempAggregate) When(event Event) error {
	if a.when != nil {
		return (*a.when)(event)
	}
	return nil
}

func (a *CommonTenantIdTempAggregate) SetWhen(when func(event Event) error) {
	a.when = &when
}

func (a *CommonTenantIdTempAggregate) IsTemporal() bool {
	return true
}

func (a *CommonTenantIdTempAggregate) HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error) {
	return nil, nil
}

func (a *CommonTenantIdTempAggregate) PrepareStreamMetadata() esdb.StreamMetadata {
	streamMetadata := esdb.StreamMetadata{}
	streamMetadata.SetMaxCount(utils.StreamMetadataMaxCount)
	streamMetadata.SetMaxAge(time.Duration(utils.StreamMetadataMaxAgeSeconds) * time.Second)
	return streamMetadata
}
