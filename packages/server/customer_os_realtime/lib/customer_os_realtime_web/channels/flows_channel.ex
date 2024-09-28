defmodule CustomerOsRealtimeWeb.FlowsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Flows entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Flows"
end
