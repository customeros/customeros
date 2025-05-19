defmodule Web.BankAccountChannel do
  @moduledoc """
  This Channel broadcasts sync events to all BankAccount entity subscribers.
  """
  use Web.EntityChannelMacro, "BankAccount"
end
