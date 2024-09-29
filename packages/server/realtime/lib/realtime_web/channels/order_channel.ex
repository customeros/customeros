defmodule RealtimeWeb.OrderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Order entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Order"
end
