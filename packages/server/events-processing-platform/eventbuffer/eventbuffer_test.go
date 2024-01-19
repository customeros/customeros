package eventbuffer

import (
	"context"
	"testing"
	"time"

	"os"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	testEventStore "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test/eventstore"
	"github.com/stretchr/testify/require"

	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/test"
)

var testDatabase *test.TestDatabase

func TestMain(m *testing.M) {
	myDatabase, shutdown := test.SetupTestDatabase()
	testDatabase = &myDatabase

	defer shutdown()

	os.Exit(m.Run())
}

// func tearDownTestCase(ctx context.Context, database *test.TestDatabase) func(tb testing.TB) {
// 	return func(tb testing.TB) {
// 		tb.Logf("Teardown test %v, cleaning neo4j DB", tb.Name())
// 		neo4jtest.CleanupAllData(ctx, database.Driver)
// 	}
// }

func TestEventBufferWatcher(t *testing.T) {
	timeNow := time.Now().UTC()
	evt := eventstore.Event{
		EventType:     "example-event-type",
		Data:          []byte("example-event-data"),
		EventID:       "example-event-id",
		Timestamp:     timeNow,
		AggregateType: "example-event-aggregate-type",
		AggregateID:   "example-event-aggregate-id",
	}
	tenant := "example-tenant"
	testLogger := test.SetupTestLogger()
	repositories := testDatabase.Repositories
	es := testEventStore.NewTestAggregateStore()

	eb := NewEventBufferWatcher(repositories, testLogger, es)

	ctx := context.Background()

	t.Run("StartStop", func(t *testing.T) {
		eb.Start(ctx)
		require.NotNil(t, eb.ticker)
		require.NotNil(t, eb.signalChannel)
		eb.Stop()
		require.Nil(t, eb.signalChannel)
	})

	t.Run("Park", func(t *testing.T) {
		tenant := "example-tenant"
		uuid := "example-uuid"
		expiryTimestamp := time.Now().UTC()

		err := eb.Park(ctx, evt, tenant, uuid, expiryTimestamp)
		require.NoError(t, err)
		parkedEvtBuffer, err := eb.repositories.EventBufferRepository.GetByUUID(uuid)
		require.NoError(t, err)
		require.Equal(t, evt.AggregateID, parkedEvtBuffer.EventAggregateID)
		require.Equal(t, string(evt.AggregateType), parkedEvtBuffer.EventAggregateType)
		require.Equal(t, evt.Data, parkedEvtBuffer.EventData)
		require.Equal(t, evt.EventID, parkedEvtBuffer.EventID)
		require.Equal(t, evt.EventType, parkedEvtBuffer.EventType)
		require.Equal(t, evt.Metadata, parkedEvtBuffer.EventMetadata)
		// require.Equal(t, evt.Timestamp, parkedEvtBuffer.EventTimestamp.UTC())
		require.Equal(t, evt.Version, parkedEvtBuffer.EventVersion)
		// require.Equal(t, expiryTimestamp, parkedEvtBuffer.ExpiryTimestamp.UTC())
		require.Equal(t, uuid, parkedEvtBuffer.UUID)
		require.Equal(t, tenant, parkedEvtBuffer.Tenant)
	})

	t.Run("Dispatch", func(t *testing.T) {
		uuid1 := "example-uuid1"
		uuid2 := "example-uuid2"
		expiryTimestamp := time.Now().UTC()

		err := eb.Park(ctx, evt, tenant, uuid1, expiryTimestamp)
		require.NoError(t, err)
		err = eb.Park(ctx, evt, tenant, uuid2, expiryTimestamp)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
		err = eb.Dispatch(ctx)
		require.NoError(t, err)
		expired, err := eb.repositories.EventBufferRepository.GetByExpired(time.Now().UTC())
		require.NoError(t, err)
		require.Equal(t, 0, len(expired))

	})

	t.Run("DispatchByUUID", func(t *testing.T) {
		uuid := "example-uuid"
		expiryTimestamp := time.Now().UTC()

		err := eb.Park(ctx, evt, tenant, uuid, expiryTimestamp)
		require.NoError(t, err)

		err = eb.DispatchByUUID(ctx, uuid)
		require.NoError(t, err)
		expired, err := eb.repositories.EventBufferRepository.GetByExpired(time.Now().UTC())
		require.NoError(t, err)
		require.Equal(t, 0, len(expired))
	})
}
