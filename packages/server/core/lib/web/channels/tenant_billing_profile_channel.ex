defmodule Web.TenantBillingProfileChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TenantBillingProfile entity subscribers.
  """
  use Web.EntityChannelMacro, "TenantBillingProfile"
end
