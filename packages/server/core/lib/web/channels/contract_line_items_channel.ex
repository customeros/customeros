defmodule Web.ContractLineItemsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all ContractLineItems entity subscribers.
  """
  use Web.EntitiesChannelMacro, "ContractLineItems"
end
