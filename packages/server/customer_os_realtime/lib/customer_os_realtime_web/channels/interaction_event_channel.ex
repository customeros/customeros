defmodule CustomerOsRealtimeWeb.InteractionEventChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionEvent entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "InteractionEvent"
end
