package eventstore

import (
	"context"
	"fmt"
)

const (
	aggregateStartVersion                = -1 // used for EventStoreDB
	aggregateAppliedEventsInitialCap     = 10
	aggregateUncommittedEventsInitialCap = 10
)

// HandleCommand Aggregate commands' handler method
// Example
//
//	func (a *ContactAggregate) HandleCommand(command interface{}) error {
//		switch c := command.(type) {
//		case *CreateContactCommand:
//			return a.handleCreateContactCommand(c)
//		case *ContactPaidCommand:
//			return a.handleContactPaidCommand(c)
//		case *SubmitContactCommand:
//			return a.handleSubmitContactCommand(c)
//		default:
//			return errors.New("invalid command type")
//		}
//	}
type HandleCommand interface {
	HandleCommand(ctx context.Context, command Command) error
}

// When process and update aggregate state on specified es.Event type
// Example:
//
// func (a *ContactAggregate) When(evt es.Event) error {
//
//		switch evt.GetEventType() {
//
//		case events.ContactCreated:
//			var eventData events.ContactCreatedEvent
//			if err := json.Unmarshal(evt.GetData(), &eventData); err != nil {
//				return err
//			}
//			a.Contact.ItemsIDs = eventData.ItemsIDs
//			a.Contact.Created = true
//			return nil
//
//		default:
//			return errors.New("invalid event type")
//		}
//	}
type When interface {
	When(event Event) error
}

type when func(event Event) error

// Apply process Aggregate Event
type Apply interface {
	Apply(event Event) error
}

// Load create Aggregate state from Event's.
type Load interface {
	Load(events []Event) error
}

type Aggregate interface {
	When
	AggregateRoot
}

// AggregateRoot contains all methods of AggregateBase
type AggregateRoot interface {
	GetUncommittedEvents() []Event
	GetID() string
	SetID(id string) *AggregateBase
	GetVersion() int64
	ClearUncommittedEvents()
	ToSnapshot()
	SetType(aggregateType AggregateType)
	GetType() AggregateType
	SetAppliedEvents(events []Event)
	GetAppliedEvents() []Event
	RaiseEvent(event Event) error
	String() string
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
}

// NewAggregateBase AggregateBase constructor, contains all main fields and methods,
// main aggregate must realize When interface and pass as argument to constructor
// Example of recommended aggregate constructor method:
//
//	func NewContactAggregate() *ContactAggregate {
//		ContactAggregate := &ContactAggregate{
//			Contact: models.NewContact(),
//		}
//		base := es.NewAggregateBase(ContactAggregate.When)
//		base.SetType(ContactAggregateType)
//		ContactAggregate.AggregateBase = base
//		return ContactAggregate
//	}
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
func (a *AggregateBase) SetID(id string) *AggregateBase {
	a.ID = fmt.Sprintf("%s-%s", a.GetType(), id)
	return a
}

// GetID get AggregateBase ID
func (a *AggregateBase) GetID() string {
	return a.ID
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

// ToSnapshot prepare AggregateBase for saving Snapshot.
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
