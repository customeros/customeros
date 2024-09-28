defmodule CustomerOsRealtimeWeb.ContractsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contracts entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Contracts"
end
