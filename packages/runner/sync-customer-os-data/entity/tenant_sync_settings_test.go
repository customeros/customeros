package entity

import (
	"testing"
	"time"
)

func TestContactSyncSettings(t *testing.T) {
	t.Run("Test RawDataSource Scan and Value methods", func(t *testing.T) {
		var as RawDataSource
		err := as.Scan("hubspot")
		if err != nil {
			t.Errorf("Failed to scan value: %v", err)
		}

		v, err := as.Value()
		if err != nil {
			t.Errorf("Failed to get value: %v", err)
		}

		if v != "hubspot" {
			t.Errorf("Unexpected value: got %v, want %v", v, "hubspot")
		}
	})

	t.Run("Test TenantSyncSettings fields and methods", func(t *testing.T) {
		css := TenantSyncSettings{
			ID:        1,
			CreatedAt: time.Now(),
			Tenant:    "tenant1",
			Source:    string(AirbyteSourceHubspot),
			Enabled:   true,
		}

		if css.ID != 1 {
			t.Errorf("Unexpected ID: got %v, want %v", css.ID, 1)
		}

		if css.Tenant != "tenant1" {
			t.Errorf("Unexpected tenant: got %v, want %v", css.Tenant, "tenant1")
		}

		if css.Source != string(AirbyteSourceHubspot) {
			t.Errorf("Unexpected source: got %v, want %v", css.Source, AirbyteSourceHubspot)
		}

		if css.Enabled != true {
			t.Errorf("Unexpected source: got %v, want %v", css.Enabled, true)
		}

		if css.TableName() != "tenant_sync_settings" {
			t.Errorf("Unexpected table name: got %v, want %v", css.TableName(), "contact_sync_settings")
		}
	})
}
