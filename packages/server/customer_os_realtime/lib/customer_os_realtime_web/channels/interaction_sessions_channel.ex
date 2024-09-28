defmodule CustomerOsRealtimeWeb.InteractionSessionsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionSessions entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntitiesChannelMacro, "InteractionSessions"
end
