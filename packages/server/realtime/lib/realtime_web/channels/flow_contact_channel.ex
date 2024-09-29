defmodule RealtimeWeb.FlowContactChannel do
  @moduledoc """
  This Channel broadcasts sync events to all FlowContact entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "FlowContact"
end
