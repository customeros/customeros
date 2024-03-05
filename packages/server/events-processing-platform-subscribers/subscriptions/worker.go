package subscriptions

import (
	"context"
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
)

type Worker func(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error
