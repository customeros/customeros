defmodule CustomerOsRealtimeWeb.OrdersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Orders entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Orders"
end
