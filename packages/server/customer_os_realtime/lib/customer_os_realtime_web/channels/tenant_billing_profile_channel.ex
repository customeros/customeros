defmodule CustomerOsRealtimeWeb.TenantBillingProfileChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TenantBillingProfile entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "TenantBillingProfile"
end
