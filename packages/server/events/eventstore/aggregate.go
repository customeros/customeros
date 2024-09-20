package eventstore

import (
	"context"
	"fmt"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
)

const (
	aggregateStartVersion                = -1 // used for EventStoreDB
	aggregateAppliedEventsInitialCap     = 10
	aggregateUncommittedEventsInitialCap = 10
)

type LoadAggregateOptions struct {
	Required       bool
	SkipLoadEvents bool
}

func NewLoadAggregateOptions() *LoadAggregateOptions {
	return &LoadAggregateOptions{
		Required:       false,
		SkipLoadEvents: false,
	}
}

func NewLoadAggregateOptionsWithRequired() *LoadAggregateOptions {
	return &LoadAggregateOptions{
		Required:       true,
		SkipLoadEvents: false,
	}
}

func (o *LoadAggregateOptions) WithSkipLoadEvents() *LoadAggregateOptions {
	o.SkipLoadEvents = true
	return o
}

type When interface {
	When(event Event) error
}

type when func(event Event) error

// Apply process Aggregate Event
type Apply interface {
	Apply(event Event) error
	ApplyAll(events []Event) error
}

// Load create Aggregate state from Event's.
type Load interface {
	Load(events []Event) error
}

type Aggregate interface {
	When
	AggregateRoot
	HandleGRPCRequest(ctx context.Context, request any, params map[string]any) (any, error)
}

// AggregateRoot contains all methods of AggregateBase
type AggregateRoot interface {
	GetUncommittedEvents() []Event
	GetTenant() string
	GetID() string
	SetID(id string) *AggregateBase
	GetVersion() int64
	SetVersion(version int64)
	ClearUncommittedEvents()
	ToSnapshot()
	SetType(aggregateType AggregateType)
	GetType() AggregateType
	SetAppliedEvents(events []Event)
	GetAppliedEvents() []Event
	RaiseEvent(event Event) error
	String() string
	IsWithAppliedEvents() bool
	IsTemporal() bool
	SetStreamMetadata(streamMetadata *esdb.StreamMetadata)
	GetStreamMetadata() *esdb.StreamMetadata
	PrepareStreamMetadata() esdb.StreamMetadata
	Load
	Apply
}

// AggregateType type of the Aggregate
type AggregateType string

// AggregateBase base aggregate contains all main necessary fields
type AggregateBase struct {
	ID                string
	Tenant            string
	Version           int64
	AppliedEvents     []Event
	UncommittedEvents []Event
	Type              AggregateType
	withAppliedEvents bool
	when              when
	streamMetadata    *esdb.StreamMetadata
}

func NewAggregateBase(when when) *AggregateBase {
	if when == nil {
		return nil
	}

	return &AggregateBase{
		Version:           aggregateStartVersion,
		AppliedEvents:     make([]Event, 0, aggregateAppliedEventsInitialCap),
		UncommittedEvents: make([]Event, 0, aggregateUncommittedEventsInitialCap),
		when:              when,
		withAppliedEvents: false,
	}
}

// SetID set AggregateBase ID
func (a *AggregateBase) SetID(suffix string) *AggregateBase {
	a.ID = fmt.Sprintf("%s-%s", a.GetType(), suffix)
	return a
}

// GetID get AggregateBase ID
func (a *AggregateBase) GetID() string {
	return a.ID
}

// GetTenant get AggregateBase Tenant
func (a *AggregateBase) GetTenant() string {
	return a.Tenant
}

func (a *AggregateBase) IsTemporal() bool {
	return false
}

// SetType set AggregateBase AggregateType
func (a *AggregateBase) SetType(aggregateType AggregateType) {
	a.Type = aggregateType
}

// GetType get AggregateBase AggregateType
func (a *AggregateBase) GetType() AggregateType {
	return a.Type
}

// GetVersion get AggregateBase version
func (a *AggregateBase) GetVersion() int64 {
	return a.Version
}

