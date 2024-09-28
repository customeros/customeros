defmodule CustomerOsRealtimeWeb.ContractLineItemChannel do
  @moduledoc """
  This Channel broadcasts sync events to all ContractLineItem entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "ContractLineItem"
end
