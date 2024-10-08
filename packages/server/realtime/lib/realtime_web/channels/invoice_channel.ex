defmodule RealtimeWeb.InvoiceChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Invoice entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "Invoice"
end
