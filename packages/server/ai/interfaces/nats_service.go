package interfaces

import "context"

type NatsService interface {
	Start(ctx context.Context) error
	Stop()
}
