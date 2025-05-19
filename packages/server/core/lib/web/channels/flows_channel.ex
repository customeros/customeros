defmodule Web.FlowsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Flows entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Flows"
end
