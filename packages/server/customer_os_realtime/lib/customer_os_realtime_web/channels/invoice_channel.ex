defmodule CustomerOsRealtimeWeb.InvoiceChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Invoice entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Invoice"
end
