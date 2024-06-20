package entity

type EventStoreAggregate struct {
	AggregateVersion *int64 `json:"aggregate_version"`
}
