defmodule Web.ContractsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contracts entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Contracts"
end
