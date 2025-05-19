defmodule Web.DocumentsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Documents entity subscribers.
  """
  use Web.EntitiesChannelMacro, "Documents"
end
