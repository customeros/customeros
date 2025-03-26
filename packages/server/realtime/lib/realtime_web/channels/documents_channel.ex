defmodule RealtimeWeb.DocumentsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Documents entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "Documents"
end
