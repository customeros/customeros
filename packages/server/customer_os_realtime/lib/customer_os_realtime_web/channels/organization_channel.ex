defmodule CustomerOsRealtimeWeb.OrganizationChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Organization entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Organization"
end
