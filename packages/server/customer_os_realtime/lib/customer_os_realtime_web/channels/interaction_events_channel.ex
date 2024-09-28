defmodule CustomerOsRealtimeWeb.InteractionEventsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionEvents entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "InteractionEvents"
end
