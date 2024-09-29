defmodule RealtimeWeb.FlowChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Flow entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Flow"
end