// SetVersion set AggregateBase version
func (a *AggregateBase) SetVersion(version int64) {
	a.Version = version
}

// ClearUncommittedEvents clear AggregateBase uncommitted Event's
func (a *AggregateBase) ClearUncommittedEvents() {
	a.UncommittedEvents = make([]Event, 0, aggregateUncommittedEventsInitialCap)
}

// GetAppliedEvents get AggregateBase applied Event's
func (a *AggregateBase) GetAppliedEvents() []Event {
	return a.AppliedEvents
}

// SetAppliedEvents set AggregateBase applied Event's
func (a *AggregateBase) SetAppliedEvents(events []Event) {
	a.AppliedEvents = events
}

// GetUncommittedEvents get AggregateBase uncommitted Event's
func (a *AggregateBase) GetUncommittedEvents() []Event {
	return a.UncommittedEvents
}

func (a *AggregateBase) WithAppliedEvents() {
	a.withAppliedEvents = true
}

func (a *AggregateBase) IsWithAppliedEvents() bool {
	return a.withAppliedEvents
}

func (a *AggregateBase) GetStreamMetadata() *esdb.StreamMetadata {
	return a.streamMetadata
}

func (a *AggregateBase) SetStreamMetadata(streamMetadata *esdb.StreamMetadata) {
	a.streamMetadata = streamMetadata
}

func (a *AggregateBase) PrepareStreamMetadata() esdb.StreamMetadata {
	streamMetadata := esdb.StreamMetadata{}
	return streamMetadata
}

// Load add existing events from event store to aggregate using When interface method
func (a *AggregateBase) Load(events []Event) error {

	for _, evt := range events {
		if evt.GetAggregateID() != a.GetID() {
			return ErrInvalidAggregate
		}

		if err := a.when(evt); err != nil {
			return err
		}

		if a.withAppliedEvents {
			a.AppliedEvents = append(a.AppliedEvents, evt)
		}
		a.Version++
	}

	return nil
}

// Apply push event to aggregate uncommitted events using When method
func (a *AggregateBase) Apply(event Event) error {
	if event.GetAggregateID() != a.GetID() {
		return ErrInvalidAggregateID
	}

	event.SetAggregateType(a.GetType())

	if err := a.when(event); err != nil {
		return err
	}

	a.Version++
	event.SetVersion(a.GetVersion())
	a.UncommittedEvents = append(a.UncommittedEvents, event)
	return nil
}

func (a *AggregateBase) ApplyAll(events []Event) error {
	for _, event := range events {
		err := a.Apply(event)
		if err != nil {
			return err
		}
	}
	return nil
}

// RaiseEvent push event to aggregate applied events using When method, used for load directly from eventstore
func (a *AggregateBase) RaiseEvent(event Event) error {
	if event.GetAggregateID() != a.GetID() {
		return ErrInvalidAggregateID
	}
	if a.GetVersion() >= event.GetVersion() {
		return ErrInvalidEventVersion
	}

	event.SetAggregateType(a.GetType())

	if err := a.when(event); err != nil {
		return err
	}

	if a.withAppliedEvents {
		a.AppliedEvents = append(a.AppliedEvents, event)
	}

	a.Version = event.GetVersion()
	return nil
}

func (a *AggregateBase) ToSnapshot() {
	if a.withAppliedEvents {
		a.AppliedEvents = append(a.AppliedEvents, a.UncommittedEvents...)
	}
	a.ClearUncommittedEvents()
}

func (a *AggregateBase) String() string {
	return fmt.Sprintf("ID: {%s}, Version: {%v}, Type: {%v}, AppliedEvents: {%v}, UncommittedEvents: {%v}",
		a.GetID(),
		a.GetVersion(),
		a.GetType(),
		len(a.GetAppliedEvents()),
		len(a.GetUncommittedEvents()),
	)
}

func IsAggregateNotFound(aggregate Aggregate) bool {
	return aggregate.GetVersion() < 0
}
