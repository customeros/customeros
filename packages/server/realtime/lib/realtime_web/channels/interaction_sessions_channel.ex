defmodule RealtimeWeb.InteractionSessionsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionSessions entity subscribers.
  """
  use RealtimeWeb.EntitiesChannelMacro, "InteractionSessions"
end
