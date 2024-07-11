package grpcErrors

import (
	"context"
	"database/sql"
	"github.com/openline-ai/openline-customer-os/packages/server/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

var (
	ErrNoCtxMetaData = errors.New("No ctx metadata")
	ErrBadRequest    = errors.New("Bad request")
	ErrMissingFields = errors.New("Missing fields")
)

func ErrMissingField(fieldName string) error {
	return errors.Errorf("missing required field: %s", fieldName)
}

// ErrResponse get gRPC error response
func ErrResponse(err error) error {
	return status.Error(GetErrStatusCode(err), err.Error())
}

// GetErrStatusCode get error status code from error
func GetErrStatusCode(err error) codes.Code {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return codes.NotFound
	case errors.Is(err, context.Canceled):
		return codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		return codes.DeadlineExceeded
	case errors.Is(err, ErrNoCtxMetaData):
		return codes.Unauthenticated
	case errors.Is(err, ErrBadRequest):
		return codes.InvalidArgument
	case errors.Is(err, ErrMissingFields):
		return codes.InvalidArgument
	case CheckErrMessage(err, events.Validate):
		return codes.InvalidArgument
	case CheckErrMessage(err, events.Redis):
		return codes.NotFound
	case CheckErrMessage(err, events.FieldValidation):
		return codes.InvalidArgument
	case CheckErrMessage(err, events.RequiredHeaders):
		return codes.Unauthenticated
	case CheckErrMessage(err, events.Base64):
		return codes.InvalidArgument
	case CheckErrMessage(err, events.Unmarshal):
		return codes.InvalidArgument
	case CheckErrMessage(err, events.Uuid):
		return codes.InvalidArgument
	case CheckErrMessage(err, events.Cookie):
		return codes.Unauthenticated
	case CheckErrMessage(err, events.Token):
		return codes.Unauthenticated
	case CheckErrMessage(err, events.Bcrypt):
		return codes.InvalidArgument
	case eventstore.IsEventStoreErrorCodeResourceNotFound(err), errors.Is(err, eventstore.ErrAggregateNotFound):
		return codes.NotFound
	case CheckErrMessage(err, "missing required field"):
		return codes.InvalidArgument
	}
	return codes.Internal
}

func CheckErrMessage(err error, msg string) bool {
	return checkErrMessages(err, msg)
}

func checkErrMessages(err error, messages ...string) bool {
	for _, message := range messages {
		if strings.Contains(strings.TrimSpace(strings.ToLower(err.Error())), strings.TrimSpace(strings.ToLower(message))) {
			return true
		}
	}
	return false
}
