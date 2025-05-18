defmodule Web.InteractionSessionsChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionSessions entity subscribers.
  """
  use Web.EntitiesChannelMacro, "InteractionSessions"
end
