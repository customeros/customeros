defmodule RealtimeWeb.ContractsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contracts entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Contracts"
end
