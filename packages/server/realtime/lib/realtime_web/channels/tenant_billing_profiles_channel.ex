defmodule RealtimeWeb.TenantBillingProfilesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TenantBillingProfiles entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "TenantBillingProfiles"
end
