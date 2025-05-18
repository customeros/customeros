defmodule Web.InvoicesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Invoices entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Invoices"
end
