defmodule RealtimeWeb.ContractLineItemsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all ContractLineItems entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "ContractLineItems"
end
