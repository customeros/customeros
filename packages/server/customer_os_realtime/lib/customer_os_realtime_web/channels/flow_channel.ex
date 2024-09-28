defmodule CustomerOsRealtimeWeb.FlowChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Flow entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Flow"
end
