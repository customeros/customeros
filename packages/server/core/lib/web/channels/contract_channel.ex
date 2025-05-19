defmodule Web.ContractChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Contract entity subscribers.
  """
  use Web.EntityChannelMacro, "Contract"
end
