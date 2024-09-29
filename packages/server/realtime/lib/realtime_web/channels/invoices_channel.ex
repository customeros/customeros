defmodule RealtimeWeb.InvoicesChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Invoices entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Invoices"
end
