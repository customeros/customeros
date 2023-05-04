package eventstore

import (
	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/pkg/errors"
)

var (
	ErrAlreadyExists       = errors.New("Already exists")
	ErrAggregateNotFound   = errors.New("aggregate not found")
	ErrInvalidEventType    = errors.New("invalid event type")
	ErrInvalidCommandType  = errors.New("invalid command type")
	ErrInvalidAggregate    = errors.New("invalid aggregate")
	ErrInvalidAggregateID  = errors.New("invalid aggregate id")
	ErrInvalidEventVersion = errors.New("invalid event version")

	ErrMissingTenant = errors.New("missing tenant")
)

func IsErrEsResourceNotFound(err error) bool {
	esdbErr, ok := esdb.FromError(err)
	if ok {
		return false
	}
	errorCode := esdbErr.Code()
	if errorCode == esdb.ErrorCodeResourceNotFound {
		return true
	}
	return false
}
