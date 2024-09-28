defmodule CustomerOsRealtimeWeb.ActionChannel do
  @moduledoc """
  This Channel broadcasts sync events to all Action entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "Action"
end
