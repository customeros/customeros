defmodule CustomerOsRealtimeWeb.ContractLineItemsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all ContractLineItems entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "ContractLineItems"
end
