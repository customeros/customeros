defmodule CustomerOsRealtimeWeb.InteractionSessionChannel do
  @moduledoc """
  This Channel broadcasts sync events to all InteractionSession entity subscribers.
  """
  use CustomerOsRealtimeWeb.EntityChannelMacro, "InteractionSession"
end
