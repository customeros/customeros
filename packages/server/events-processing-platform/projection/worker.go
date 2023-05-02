package projection

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"
	"golang.org/x/net/context"
)

type Worker func(ctx context.Context, stream *esdb.PersistentSubscription, workerID int) error
