package caches

import (
	"testing"
)

func TestCache_AddExternalSystem(t *testing.T) {
	cache := NewCache()

	tenant := "tenant1"
	externalSystem := "system1"

	// Add external system to cache
	cache.AddExternalSystem(tenant, externalSystem)

	// Check if the external system exists in cache
	if !cache.CheckExternalSystem(tenant, externalSystem) {
		t.Errorf("Expected external system to be added to cache, but it was not found")
	}

	// Check if a non-existing external system returns false
	if cache.CheckExternalSystem(tenant, "nonexisting") {
		t.Errorf("Expected non-existing external system to return false, but it returned true")
	}
}

func TestCache_CheckExternalSystem(t *testing.T) {
	cache := NewCache()

	tenant := "tenant1"
	externalSystem := "system1"

	// Check non-existing external system in cache
	if cache.CheckExternalSystem(tenant, externalSystem) {
		t.Errorf("Expected non-existing external system to return false, but it returned true")
	}

	// Add external system to cache
	cache.AddExternalSystem(tenant, externalSystem)

	// Check if the external system exists in cache
	if !cache.CheckExternalSystem(tenant, externalSystem) {
		t.Errorf("Expected external system to be found in cache, but it was not found")
	}
}

func TestAddTenant(t *testing.T) {
	cache := NewCache()

	// Positive test case
	tenant := "tenant1"
	cache.AddTenant(tenant)

	// Verify that the cache contains the added tenant
	if !cache.CheckTenant(tenant) {
		t.Errorf("Expected tenant %s to be added to the cache, but it was not found", tenant)
	}

	// Negative test case
	tenant = "tenant2"

	// Verify that the cache does not contain the tenant that was not added
	if cache.CheckTenant(tenant) {
		t.Errorf("Expected tenant %s to not be added to the cache, but it was found", tenant)
	}
}
