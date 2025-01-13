package events

import (
	"context"
	"reflect"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
)

// EventHandler represents a registered event handler

// HandlerRegistry manages event handlers
type HandlerRegistry struct {
	handlers map[string]interfaces.EventHandler
}

// NewHandlerRegistry creates a new handler registry
func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]interfaces.EventHandler),
	}
}

// NewHandler creates a new EventHandler with proper type handling
func NewHandler(eventType interface{}, handlerFunc interface{}) interfaces.EventHandler {
	// Validate the handler function signature
	handlerType := reflect.TypeOf(handlerFunc)
	if handlerType.Kind() != reflect.Func {
		panic("handlerFunc must be a function")
	}

	// Ensure the function takes a context and an event, and returns an error
	if handlerType.NumIn() != 2 ||
		handlerType.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() ||
		handlerType.NumOut() != 1 ||
		handlerType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic("handler must have signature func(context.Context, EventType) error")
	}

	// Create a wrapper function that matches the EventHandler signature
	wrappedFunc := func(ctx context.Context, event any) error {
		// Type assert the event
		eventVal := reflect.ValueOf(event)

		// Call the original handler
		results := reflect.ValueOf(handlerFunc).Call([]reflect.Value{
			reflect.ValueOf(ctx),
			eventVal,
		})

		// Return the error (if any)
		if !results[0].IsNil() {
			return results[0].Interface().(error)
		}
		return nil
	}

	return interfaces.EventHandler{
		HandlerFunc: wrappedFunc,
		EventType:   reflect.TypeOf(eventType).Name(),
		DataType:    reflect.TypeOf(eventType),
	}
}
