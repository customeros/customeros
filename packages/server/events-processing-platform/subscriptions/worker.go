package subscriptions

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"golang.org/x/net/context"
)

type Worker func(ctx context.Context, sub *esdb.PersistentSubscription, workerID int) error
