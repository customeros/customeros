defmodule RealtimeWeb.InteractionEventsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionEvents entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "InteractionEvents"
end
