package validator

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"reflect"
)

func Validate(obj any, span opentracing.Span) (error, bool) {
	if err := GetValidator().Struct(obj); err != nil {
		// Use reflection to get the name of the type of cmd.
		typeOfCmd := reflect.TypeOf(obj)
		if typeOfCmd.Kind() == reflect.Ptr {
			typeOfCmd = typeOfCmd.Elem() // Dereference the pointer to get the actual type.
		}
		objectName := typeOfCmd.Name()

		wrappedErr := errors.Wrap(err, fmt.Sprintf("failed validation for %s", objectName))
		tracing.TraceErr(span, wrappedErr)
		return wrappedErr, true
	}
	return nil, false
}
