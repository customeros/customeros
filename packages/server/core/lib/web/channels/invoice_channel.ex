defmodule Web.InvoiceChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Invoice entity subscribers.
  """
  use Web.EntityChannelMacro, "Invoice"
end
