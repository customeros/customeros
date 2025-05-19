defmodule Web.TenantBillingProfilesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all TenantBillingProfiles entity subscribers.
  """
  use Web.EntitiesChannelMacro, "TenantBillingProfiles"
end
