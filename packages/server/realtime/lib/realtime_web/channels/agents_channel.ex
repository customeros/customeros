defmodule RealtimeWeb.AgentsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Agents entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Agents"
end
