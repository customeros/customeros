defmodule CustomerOsRealtimeWeb.WorkFlowChannel do
  @moduledoc """
  This Channel broadcasts sync events to all WorkFlow entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "WorkFlow"
end
