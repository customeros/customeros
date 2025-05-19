defmodule Web.WorkFlowsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all WorkFlows entity subscribers.
  """
  use Web.EntitiesChannelMacro, "WorkFlows"
end
