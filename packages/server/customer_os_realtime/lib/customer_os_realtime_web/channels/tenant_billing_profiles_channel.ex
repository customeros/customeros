defmodule CustomerOsRealtimeWeb.TenantBillingProfilesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TenantBillingProfiles entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "TenantBillingProfiles"
end
