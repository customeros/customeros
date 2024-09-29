defmodule RealtimeWeb.BankAccountsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all BankAccounts entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "BankAccounts"
end
