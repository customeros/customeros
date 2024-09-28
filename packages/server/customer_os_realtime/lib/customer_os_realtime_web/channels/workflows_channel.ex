defmodule CustomerOsRealtimeWeb.WorkFlowsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all WorkFlows entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "WorkFlows"
end
