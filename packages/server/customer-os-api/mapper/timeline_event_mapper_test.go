package mapper_test

import (
	"testing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
)

func TestMapEntityToTimelineEvent(t *testing.T) {
	// Test nil input
	result := mapper.MapEntityToTimelineEvent(nil)
	if result != nil {
		t.Errorf("Expected nil result for nil input, got %+v", result)
	}
}

func TestMapEntitiesToTimelineEvents(t *testing.T) {
	// Test nil input
	result := mapper.MapEntitiesToTimelineEvents(nil)
	if result != nil {
		t.Errorf("Expected empty result for nil input, got %+v", result)
	}
}
