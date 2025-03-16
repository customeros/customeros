package postgres_repository

import (
	"context"
	"testing"

	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

func TestSkuRepository_SaveAndGet(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	// Create a new SkuEntity to save
	sku := &postgresentity.SkuEntity{
		Tenant:   "tenant1",
		Type:     postgresentity.SkuTypeSubscription,
		Name:     "Test Sku",
		Price:    99.99,
		Archived: false,
	}

	// Save the SKU
	savedSku, err := repositories.SkuRepository.Save(ctx, sku)
	if err != nil {
		t.Fatalf("failed to save sku: %v", err)
	}
	if savedSku.ID == "" {
		t.Errorf("expected sku to have an ID, got empty")
	}

	// Retrieve the SKU
	fetched, err := repositories.SkuRepository.Get(ctx, "tenant1", savedSku.ID)
	if err != nil {
		t.Fatalf("failed to get sku: %v", err)
	}
	if fetched == nil {
		t.Fatalf("expected sku to be fetched, got nil")
	}
	if fetched.Name != "Test Sku" {
		t.Errorf("expected sku name 'Test Sku', got '%s'", fetched.Name)
	}
}

func TestSkuRepository_GetAll(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	// Create two SKU records with different archived statuses.
	sku1 := &postgresentity.SkuEntity{
		Tenant:   "tenant1",
		Type:     postgresentity.SkuTypeSubscription,
		Name:     "Active Sku",
		Price:    50.00,
		Archived: false,
	}
	sku2 := &postgresentity.SkuEntity{
		Tenant:   "tenant1",
		Type:     postgresentity.SkuTypeOneTime,
		Name:     "Archived Sku",
		Price:    150.00,
		Archived: true,
	}

	// Save both SKUs
	if _, err := repositories.SkuRepository.Save(ctx, sku1); err != nil {
		t.Fatalf("failed to save sku1: %v", err)
	}
	if _, err := repositories.SkuRepository.Save(ctx, sku2); err != nil {
		t.Fatalf("failed to save sku2: %v", err)
	}

	// Retrieve all SKUs for tenant1
	allSkus, err := repositories.SkuRepository.GetAll(ctx, "tenant1", nil)
	if err != nil {
		t.Fatalf("failed to get all skus: %v", err)
	}
	if len(allSkus) != 2 {
		t.Errorf("expected 2 skus, got %d", len(allSkus))
	}

	// Retrieve only archived SKUs
	archived := true
	archivedSkus, err := repositories.SkuRepository.GetAll(ctx, "tenant1", &archived)
	if err != nil {
		t.Fatalf("failed to get archived skus: %v", err)
	}
	if len(archivedSkus) != 1 {
		t.Errorf("expected 1 archived sku, got %d", len(archivedSkus))
	}
	if archivedSkus[0].Name != "Archived Sku" {
		t.Errorf("expected archived sku to be 'Archived Sku', got '%s'", archivedSkus[0].Name)
	}
}

func TestSkuRepository_Archive(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase()(t)

	// Create a new SKU that is not archived
	sku := &postgresentity.SkuEntity{
		Tenant:   "tenant1",
		Type:     postgresentity.SkuTypeSubscription,
		Name:     "Sku To Archive",
		Price:    75.00,
		Archived: false,
	}

	// Save the SKU
	savedSku, err := repositories.SkuRepository.Save(ctx, sku)
	if err != nil {
		t.Fatalf("failed to save sku: %v", err)
	}

	// Archive the SKU
	if err := repositories.SkuRepository.Archive(ctx, "tenant1", savedSku.ID); err != nil {
		t.Fatalf("failed to archive sku: %v", err)
	}

	// Retrieve the SKU and check its archived status
	archivedSku, err := repositories.SkuRepository.Get(ctx, "tenant1", savedSku.ID)
	if err != nil {
		t.Fatalf("failed to get sku: %v", err)
	}
	if archivedSku == nil {
		t.Fatalf("expected sku to exist after archiving")
	}
	if !archivedSku.Archived {
		t.Errorf("expected sku to be archived")
	}
}
