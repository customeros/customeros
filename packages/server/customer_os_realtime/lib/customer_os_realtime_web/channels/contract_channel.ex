defmodule CustomerOsRealtimeWeb.ContractChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contract entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Contract"
end
