defmodule RealtimeWeb.WorkFlowsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all WorkFlows entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "WorkFlows"
end
