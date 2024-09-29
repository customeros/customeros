defmodule RealtimeWeb.FlowsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Flows entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Flows"
end
