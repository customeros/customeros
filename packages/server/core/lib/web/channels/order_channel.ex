defmodule Web.OrderChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Order entity subscribers.
  """
  use Web.EntityChannelMacro, "Order"
end
