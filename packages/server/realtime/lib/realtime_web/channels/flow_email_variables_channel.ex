defmodule RealtimeWeb.FlowEmailVariablesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowEmailVariables entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "FlowEmailVariables"
end
