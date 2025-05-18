defmodule Web.BankAccountsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all BankAccounts entity subscribers.
  """
  use Web.EntitiesChannelMacro, "BankAccounts"
end
