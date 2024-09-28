defmodule CustomerOsRealtimeWeb.BankAccountsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all BankAccounts entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "BankAccounts"
end
