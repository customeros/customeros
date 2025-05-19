defmodule Web.OrdersChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Orders entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Orders"
end
