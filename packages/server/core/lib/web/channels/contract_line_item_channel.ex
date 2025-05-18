defmodule Web.ContractLineItemChannel do
  @moduledoc """
  This Channel broadcasts sync events to all ContractLineItem entity subscribers.
  """
  use Web.EntityChannelMacro, "ContractLineItem"
end
