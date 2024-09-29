defmodule RealtimeWeb.InteractionEventChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionEvent entity subscribers.
  """
  use RealtimeWeb.EntityChannelMacro, "InteractionEvent"
end
