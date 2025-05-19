defmodule Web.FlowEmailVariablesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowEmailVariables entity subscribers.
  """
  use Web.EntitiesChannelMacro, "FlowEmailVariables"
end
