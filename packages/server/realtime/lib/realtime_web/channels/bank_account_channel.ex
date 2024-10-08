defmodule RealtimeWeb.BankAccountChannel do
  @moduledoc """
  This Channel broadcasts sync events to all BankAccount entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "BankAccount"
end
