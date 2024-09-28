defmodule CustomerOsRealtimeWeb.BankAccountChannel do
  @moduledoc """
  This Channel broadcasts sync events to all BankAccount entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "BankAccount"
end
