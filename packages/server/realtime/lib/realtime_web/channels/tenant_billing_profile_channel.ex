defmodule RealtimeWeb.TenantBillingProfileChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TenantBillingProfile entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "TenantBillingProfile"
end
