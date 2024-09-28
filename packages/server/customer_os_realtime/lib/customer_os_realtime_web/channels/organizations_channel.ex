defmodule CustomerOsRealtimeWeb.OrganizationsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Organizations entity subscribers.
  """

  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Organizations"
end
