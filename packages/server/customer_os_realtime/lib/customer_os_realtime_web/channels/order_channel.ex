defmodule CustomerOsRealtimeWeb.OrderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Order entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Order"
end
