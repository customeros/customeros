defmodule RealtimeWeb.ContractChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contract entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Contract"
end
