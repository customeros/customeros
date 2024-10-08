defmodule RealtimeWeb.ContractLineItemChannel do
  @moduledoc """
  This Channel broadcasts sync events to all ContractLineItem entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "ContractLineItem"
end
