defmodule RealtimeWeb.OrdersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Orders entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Orders"
end
