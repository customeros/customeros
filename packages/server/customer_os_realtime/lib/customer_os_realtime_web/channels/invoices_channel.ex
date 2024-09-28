defmodule CustomerOsRealtimeWeb.InvoicesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Invoices entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "Invoices"
end
