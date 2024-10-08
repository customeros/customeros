defmodule RealtimeWeb.OrganizationsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Organizations entity subscribers.
  """

  use RealtimeWeb.EntitiesChannelMacro, "Organizations"
end
